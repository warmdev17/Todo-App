package main

import (
	"errors"
	"net/http"
	"strconv"
)

var users = []User{
	{1, "hoangmaiphuongtin@gmail.com", "maiphuong", "$2a$10$08/9rz35z3xJ0X0mqikMf.1cgPRiC6Vhi6A7W4dbRixDAtViOKBd."},
	{2, "warmdevofficial@gmail.com", "warmdev", "$2a$10$BVWV36D.NghpfB9O5gd4muhxzsiXXxTaAnGt6tFA/gkf2xcDkfoN6"},
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

	return User{}, errors.New("user does not exist")
}

func findUserByUsername(username string) (User, error) {
	for _, user := range users {
		if user.Username == username {
			return user, nil
		}
	}

	return User{}, errors.New("user does not exist")
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
