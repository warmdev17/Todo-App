package main

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"
)

func tasksHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodOptions {
		w.WriteHeader(http.StatusNoContent)
		return
	}

	if r.Method != http.MethodGet && r.Method != http.MethodPost {
		writeError(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	currentUser, statusCode, err := getCurrentUser(r)
	if err != nil {
		writeError(w, statusCode, err.Error())
		return
	}

	if r.Method == http.MethodGet {

		tasks, err := getTasksByUserID(currentUser.ID)

		if err != nil {
			writeError(w, http.StatusInternalServerError, "Failed to get tasks")
			return
		}

		writeJSON(w, http.StatusOK, map[string]any{
			"success": true,
			"data":    tasks,
		})
		return
	}

	if r.Method == http.MethodPost {
		var input TaskCreate

		err = json.NewDecoder(r.Body).Decode(&input)
		if err != nil {
			writeError(w, http.StatusBadRequest, "Invalid JSON")
			return
		}

		if isBlankPointer(input.Title) {
			writeError(w, http.StatusBadRequest, "Title is required")
			return
		}

		trimmedTitle := strings.TrimSpace(*input.Title)

		newTask, err := createNewTask(currentUser.ID, trimmedTitle)

		if err != nil {
			writeError(w, http.StatusInternalServerError, err.Error())
			return
		}

		writeJSON(w, http.StatusCreated, map[string]any{
			"success": true,
			"data":    newTask,
		})
		return
	}
}

func taskByIDHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodOptions {
		w.WriteHeader(http.StatusNoContent)
		return
	}

	if r.Method != http.MethodGet && r.Method != http.MethodDelete && r.Method != http.MethodPatch {
		writeError(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	currentUser, statusCode, err := getCurrentUser(r)
	if err != nil {
		writeError(w, statusCode, err.Error())
		return
	}
	switch r.Method {
	case http.MethodGet:
		{
			id, err := getTaskIDFromPath(r)
			if err != nil {
				writeError(w, http.StatusBadRequest, "Invalid task ID")
				return
			}

			task, err := getTaskByIDAndUserID(id, currentUser.ID)
			if err != nil {
				writeError(w, http.StatusNotFound, err.Error())
				return
			}
			writeJSON(w, http.StatusOK, map[string]any{
				"success": true,
				"data":    task,
			})
			return
		}
	case http.MethodDelete:
		{
			id, err := getTaskIDFromPath(r)
			if err != nil {
				writeError(w, http.StatusBadRequest, "Invalid task id")
				return
			}

			deletedTask, err := deleteTaskByIDAndUserID(id, currentUser.ID)

			if err != nil {
				writeError(w, http.StatusNotFound, "Task not found")
				return
			}

			writeJSON(w, http.StatusOK, map[string]any{
				"success": true,
				"data":    deletedTask,
			})
			return
		}
	case http.MethodPatch:
		{
			var input UpdateTaskInput

			err = json.NewDecoder(r.Body).Decode(&input)
			if err != nil {
				writeError(w, http.StatusBadRequest, "Invalid JSON")
				return
			}

			if input.Title == nil && input.Completed == nil {
				writeError(w, http.StatusBadRequest, "No fields to update")
				return
			}

			taskID, err := getTaskIDFromPath(r)
			if err != nil {
				writeError(w, http.StatusBadRequest, "Invalid task id")
				return
			}

			var trimmedTitle string
			if input.Title != nil {
				trimmedTitle = strings.TrimSpace(*input.Title)
				if trimmedTitle == "" {
					writeError(w, http.StatusBadRequest, "Title cannot be empty")
					return
				}
				*input.Title = trimmedTitle
			}

			task, err := updateTaskByIDAndUserID(taskID, currentUser.ID, input)

			if err != nil {
				writeError(w, http.StatusNotFound, err.Error())
				return
			}

			writeJSON(w, http.StatusOK, map[string]any{
				"success": true,
				"data":    task,
			})

			return

		}
	}
}

func getTaskIDFromPath(r *http.Request) (int, error) {
	idText := strings.TrimPrefix(r.URL.Path, "/tasks/")
	id, err := strconv.Atoi(idText)
	if err != nil {
		return 0, err
	}

	return id, nil
}
