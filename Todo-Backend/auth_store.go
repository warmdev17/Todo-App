package main

import (
	"database/sql"
	"errors"
	"log"
	"net/http"

	"github.com/google/uuid"
)

func fetchUser(query string, args ...any) (User, error) {
	var user User

	err := DB.Get(&user, query, args...)

	if err != nil {
		if err == sql.ErrNoRows {
			return User{}, errors.New("user does not exist")
		}

		return User{}, errors.New("user does not exist")
	}

	return user, nil
}

func findUserByEmail(email string) (User, error) {
	return fetchUser("SELECT id, username, email, hash_password FROM users WHERE email = $1", email)
}

func findUserByUsername(username string) (User, error) {
	return fetchUser("SELECT id, username, email, hash_password FROM users WHERE username = $1", username)
}

func getUserByID(id uuid.UUID) (User, error) {
	return fetchUser("SELECT id, username, email, hash_password FROM users WHERE id = $1", id)
}

func getCurrentUser(r *http.Request) (User, int, error) {
	userIDAny := r.Context().Value(userIDKey)

	if userIDAny == nil {
		return User{}, http.StatusUnauthorized, errors.New("unauthorized")
	}

	userID, ok := userIDAny.(uuid.UUID)
	if !ok {
		return User{}, http.StatusUnauthorized, errors.New("invalid user id in context")
	}

	user, err := getUserByID(userID)
	if err != nil {
		return User{}, http.StatusUnauthorized, errors.New("user not found")
	}

	return user, http.StatusOK, nil
}

func createNewUser(user User) (uuid.UUID, error) {
	query := `
		INSERT INTO users (username, email, hash_password) 
		VALUES (:username, :email, :hash_password)
		RETURNING id`

	rows, err := DB.NamedQuery(query, user)

	if err != nil {
		return uuid.Nil, err
	}

	defer func() {
		err := rows.Close()

		if err != nil {
			log.Println("Failed to close db connection pool")
		}
	}()

	var newID uuid.UUID

	if rows.Next() {
		err := rows.Scan(&newID)

		if err != nil {
			return uuid.Nil, err
		}
	}
	return newID, nil
}
