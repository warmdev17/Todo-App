package main

import (
	"encoding/json"
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"
)

type Task struct {
	ID        int    `json:"id"`
	Title     string `json:"title"`
	Completed bool   `json:"completed"`
}

var tasks = []Task{
	{ID: 1, Title: "Learn Go net/http", Completed: false},
	{ID: 2, Title: "Build TODO REST API app", Completed: false},
	{ID: 3, Title: "Hello bro", Completed: true},
}

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	http.HandleFunc("/tasks", tasksHandler)
	log.Println("Server running on http://localhost:" + os.Getenv("APP_PORT"))
	err = http.ListenAndServe(":"+os.Getenv("APP_PORT"), nil)
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
			http.Error(w, "Invalid JSON", http.StatusBadRequest)
			return
		}
		return
	}

	if r.Method == http.MethodPost {
		var input struct {
			Title string `json:"title"`
		}

		log.Print("/POST")
		err := json.NewDecoder(r.Body).Decode(&input)
		if err != nil {
			http.Error(w, "Invalid JSON", http.StatusBadRequest)
			return
		}

		if input.Title == "" {
			http.Error(w, "Title is required", http.StatusBadRequest)
			return
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
			http.Error(w, "Invalid JSON", http.StatusBadRequest)
			return
		}

		return
	}

	http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
}
