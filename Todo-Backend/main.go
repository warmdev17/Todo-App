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

var users = []User{
	{1, "hoangmaiphuongtin@gmail.com", "maiphuong", "$2a$10$08/9rz35z3xJ0X0mqikMf.1cgPRiC6Vhi6A7W4dbRixDAtViOKBd."},
	{2, "warmdevofficial@gmail.com", "warmdev", "$2a$10$BVWV36D.NghpfB9O5gd4muhxzsiXXxTaAnGt6tFA/gkf2xcDkfoN6"},
}

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

	currentUser, statusCode, err := getCurrentUser(r)
	if r.Method == http.MethodGet {
		if err != nil {
			writeError(w, statusCode, err.Error())
			return
		}

		var userTask []Task

		for _, task := range tasks {
			if task.UserID == currentUser.ID {
				userTask = append(userTask, task)
			}
		}

		writeJSON(w, http.StatusOK, map[string]any{
			"success": true,
			"data":    userTask,
		})
		return
	}

	if r.Method == http.MethodPost {
		var input struct {
			Title *string `json:"title"`
		}

		if err != nil {
			writeError(w, http.StatusUnauthorized, err.Error())
			return
		}

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

		newTask := Task{
			ID:        nextTaskID(),
			Title:     trimmedTitle,
			UserID:    currentUser.ID,
			Completed: false,
		}

		tasks = append(tasks, newTask)
		writeJSON(w, http.StatusCreated, map[string]any{
			"success": true,
			"data":    newTask,
		})
		return
	}

	writeError(w, http.StatusMethodNotAllowed, "Method not allowed")
}

func taskByIDHandler(w http.ResponseWriter, r *http.Request) {
	setCORSHeader(w, "GET, DELETE, PATCH, OPTIONS")

	if r.Method == http.MethodOptions {
		w.WriteHeader(http.StatusNoContent)
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

			index, err := getTaskIndexByIDAndUserID(id, currentUser.ID)
			if err != nil {
				writeError(w, http.StatusNotFound, err.Error())
				return
			}

			deletedTask := tasks[index]
			tasks = append(tasks[:index], tasks[index+1:]...)

			writeJSON(w, http.StatusOK, map[string]any{
				"success": true,
				"data":    deletedTask,
			})
			return
		}
	case http.MethodPatch:
		{
			var input struct {
				Completed *bool   `json:"completed"`
				Title     *string `json:"title"`
			}

			err = json.NewDecoder(r.Body).Decode(&input)
			if err != nil {
				writeError(w, http.StatusBadRequest, "Invalid JSON")
				return
			}

			if input.Title == nil && input.Completed == nil {
				writeError(w, http.StatusBadRequest, "No fields to update")
				return
			}

			id, err := getTaskIDFromPath(r)
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
			}

			index, err := getTaskIndexByIDAndUserID(id, currentUser.ID)
			if err != nil {
				writeError(w, http.StatusNotFound, err.Error())
				return
			}

			if input.Title != nil {
				tasks[index].Title = trimmedTitle
			}

			if input.Completed != nil {
				tasks[index].Completed = *input.Completed
			}

			writeJSON(w, http.StatusOK, map[string]any{
				"success": true,
				"data":    tasks[index],
			})

			return

		}
	default:
		{
			writeError(w, http.StatusMethodNotAllowed, "Method not allowed")
			return
		}
	}
}

