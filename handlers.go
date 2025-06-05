package main

import (
	"encoding/json"
	"net/http"
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
