package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	_ "github.com/lib/pq"
)

var DB *sql.DB

func InitDB() {
	// Получаем настройки из переменных окружения
	host := os.Getenv("DB_HOST")
	if host == "" {
		host = "localhost"
	}

	port := os.Getenv("DB_PORT")
	if port == "" {
		port = "5432"
	}

	user := os.Getenv("DB_USER")
	if user == "" {
		user = "postgres"
	}

	password := os.Getenv("DB_PASSWORD")
	if password == "" {
		password = "G174632"
	}

	dbname := os.Getenv("DB_NAME")
	if dbname == "" {
		dbname = "todo_db"
	}

	sslmode := os.Getenv("DB_SSLMODE")
	if sslmode == "" {
		sslmode = "disable"
	}

	// Формируем строку подключения
	connStr := fmt.Sprintf("user=%s password=%s host=%s port=%s dbname=%s sslmode=%s",
		user, password, host, port, dbname, sslmode)

	var err error
	DB, err = sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal("Ошибка открытия БД:", err)
	}

	// Проверка подключения
	err = DB.Ping()
	if err != nil {
		log.Fatal("Ошибка подключения к БД:", err)
	}

	fmt.Println("Подключено к PostgreSQL!")

	// Создаем таблицу пользователей
	userTable := `
	CREATE TABLE IF NOT EXISTS users (
		id SERIAL PRIMARY KEY,
		username TEXT UNIQUE NOT NULL,
		password TEXT NOT NULL
	);`

	// Создаем таблицу задач
	taskTable := `
	CREATE TABLE IF NOT EXISTS tasks (
		id SERIAL PRIMARY KEY,
		user_id INT REFERENCES users(id),
		title TEXT NOT NULL,
		done BOOLEAN DEFAULT FALSE
	);`

	_, err = DB.Exec(userTable)
	if err != nil {
		log.Fatal("Ошибка создания таблицы users:", err)
	}
	_, err = DB.Exec(taskTable)
	if err != nil {
		log.Fatal("Ошибка создания таблицы tasks:", err)
	}

	fmt.Println("Таблицы созданы (users, tasks)")
}