func loginHandler(w http.ResponseWriter, r *http.Request) {
	setCORSHeader(w, "POST, OPTIONS")

	var LoginInput struct {
		Email    *string `json:"email"`
		Password *string `json:"password"`
		Username *string `json:"username"`
	}

	if r.Method == http.MethodOptions {
		w.WriteHeader(http.StatusNoContent)
		return
	}

	if r.Method == http.MethodPost {
		err := json.NewDecoder(r.Body).Decode(&LoginInput)
		if err != nil {
			writeError(w, http.StatusBadRequest, "Invalid JSON")
			return
		}

		if LoginInput.Email == nil &&
			LoginInput.Username == nil &&
			LoginInput.Password == nil {
			writeError(w, http.StatusBadRequest, "request body is required")
			return
		}

		var user User

		if isBlankPointer(LoginInput.Password) {
			writeError(w, http.StatusBadRequest, "password is required")
			return
		}

		hasEmail := !isBlankPointer(LoginInput.Email)
		hasUsername := !isBlankPointer(LoginInput.Username)

		if hasEmail == hasUsername {
			writeError(w, http.StatusBadRequest, "use either username or email")
			return
		}
		if hasEmail {
			trimmedEmail := strings.TrimSpace(*LoginInput.Email)
			user, err = findUserByEmail(trimmedEmail)
			if err != nil {
				writeError(w, http.StatusUnauthorized, "Invalid credentials")
				return
			}
			err := bcrypt.CompareHashAndPassword([]byte(user.HashPassword), []byte(*LoginInput.Password))
			if err != nil {
				writeError(w, http.StatusUnauthorized, "Invalid credentials")
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
			trimmedUsername := strings.TrimSpace(*LoginInput.Username)
			user, err = findUserByUsername(trimmedUsername)
			if err != nil {
				writeError(w, http.StatusUnauthorized, "Invalid credentials")
				return
			}

			err := bcrypt.CompareHashAndPassword([]byte(user.HashPassword), []byte(*LoginInput.Password))
			if err != nil {
				writeError(w, http.StatusUnauthorized, "Invalid credentials")
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

	writeError(w, http.StatusMethodNotAllowed, "Method not allowed")
}

func registerHandler(w http.ResponseWriter, r *http.Request) {
	setCORSHeader(w, "POST, OPTIONS")

	if r.Method == http.MethodOptions {
		w.WriteHeader(http.StatusNoContent)
		return
	}

	var input RegisterInput

	if r.Method == http.MethodPost {
		if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
			writeError(w, http.StatusBadRequest, "Invalid JSON")
			return
		}

		if input.Username == nil &&
			input.Email == nil &&
			input.Password == nil &&
			input.ConfirmPassword == nil {
			writeError(w, http.StatusBadRequest, "request body is required")
			return
		}

		errs := validateCreateUser(input)
		if len(errs) > 0 {
			writeValidationError(w, http.StatusBadRequest, errs)
			return
		}

		trimmedEmail := strings.TrimSpace(*input.Email)
		trimmedUsername := strings.TrimSpace(*input.Username)

		if _, err := findUserByEmail(trimmedEmail); err == nil {
			writeError(w, http.StatusBadRequest, "email already exists")
			return
		}

		if _, err := findUserByUsername(trimmedUsername); err == nil {
			writeError(w, http.StatusBadRequest, "username already exists")
			return
		}

		hashPassword, err := bcrypt.GenerateFromPassword([]byte(*input.Password), bcrypt.DefaultCost)
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
		return
	}

	writeError(w, http.StatusMethodNotAllowed, "Method not allowed")
}

func getTaskIDFromPath(r *http.Request) (int, error) {
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
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type, X-User-ID")
}

func setJSONHeader(w http.ResponseWriter) {
	w.Header().Set("Content-Type", "application/json")
}

func writeJSON(w http.ResponseWriter, statusCode int, data any) {
	setJSONHeader(w)

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

func writeValidationError(w http.ResponseWriter, statusCode int, errors []string) {
	writeJSON(w, statusCode, map[string]any{
		"success": false,
		"data":    nil,
		"errors":  errors,
	})
}

func getUserByID(id int) (User, error) {
	for _, user := range users {
		if user.ID == id {
			return user, nil
		}
	}

	return User{}, errors.New("user not found")
}

func getCurrentUser(r *http.Request) (User, int, error) {
	userIDText := r.Header.Get("X-User-ID")
	if userIDText == "" {
		return User{}, http.StatusUnauthorized, errors.New("missing user id")
	}

	userID, err := strconv.Atoi(userIDText)
	if err != nil {
		return User{}, http.StatusBadRequest, errors.New("invalid user id")
	}

	user, err := getUserByID(userID)
	if err != nil {
		return User{}, http.StatusUnauthorized, errors.New("unauthorized")
	}

	return user, http.StatusOK, nil
}

func getTaskByIDAndUserID(id int, userID int) (Task, error) {
	for _, task := range tasks {
		if task.ID == id && task.UserID == userID {
			return task, nil
		}
	}
	return Task{}, errors.New("task not found")
}

func getTaskIndexByIDAndUserID(id int, userID int) (int, error) {
	for index, task := range tasks {
		if task.ID == id && task.UserID == userID {
			return index, nil
		}
	}

	return -1, errors.New("task not found")
}
