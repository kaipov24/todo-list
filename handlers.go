package main

import (
	"encoding/json"
	"net/http"
)

func GetTasks(w http.ResponseWriter, r *http.Request) {
	tasks := []Task{{ID: 1, Title: "To do list", Done: false}, {ID: 2, Title: "Build something bigger", Done: false}}

	json.NewEncoder(w).Encode(tasks)
}
