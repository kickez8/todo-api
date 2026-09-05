package main

import (
	"encoding/json"
	"net/http"
)

// Структура для парсинга данных пользователя при регистрации/логине
type User struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

// Структура для работы с задачами
type Task struct {
	ID    int    `json:"id"`
	Title string `json:"title"`
	Done  bool   `json:"done"`
}

// 1. Регистрация: добавление нового пользователя в базу
func RegisterHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Метод не поддерживается", http.StatusMethodNotAllowed)
		return
	}

	var u User
	err := json.NewDecoder(r.Body).Decode(&u)
	if err != nil || u.Username == "" || u.Password == "" {
		http.Error(w, "Неверные данные", http.StatusBadRequest)
		return
	}

	_, err = DB.Exec("INSERT INTO users (username, password) VALUES ($1, $2)", u.Username, u.Password)
	if err != nil {
		http.Error(w, "Пользователь уже существует", http.StatusConflict)
		return
	}

	w.WriteHeader(http.StatusCreated)
}

// 2. Авторизация: вспомогательная функция проверки HTTP Basic Auth
func authenticate(r *http.Request) (int, bool) {
	username, password, ok := r.BasicAuth()
	if !ok {
		return 0, false
	}

	var id int
	var dbPassword string
	err := DB.QueryRow("SELECT id, password FROM users WHERE username = $1", username).Scan(&id, &dbPassword)
	if err != nil || dbPassword != password {
		return 0, false
	}

	return id, true
}

// 3. Логин: проверка корректности данных входа
func LoginHandler(w http.ResponseWriter, r *http.Request) {
	_, authorized := authenticate(r)
	if !authorized {
		http.Error(w, "Неверный логин или пароль", http.StatusUnauthorized)
		return
	}
	w.Write([]byte("Успешный вход!"))
}

// 4. Менеджер задач: обработка получения, создания и удаления напоминаний
func TasksHandler(w http.ResponseWriter, r *http.Request) {
	// проверяем авторизацию перед любым действием с задачами
	userID, authorized := authenticate(r)
	if !authorized {
		w.Header().Set("WWW-Authenticate", `Basic realm="Restricted"`)
		http.Error(w, "Неавторизован", http.StatusUnauthorized)
		return
	}

	switch r.Method {
	case http.MethodGet: // получить все задачи текущего пользователя
		rows, err := DB.Query("SELECT id, title, done FROM tasks WHERE user_id = $1", userID)
		if err != nil {
			http.Error(w, "Ошибка БД", http.StatusInternalServerError)
			return
		}
		defer rows.Close()

		tasks := []Task{}
		for rows.Next() {
			var t Task
			err := rows.Scan(&t.ID, &t.Title, &t.Done)
			if err != nil {
				continue
			}
			tasks = append(tasks, t)
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(tasks)

	case http.MethodPost: // создать новую задачу
		var t Task
		err := json.NewDecoder(r.Body).Decode(&t)
		if err != nil || t.Title == "" {
			http.Error(w, "Неверное тело запроса", http.StatusBadRequest)
			return
		}

		_, err = DB.Exec("INSERT INTO tasks (user_id, title) VALUES ($1, $2)", userID, t.Title)
		if err != nil {
			http.Error(w, "Ошибка сохранения", http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusCreated)

	case http.MethodDelete: // удалить задачу
		// Ожидаем id в строке запроса, например: /tasks?id=5
		taskID := r.URL.Query().Get("id")
		if taskID == "" {
			http.Error(w, "Не указан id задачи", http.StatusBadRequest)
			return
		}

		// Удаляем задачу только если id совпадает И она принадлежит этому пользователю (user_id)
		result, err := DB.Exec("DELETE FROM tasks WHERE id = $1 AND user_id = $2", taskID, userID)
		if err != nil {
			http.Error(w, "Ошибка удаления", http.StatusInternalServerError)
			return
		}

		// Проверяем, удалилась ли хоть одна строчка (вдруг id не существует или задача чужая)
		rowsAffected, _ := result.RowsAffected()
		if rowsAffected == 0 {
			http.Error(w, "Задача не найдена или нет доступа", http.StatusNotFound)
			return
		}

		// 204 No Content означает, что всё прошло успешно, но возвращать в теле ответа нечего
		w.WriteHeader(http.StatusNoContent)

	default:
		http.Error(w, "Метод не поддерживается", http.StatusMethodNotAllowed)
	}
}
