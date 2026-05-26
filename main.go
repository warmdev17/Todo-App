package main

import (
	"encoding/json"
	"log"
	"net/http"
)

type Task struct {
	ID        int    `json:"id"`
	Title     string `json:"title"`
	Completed bool   `json:"completed"`
}

var input struct {
	Title string `json:"title"`
}

var tasks = []Task{
	{ID: 1, Title: "Learn Go net/http", Completed: false},
	{ID: 2, Title: "Build TODO REST API app", Completed: false},
	{ID: 3, Title: "Hello bro", Completed: true},
}

func main() {
	http.HandleFunc("/tasks", tasksHandler)
	log.Print("Server running on http://localhost:8080")
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatal("Failed to start server:", err)
	}
}

func tasksHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
	w.Header().Set("Content-Type", "application/json")

	if r.Method == http.MethodOptions {
		w.WriteHeader(http.StatusNoContent)
		return
	}

	if r.Method == http.MethodGet {
		log.Print("/GET")
		err := json.NewEncoder(w).Encode(map[string]any{
			"success": true,
			"data":    tasks,
		})
		if err != nil {
			log.Fatal("Failed to write JSON", err)
		}
		return
	}

	if r.Method == http.MethodPost {
		log.Print("/POST")
		err := json.NewDecoder(r.Body).Decode(&input)
		if err != nil {
			log.Fatal("Failed to read json:", err)
		}

		newTask := Task{
			ID:        len(tasks) + 1,
			Title:     input.Title,
			Completed: false,
		}

		tasks = append(tasks, newTask)
		w.WriteHeader(http.StatusCreated)
		if err := json.NewEncoder(w).Encode(map[string]any{
			"success": true,
			"data":    newTask,
		}); err != nil {
			log.Fatal("Failed to writer json:", err)
		}

		return
	}

	http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
}
