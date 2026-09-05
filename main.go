package main

import (
	"log"
	"net/http"
)

func main() {
	InitDB()
	defer DB.Close()

	http.HandleFunc("/register", RegisterHandler)
	http.HandleFunc("/login", LoginHandler)
	http.HandleFunc("/tasks", TasksHandler)

	log.Println("Сервер запущен на порту :8080...")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
