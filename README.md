# Todo API

REST API для управления задачами (Todo List) с аутентификацией пользователей.  
Проект написан на **Go**, использует **PostgreSQL** и упакован в **Docker** для лёгкого запуска.

Что умеет:

- Регистрация и вход по логину/паролю
- Каждый пользователь видит только свои задачи
- Можно создавать, смотреть, отмечать выполненные и удалять задачи
- Всё работает через HTTP, авторизация — Basic Auth

Как запустить:

Если есть Docker:
git clone https://github.com/kickez8/todo-api.git
cd todo-api
docker-compose up -d --build
Сервер поднимется на http://localhost:8080

Если Docker нет:
Поднять PostgreSQL, создать базу todo_db
Выполнить go mod tidy и go run .

Как тестировать:

В папке лежит файл test.http. Если в VS Code стоит расширение REST Client, можно просто открыть этот файл и нажать Send Request над любым запросом — ответ покажется сразу рядом. Все основные запросы уже там есть: регистрация, создание задачи, получение списка, обновление статуса, удаление.

Эндпоинты:

POST /register — создать пользователя (пароль нужен)
POST /login — проверить вход (нужен пароль)
GET /tasks — список задач (нужен пароль)
POST /tasks — новая задача (нужен пароль)
PUT /tasks/update?id= — отметить done (нужен пароль)
DELETE /tasks?id= — удалить задачу (нужен пароль)

Примеры:

Регистрация:
POST /register
{"username": "ivan", "password": "123"}

Создать задачу:
POST /tasks
Authorization: Basic ivan 123
{"title": "Купить хлеб"}

Получить задачи:
GET /tasks
Authorization: Basic ivan 123

Отметить выполненной:
PUT /tasks/update?id=1
Authorization: Basic ivan 123
{"done": true}

Удалить:
DELETE /tasks?id=1
Authorization: Basic ivan 123

Структура проекта:
```text
todo-api/
├── main.go             # Точка входа, инициализация роутов
├── database.go         # Подключение к БД, миграции и таблицы
├── handlers.go         # Обработчики запросов и вся бизнес-логика
├── test.http           # Запросы для быстрого тестирования (REST Client)
├── Dockerfile          # Инструкция для сборки Docker-образа
├── docker-compose.yml  # Конфигурация локального окружения
└── .env                # Конфигурационные данные и секреты (не для Git)
```

Технологии: Go 1.21, PostgreSQL, Docker, REST Client (для тестов)

Автор: @kickez8
