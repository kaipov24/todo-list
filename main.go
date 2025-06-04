package main

import (
	"log"
	"net/http"

	"github.com/go-chi/chi"
)

func main() {
	r := chi.NewRouter()

	r.Get("/tasks", GetTasks)

	log.Fatal(http.ListenAndServe(":8080", r))
}
