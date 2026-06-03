package main

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/joho/godotenv"
	"golang.org/x/crypto/bcrypt"
)

type Task struct {
	ID        int    `json:"id"`
	UserID    int    `json:"userId"`
	Title     string `json:"title"`
	Completed bool   `json:"completed"`
}

type User struct {
	ID           int    `json:"id"`
	Email        string `json:"email"`
	Username     string `json:"username"`
	HashPassword string `json:"-"`
}

type AuthUser struct {
	ID       int    `json:"id"`
	Email    string `json:"email"`
	Username string `json:"username"`
}

type RegisterInput struct {
	Email           *string `json:"email"`
	Username        *string `json:"username"`
	Password        *string `json:"password"`
	ConfirmPassword *string `json:"confirmPassword"`
}

var tasks = []Task{
	{ID: 1, UserID: 1, Title: "Learn Go net/http", Completed: false},
	{ID: 2, UserID: 1, Title: "Learn Gin framework", Completed: false},
	{ID: 3, UserID: 2, Title: "Build TODO REST API app", Completed: false},
}

var users = []User{}

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	http.HandleFunc("/tasks", tasksHandler)
	http.HandleFunc("/tasks/", taskByIDHandler)
	http.HandleFunc("/login", loginHandler)
	http.HandleFunc("/register", registerHandler)
	log.Println("Server running on http://localhost:" + os.Getenv("APP_PORT"))
	err = http.ListenAndServe(":"+os.Getenv("APP_PORT"), nil)
	if err != nil {
		log.Fatal("Failed to start server:", err)
	}
}

