package main

import (
	"encoding/json"
	"net/http"
	"strings"

	"golang.org/x/crypto/bcrypt"
)

func loginHandler(w http.ResponseWriter, r *http.Request) {
	setCORSHeader(w, "POST, OPTIONS")

	var loginInput struct {
		Email    *string `json:"email"`
		Password *string `json:"password"`
		Username *string `json:"username"`
	}

	if r.Method == http.MethodOptions {
		w.WriteHeader(http.StatusNoContent)
		return
	}

	if r.Method == http.MethodPost {
		err := json.NewDecoder(r.Body).Decode(&loginInput)
		if err != nil {
			writeError(w, http.StatusBadRequest, "Invalid JSON")
			return
		}

		if loginInput.Email == nil &&
			loginInput.Username == nil &&
			loginInput.Password == nil {
			writeError(w, http.StatusBadRequest, "request body is required")
			return
		}

		var user User

		if isBlankPointer(loginInput.Password) {
			writeError(w, http.StatusBadRequest, "password is required")
			return
		}

		hasEmail := !isBlankPointer(loginInput.Email)
		hasUsername := !isBlankPointer(loginInput.Username)

		if hasEmail == hasUsername {
			writeError(w, http.StatusBadRequest, "use either username or email")
			return
		}
		if hasEmail {
			trimmedEmail := strings.TrimSpace(*loginInput.Email)
			user, err = findUserByEmail(trimmedEmail)
			if err != nil {
				writeError(w, http.StatusUnauthorized, "Invalid credentials")
				return
			}
			err := bcrypt.CompareHashAndPassword([]byte(user.HashPassword), []byte(*loginInput.Password))
			if err != nil {
				writeError(w, http.StatusUnauthorized, "Invalid credentials")
				return
			} else {
				tokenString, err := generateToken(user.ID, user.Username)

				if err != nil {
					writeError(w, http.StatusInternalServerError, "Failed to generate token")
					return
				}

				writeJSON(w, http.StatusOK, map[string]any{
					"success": true,
					"data": map[string]any{
						"token":    tokenString,
						"userId":   user.ID,
						"username": user.Username,
					},
				})
				return
			}
		} else {
			trimmedUsername := strings.TrimSpace(*loginInput.Username)
			user, err = findUserByUsername(trimmedUsername)
			if err != nil {
				writeError(w, http.StatusUnauthorized, "Invalid credentials")
				return
			}

			err := bcrypt.CompareHashAndPassword([]byte(user.HashPassword), []byte(*loginInput.Password))
			if err != nil {
				writeError(w, http.StatusUnauthorized, "Invalid credentials")
				return
			} else {
				tokenString, err := generateToken(user.ID, user.Username)

				if err != nil {
					writeError(w, http.StatusInternalServerError, "Failed to generate token")
					return
				}
				writeJSON(w, http.StatusOK, map[string]any{
					"success": true,
					"data": map[string]any{
						"token":    tokenString,
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

		newUser := User{Email: trimmedEmail, Username: trimmedUsername, HashPassword: string(hashPassword)}

		userID, err := createNewUser(newUser)
		if err != nil {
			writeError(w, http.StatusConflict, err.Error())
			return
		}
		writeJSON(w, http.StatusCreated, map[string]any{
			"success": true,
			"data": AuthUser{
				ID:       userID,
				Email:    newUser.Email,
				Username: newUser.Username,
			},
		})
		return
	}

	writeError(w, http.StatusMethodNotAllowed, "Method not allowed")
}
