package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/go-chi/chi"
	"github.com/joho/godotenv"

	_ "github.com/lib/pq"
)

var db *sql.DB

func main() {
	godotenv.Load()
	r := chi.NewRouter()

	dbURL := os.Getenv("DATABASE_URL")
	if dbURL == "" {
		log.Fatal("DATABASE_URL must be set")
	}

	var err error
	db, err = sql.Open("postgres", dbURL)
	if err != nil {
		log.Fatalf("Error opening database: %s", err)
	}

	err = db.Ping()
	if err != nil {
		log.Fatal("Cannot connect to DB: ", err)
	} else {
		fmt.Println("Connected to database!")
	}

	r.Get("/tasks", GetTasks)
	r.Post("/tasks", CreateTask)
	r.Patch("/tasks/{id}", UpdateTask)
	r.Patch("/tasks/{id}/status", UpdateTaskStatus)

	log.Fatal(http.ListenAndServe(":8080", r))
}
