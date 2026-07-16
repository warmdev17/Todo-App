package main

import (
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

func getTaskByIDAndUserID(id int, userID uuid.UUID) (Task, error) {
	for _, task := range tasks {
		if task.ID == id && task.UserID == userID {
			return task, nil
		}
	}
	return Task{}, errors.New("task not found")
}

func getTaskIndexByIDAndUserID(id int, userID uuid.UUID) (int, error) {
	for index, task := range tasks {
		if task.ID == id && task.UserID == userID {
			return index, nil
		}
	}

	return -1, errors.New("task not found")
}
