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
)

type Task struct {
	ID        int    `json:"id"`
	Title     string `json:"title"`
	Completed bool   `json:"completed"`
}

type User struct {
	ID       int    `json:"id"`
	Email    string `json:"email"`
	Username string `json:"username"`
	Password string `json:"-"`
}

type AuthUser struct {
	ID    int    `json:"id"`
	Email string `json:"email"`
}

type RegisterInput struct {
	Email           *string `json:"email"`
	Username        *string `json:"username"`
	Password        *string `json:"password"`
	ConfirmPassword *string `json:"confirmPassword"`
}

var tasks = []Task{
	{ID: 0, Title: "Learn Go net/http", Completed: false},
	{ID: 1, Title: "Build TODO REST API app", Completed: false},
}

var users = []User{
	{0, "warmdevofficial@gmail.com", "warmdev", "Warmdev17@todo"},
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
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
	w.Header().Set("Content-Type", "application/json")

	if r.Method == http.MethodOptions {
		w.WriteHeader(http.StatusNoContent)
		return
	}

	if r.Method == http.MethodGet {
		log.Printf("method = %s, path = %s", r.Method, r.URL.Path)
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

		log.Printf("method = %s, path = %s", r.Method, r.URL.Path)
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
			ID:        nextTaskID(),
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

func taskByIDHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "DELETE, POST, GET, PATCH, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
	w.Header().Set("Content-Type", "application/json")

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
					w.WriteHeader(http.StatusOK)
					if err := json.NewEncoder(w).Encode(map[string]any{
						"success": true,
						"data":    task,
					}); err != nil {
						http.Error(w, "Invalid JSON", http.StatusBadRequest)
						return
					}
					tasks = append(tasks[:index], tasks[index+1:]...)
				}
			}
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

			for index, task := range tasks {
				if task.ID == id {
					if input.Title != nil {
						tasks[index].Title = strings.TrimSpace(*input.Title)
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
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "GET, POST")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
	w.Header().Set("Content-Type", "application/json")
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
			log.Println("Email login")
			user, err = findUserByEmail(*LoginInput.Email)
			if err != nil {
				http.Error(w, err.Error(), http.StatusUnauthorized)
				return
			}
			if user.Password == *LoginInput.Password {
				w.WriteHeader(http.StatusCreated)
				err := json.NewEncoder(w).Encode(map[string]any{
					"success": true,
					"data": map[string]any{
						"token":    "fake-token-1",
						"username": user.Username,
					},
				})
				if err != nil {
					http.Error(w, err.Error(), http.StatusBadRequest)
					return
				}
			} else {
				w.WriteHeader(http.StatusUnauthorized)
				err := json.NewEncoder(w).Encode(map[string]any{
					"success": false,
					"data":    nil,
					"error":   "Invalid email or password",
				})
				if err != nil {
					http.Error(w, err.Error(), http.StatusBadRequest)
					return
				}
			}
		} else {
			log.Println("Username login")
			user, err = findUserByUsername(*LoginInput.Username)
			if err != nil {
				http.Error(w, err.Error(), http.StatusUnauthorized)
				return
			}

			if user.Password == *LoginInput.Password {
				w.WriteHeader(http.StatusOK)
				err := json.NewEncoder(w).Encode(map[string]any{
					"success": true,
					"data": map[string]any{
						"token":    "fake-token-1",
						"username": user.Username,
					},
				})
				if err != nil {
					http.Error(w, err.Error(), http.StatusBadRequest)
					return
				}
			} else {
				w.WriteHeader(http.StatusUnauthorized)
				err := json.NewEncoder(w).Encode(map[string]any{
					"success": false,
					"data":    nil,
					"error":   "Invalid username or password",
				})
				if err != nil {
					http.Error(w, err.Error(), http.StatusBadRequest)
					return
				}

			}
		}
	}
}

func registerHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
	w.Header().Set("Content-Type", "application/json")

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

		errs := validateCreateUser(input)

		if len(errs) > 0 {
			w.WriteHeader(http.StatusBadRequest)
			if err := json.NewEncoder(w).Encode(map[string]any{
				"success": false,
				"data":    nil,
				"errors":  errs,
			}); err != nil {
				http.Error(w, "Invalid JSON Body", http.StatusBadRequest)
				return
			}
			return
		}

		newUser := User{
			ID:       nextUserID(),
			Username: *input.Username,
			Email:    *input.Email,
			Password: *input.Password,
		}
		users = append(users, newUser)
		w.WriteHeader(http.StatusCreated)
		err := json.NewEncoder(w).Encode(map[string]any{
			"success": true,
			"data":    newUser,
		})
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
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
	max := tasks[0].ID

	for _, task := range tasks {
		if task.ID > max {
			max = task.ID
		}
	}

	return max + 1
}

func nextUserID() int {
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
