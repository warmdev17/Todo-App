package main

import (
	"log"
	"testing"

	"github.com/google/uuid"
	"github.com/joho/godotenv"
)

func TestGetTasksByUserID(t *testing.T) {
	err := godotenv.Load()
	initDB()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	defer func() {
		if err := DB.Close(); err != nil {
			log.Println("Failed to close DB connection pool:", err)
		}
	}()

	userID, err := uuid.Parse("7da47fd2-d479-45cf-9ae0-8f0574be276b")

	if err != nil {
		t.Fatalf("failed to parse uuid: %v", err)
	}
	tasks, err := getTasksByUserID(userID)

	if err != nil {
		t.Fatalf("getTasksByUserID returned error: %v", err)
	}

	t.Logf("tasks: %v", tasks)
}
