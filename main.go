package main

import (
	"database/sql"
	"log"
	"net/http"
	"os"

	"github.com/go-chi/chi"
	"github.com/joho/godotenv"
)

var db *sql.DB

func main() {
	godotenv.Load()
	r := chi.NewRouter()

	dbURL := os.Getenv("DB_URL")
	if dbURL == "" {
		log.Fatal("DB_URL must be set")
	}

	_, err := sql.Open("postgres", dbURL)
	if err != nil {
		log.Fatalf("Error opening database: %s", err)
	}

	r.Get("/tasks", GetTasks)
	// r.Post("/tasks", CreateTask)

	log.Fatal(http.ListenAndServe(":8080", r))
}
