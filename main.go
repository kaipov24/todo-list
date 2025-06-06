package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/go-chi/chi"
	"github.com/go-chi/cors"
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

	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		AllowCredentials: true,
	}))

	r.Get("/tasks", GetTasks)
	r.Post("/tasks", CreateTask)
	r.Patch("/tasks/{id}", UpdateTask)
	r.Patch("/tasks/{id}/status", UpdateTaskStatus)
	r.Delete("/tasks/{id}", DeleteTask)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	log.Fatal(http.ListenAndServe(":"+port, r))
}
