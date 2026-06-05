package main

import "errors"

var tasks = []Task{
	{ID: 1, UserID: 1, Title: "Learn Go net/http", Completed: false},
	{ID: 2, UserID: 1, Title: "Learn Gin framework", Completed: false},
	{ID: 3, UserID: 2, Title: "Build TODO REST API app", Completed: false},
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
