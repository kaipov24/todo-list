package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/go-chi/chi"
)

func GetTasks(w http.ResponseWriter, r *http.Request) {
	rows, err := db.Query("SELECT id, title, done FROM tasks ORDER BY id")
	if err != nil {
		http.Error(w, "Failed to fetch tasks", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var tasks []Task

	for rows.Next() {
		var t Task
		if err := rows.Scan(&t.ID, &t.Title, &t.Done); err != nil {
			http.Error(w, "Failed to parse task", http.StatusInternalServerError)
			return
		}
		tasks = append(tasks, t)
	}
	respondWithJSON(w, 200, tasks)

}

func CreateTask(w http.ResponseWriter, r *http.Request) {

	var t Task

	if err := json.NewDecoder(r.Body).Decode(&t); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	query := `INSERT INTO tasks (title, done) VALUES ($1, $2)`

	_, err := db.Exec(query, t.Title, t.Done)
	if err != nil {
		http.Error(w, "Failed to insert task", http.StatusInternalServerError)
		return
	}
	respondWithJSON(w, 200, t)
}

func UpdateTask(w http.ResponseWriter, r *http.Request) {

	idStr := chi.URLParam(r, "id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid task ID", http.StatusBadRequest)
		return
	}

	var input struct {
		Title *string `json:"title"`
	}
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	if input.Title == nil {
		http.Error(w, "Title field is required", http.StatusBadRequest)
		return
	}

	query := "UPDATE tasks SET title = $1 WHERE id = $2"
	result, err := db.Exec(query, *input.Title, id)
	if err != nil {
		http.Error(w, "Failed to update title: "+err.Error(), http.StatusInternalServerError)
		return
	}

	rows, _ := result.RowsAffected()
	if rows == 0 {
		http.Error(w, "Task not found", http.StatusNotFound)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func UpdateTaskStatus(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid task ID", http.StatusBadRequest)
		return
	}

	var input struct {
		Done *bool `json:"done"`
	}

	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	if input.Done == nil {
		http.Error(w, "Field 'done' is required", http.StatusBadRequest)
		return
	}

	query := "UPDATE tasks SET done = $1 WHERE id = $2"
	result, err := db.Exec(query, *input.Done, id)
	if err != nil {
		http.Error(w, "Failed to update task status: "+err.Error(), http.StatusInternalServerError)
		return
	}

	rows, err := result.RowsAffected()
	if err != nil {
		http.Error(w, "Failed to check update result", http.StatusInternalServerError)
		return
	}
	if rows == 0 {
		http.Error(w, "Task not found", http.StatusNotFound)
		return
	}

	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "Task %d marked as %v\n", id, *input.Done)
}

func DeleteTask(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid task ID", http.StatusBadRequest)
		return
	}

	query := "DELETE FROM tasks WHERE id = $1"
	_, err = db.Exec(query, id)
	if err != nil {
		http.Error(w, "Failed to update task status: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)

}