func tasksHandler(w http.ResponseWriter, r *http.Request) {
	setCORSHeader(w, "GET, POST, OPTIONS")

	if r.Method == http.MethodOptions {
		w.WriteHeader(http.StatusNoContent)
		return
	}

	if r.Method == http.MethodGet {
		log.Printf("method = %s, path = %s", r.Method, r.URL.Path)
		userID, err := getCurrentUserID(r)
		if err != nil {
			writeError(w, http.StatusBadRequest, err.Error())
			return
		}

		var userTask []Task

		for _, task := range tasks {
			if task.UserID == userID {
				userTask = append(userTask, task)
			}
		}

		writeJSON(w, http.StatusAccepted, map[string]any{
			"success": true,
			"data":    userTask,
		})
		return
	}

	if r.Method == http.MethodPost {
		var input struct {
			Title string `json:"title"`
		}

		log.Printf("method = %s, path = %s", r.Method, r.URL.Path)
		err := json.NewDecoder(r.Body).Decode(&input)
		if err != nil {
			http.Error(w, "Invalid JSON", http.StatusBadRequest)
			return
		}

		trimmedTitle := strings.TrimSpace(input.Title)

		if trimmedTitle == "" {
			http.Error(w, "Title is required", http.StatusBadRequest)
			return
		}

		newTask := Task{
			ID:        nextTaskID(),
			Title:     trimmedTitle,
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

func taskByIDHandler(w http.ResponseWriter, r *http.Request) {
	setCORSHeader(w, "GET, DELETE, PATCH, OPTIONS")

	switch r.Method {
	case http.MethodGet:
		{
			log.Printf("method = %s, path = %s", r.Method, r.URL.Path)
			id, err := getIDFromPath(r)
			if err != nil {
				http.Error(w, "Invalid task ID", http.StatusBadRequest)
				return
			}

			task, err := getTaskByID(id)
			if err != nil {
				http.Error(w, "Task not found", http.StatusNotFound)
				return
			}

			err = json.NewEncoder(w).Encode(map[string]any{
				"success": true,
				"data":    task,
			})
			if err != nil {
				http.Error(w, "Invalid JSON", http.StatusBadRequest)
				return
			}

		}
	case http.MethodDelete:
		{
			log.Printf("method = %s, path = %s", r.Method, r.URL.Path)

			id, err := getIDFromPath(r)
			if err != nil {
				http.Error(w, "Invalid task id", http.StatusBadRequest)
				return
			}

			for index, task := range tasks {
				if task.ID == id {
					tasks = append(tasks[:index], tasks[index+1:]...)
					writeJSON(w, http.StatusOK, map[string]any{
						"success": true,
						"data":    task,
					})
					return
				}
			}
			writeError(w, http.StatusNotFound, "Task not found")
			return
		}
	case http.MethodPatch:
		{
			var input struct {
				Completed *bool   `json:"completed"`
				Title     *string `json:"title"`
			}
			log.Printf("method = %s, path = %s", r.Method, r.URL.Path)
			id, err := getIDFromPath(r)
			if err != nil {
				http.Error(w, "Invalid task id", http.StatusBadRequest)
				return
			}

			err = json.NewDecoder(r.Body).Decode(&input)
			if err != nil {
				http.Error(w, "Invalid JSON", http.StatusBadRequest)
				return
			}

			var trimmedTitle string
			if input.Title != nil {
				trimmedTitle = strings.TrimSpace(*input.Title)
				if trimmedTitle == "" {
					writeError(w, http.StatusBadRequest, "Title cannot be empty")
					return
				}
			}

			for index, task := range tasks {
				if task.ID == id {
					if input.Title != nil {
						tasks[index].Title = trimmedTitle
					}

					if input.Completed != nil {
						tasks[index].Completed = *input.Completed
					}

					w.WriteHeader(http.StatusAccepted)
					if err := json.NewEncoder(w).Encode(map[string]any{
						"success": true,
						"data":    tasks[index],
					}); err != nil {
						http.Error(w, "Invalid JSON", http.StatusBadRequest)
						return
					}

					return

				}
			}

		}
	case http.MethodOptions:
		{
			log.Printf("method = %s, path = %s", r.Method, r.URL.Path)
			w.WriteHeader(http.StatusNoContent)
			return
		}
	default:
		{
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}
	}
}

func loginHandler(w http.ResponseWriter, r *http.Request) {
	setCORSHeader(w, "GET, POST, OPTIONS")

	var LoginInput struct {
		Email    *string `json:"email"`
		Password *string `json:"password"`
		Username *string `json:"username"`
	}

	if r.Method == http.MethodOptions {
		log.Printf("method = %s, path = %s", r.Method, r.URL.Path)
		w.WriteHeader(http.StatusNoContent)
		return
	}

	if r.Method == http.MethodPost {
		log.Printf("method = %s, path = %s", r.Method, r.URL.Path)
		err := json.NewDecoder(r.Body).Decode(&LoginInput)
		if err != nil {
			http.Error(w, "Invalid JSON", http.StatusBadRequest)
			return
		}

		var user User

		if isBlankPointer(LoginInput.Password) {
			http.Error(w, "password is required", http.StatusBadRequest)
			return
		}

		hasEmail := !isBlankPointer(LoginInput.Email)
		hasUsername := !isBlankPointer(LoginInput.Username)

		if hasEmail == hasUsername {
			http.Error(w, "use either username or email", http.StatusBadRequest)
			return
		}
		if hasEmail {
			user, err = findUserByEmail(*LoginInput.Email)
			if err != nil {
				writeError(w, http.StatusUnauthorized, "Invalid credentials")
				return
			}
			err := bcrypt.CompareHashAndPassword([]byte(user.HashPassword), []byte(*LoginInput.Password))
			if err != nil {
				writeError(w, http.StatusUnauthorized, "Invalid email or password")
				return
			} else {
				writeJSON(w, http.StatusOK, map[string]any{
					"success": true,
					"data": map[string]any{
						"token":    "fake-token-1",
						"userId":   user.ID,
						"username": user.Username,
					},
				})
				return
			}
		} else {
			user, err = findUserByUsername(*LoginInput.Username)
			if err != nil {
				writeError(w, http.StatusUnauthorized, "Invalid credentials")
				return
			}

			err := bcrypt.CompareHashAndPassword([]byte(user.HashPassword), []byte(*LoginInput.Password))
			if err != nil {
				writeError(w, http.StatusUnauthorized, "Invalid username or password")
				return
			} else {
				writeJSON(w, http.StatusOK, map[string]any{
					"success": true,
					"data": map[string]any{
						"token":    "fake-token-1",
						"userId":   user.ID,
						"username": user.Username,
					},
				})
				return
			}
		}
	}
}

func registerHandler(w http.ResponseWriter, r *http.Request) {
	setCORSHeader(w, "POST, OPTIONS")

	if r.Method == http.MethodOptions {
		w.WriteHeader(http.StatusNoContent)
		return
	}

	var input RegisterInput

	if r.Method == http.MethodPost {
		log.Printf("method = %s, path = %s", r.Method, r.URL.Path)
		if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
			http.Error(w, "Invalid JSON Body", http.StatusBadRequest)
			return
		}

		if input.Username == nil &&
			input.Email == nil &&
			input.Password == nil &&
			input.ConfirmPassword == nil {
			http.Error(w, "request body is required", http.StatusBadRequest)
			return
		}

		trimmedEmail := strings.TrimSpace(*input.Email)
		trimmedUsername := strings.TrimSpace(*input.Username)

		errs := validateCreateUser(input)

		if len(errs) > 0 {
			writeError(w, http.StatusBadRequest, "Invalid JSON Body")
			return
		}

		if _, err := findUserByEmail(trimmedEmail); err == nil {
			writeError(w, http.StatusBadRequest, "email already exists")
			return
		}

		if _, err := findUserByUsername(trimmedUsername); err == nil {
			writeError(w, http.StatusBadRequest, "username already exists")
			return
		}

		hashPassword, err := bcrypt.GenerateFromPassword([]byte(*input.Password), bcrypt.DefaultCost)
		log.Println(string(hashPassword))
		if err != nil {
			writeError(w, http.StatusBadRequest, err.Error())
			return
		}

		newUser := User{
			ID:           nextUserID(),
			Username:     trimmedUsername,
			Email:        trimmedEmail,
			HashPassword: string(hashPassword),
		}
		users = append(users, newUser)
		writeJSON(w, http.StatusCreated, map[string]any{
			"success": true,
			"data": AuthUser{
				ID:       newUser.ID,
				Email:    newUser.Email,
				Username: newUser.Username,
			},
		})
	}
}

func getTaskByID(id int) (Task, error) {
	for _, task := range tasks {
		if task.ID == id {
			return task, nil
		}
	}

	return Task{}, errors.New("task not found")
}

func getIDFromPath(r *http.Request) (int, error) {
	idText := strings.TrimPrefix(r.URL.Path, "/tasks/")
	id, err := strconv.Atoi(idText)
	if err != nil {
		return 0, err
	}

	return id, nil
}

func nextTaskID() int {
	if len(tasks) == 0 {
		return 1
	}

	max := tasks[0].ID

	for _, task := range tasks {
		if task.ID > max {
			max = task.ID
		}
	}

	return max + 1
}

func nextUserID() int {
	if len(users) == 0 {
		return 1
	}

	max := users[0].ID

	for _, user := range users {
		if user.ID > max {
			max = user.ID
		}
	}

	return max + 1
}

func findUserByEmail(email string) (User, error) {
	for _, user := range users {
		if user.Email == email {
			return user, nil
		}
	}

	return User{}, errors.New("user does not exists")
}

func findUserByUsername(username string) (User, error) {
	for _, user := range users {
		if user.Username == username {
			return user, nil
		}
	}

	return User{}, errors.New("user does not exists")
}

func isBlankPointer(value *string) bool {
	return value == nil || strings.TrimSpace(*value) == ""
}

func validateCreateUser(input RegisterInput) []string {
	var errorsList []string
	if isBlankPointer(input.Username) {
		errorsList = append(errorsList, "username is required")
	}

	if isBlankPointer(input.Email) {
		errorsList = append(errorsList, "email is required")
	}

	if isBlankPointer(input.Password) {
		errorsList = append(errorsList, "password is required")
	}

	if isBlankPointer(input.ConfirmPassword) {
		errorsList = append(errorsList, "confirm password is required")
	}

	if len(errorsList) == 0 && *input.Password != *input.ConfirmPassword {
		errorsList = append(errorsList, "passwords do not match")
	}

	return errorsList
}

func setCORSHeader(w http.ResponseWriter, methods string) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", methods)
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
	w.Header().Set("Content-Type", "application/json")
}

func writeJSON(w http.ResponseWriter, statusCode int, data any) {
	w.WriteHeader(statusCode)
	err := json.NewEncoder(w).Encode(data)
	if err != nil {
		log.Println("invalid json: ", err)
	}
}

func writeError(w http.ResponseWriter, statusCode int, message string) {
	writeJSON(w, statusCode, map[string]any{
		"success": false,
		"data":    nil,
		"errors":  message,
	})
}

func getCurrentUserID(r *http.Request) (int, error) {
	userIDText := r.Header.Get("X-User-ID")

	if userIDText == "" {
		return 0, errors.New("missing user id")
	}

	return strconv.Atoi(userIDText)
}
