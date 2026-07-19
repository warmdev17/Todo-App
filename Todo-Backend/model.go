package main

import "github.com/google/uuid"

type Task struct {
	ID        int       `json:"id" db:"id"`
	UserID    uuid.UUID `json:"userId" db:"user_id"`
	Title     string    `json:"title" db:"title"`
	Completed bool      `json:"completed" db:"completed"`
}

type User struct {
	ID           uuid.UUID `json:"id" db:"id"`
	Email        string    `json:"email" db:"email"`
	Username     string    `json:"username" db:"username"`
	HashPassword string    `json:"-" db:"hash_password"`
}

type AuthUser struct {
	ID       uuid.UUID `json:"id"`
	Email    string    `json:"email"`
	Username string    `json:"username"`
}

type RegisterInput struct {
	Email           *string `json:"email"`
	Username        *string `json:"username"`
	Password        *string `json:"password"`
	ConfirmPassword *string `json:"confirmPassword"`
}

type UpdateTaskInput struct {
	Completed *bool   `json:"completed" db:"completed"`
	Title     *string `json:"title" db:"title"`
}
