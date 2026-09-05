package main

import (
	"database/sql"
	"log"

	_ "github.com/lib/pq"
)

var DB *sql.DB

func InitDB() {
	connStr := "user=postgres password=G174632 dbname=todo_db sslmode=disable"
	var err error
	DB, err = sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal(err)
	}

	// Создаем таблицу пользователей
	userTable := `
	CREATE TABLE IF NOT EXISTS users (
		id SERIAL PRIMARY KEY,
		username TEXT UNIQUE NOT NULL,
		password TEXT NOT NULL
	);`

	// Создаем таблицу задач (напоминаний)
	taskTable := `
	CREATE TABLE IF NOT EXISTS tasks (
		id SERIAL PRIMARY KEY,
		user_id INT REFERENCES users(id),
		title TEXT NOT NULL,
		done BOOLEAN DEFAULT FALSE
	);`

	_, err = DB.Exec(userTable)
	if err != nil {
		log.Fatal(err)
	}
	_, err = DB.Exec(taskTable)
	if err != nil {
		log.Fatal(err)
	}
}
