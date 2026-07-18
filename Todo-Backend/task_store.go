package main

import (
	"database/sql"
	"errors"

	"github.com/google/uuid"
)

var tasks = []Task{}

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

func getTasksByUserID(userID uuid.UUID) ([]Task, error) {
	var tasks []Task

	query := `SELECT id, title, completed FROM tasks WHERE user_id = $1`

	err := DB.Select(&tasks, query, userID)

	if err != nil {
		return nil, err
	}

	return tasks, nil
}
