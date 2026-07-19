package main

import (
	"database/sql"
	"errors"
	"fmt"
	"strings"

	"github.com/google/uuid"
)

func getTaskByIDAndUserID(taskID int, userID uuid.UUID) (Task, error) {
	var task Task
	query := `
		SELECT id, title, completed FROM tasks WHERE id = $1 AND user_id = $2
	`

	err := DB.Get(&task, query, taskID, userID)
	if err != nil {
		return Task{}, err
	}

	return task, nil
}

func createNewTask(userID uuid.UUID, title string) (Task, error) {
	var newTask Task
	query := `
		INSERT INTO tasks (title, user_id)
		VALUES ($1, $2) RETURNING *
	`

	err := DB.Get(&newTask, query, title, userID)

	if err != nil {
		return Task{}, err
	}

	return newTask, nil
}
func deleteTaskByIDAndUserID(taskID int, userID uuid.UUID) (Task, error) {
	var deletedTask Task

	query := `
		DELETE FROM tasks
		WHERE id = $1 AND user_id = $2 
		RETURNING *
	`

	err := DB.Get(&deletedTask, query, taskID, userID)

	if err != nil {
		if err == sql.ErrNoRows {
			return Task{}, errors.New("Task not found")
		}

		return Task{}, err
	}

	return deletedTask, nil
}

func updateTaskByIDAndUserID(taskID int, userID uuid.UUID, input UpdateTaskInput) (Task, error) {
	var udpatedTask Task

	setClauses := []string{}
	args := []any{}
	argID := 1

	if input.Title != nil {
		setClauses = append(setClauses, fmt.Sprintf("title = $%d", argID))
		args = append(args, *input.Title)
		argID++
	}

	if input.Completed != nil {
		setClauses = append(setClauses, fmt.Sprintf("completed = $%d", argID))
		args = append(args, *input.Completed)
		argID++
	}

	query := `UPDATE tasks SET ` + strings.Join(setClauses, ", ")

	query += fmt.Sprintf(` WHERE id = $%d AND user_id = $%d RETURNING *`, argID, argID+1)
	args = append(args, taskID, userID)

	err := DB.Get(&udpatedTask, query, args...)

	if err != nil {
		if err == sql.ErrNoRows {
			return Task{}, errors.New("Task not found")
		}
		return Task{}, err
	}

	return udpatedTask, nil
}

func getTasksByUserID(userID uuid.UUID) ([]Task, error) {
	var tasks []Task

	query := `SELECT id, title, completed FROM tasks WHERE user_id = $1`

	err := DB.Select(&tasks, query, userID)

	if err != nil {
		return nil, err
	}

	return tasks, nil
}
