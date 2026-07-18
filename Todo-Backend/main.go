package main

import (
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	initDB()
	defer func() {
		if err := DB.Close(); err != nil {
			log.Println("Failed to close DB connection pool:", err)
		}
	}()

	if err != nil {
		log.Fatal("Error loading .env file")
	}
	mux := http.NewServeMux()
	mux.Handle("/tasks", Auth(http.HandlerFunc(tasksHandler)))
	mux.Handle("/tasks/", Auth(http.HandlerFunc(taskByIDHandler)))
	mux.HandleFunc("/login", loginHandler)
	mux.HandleFunc("/register", registerHandler)
	log.Println("Server running on http://localhost:" + os.Getenv("APP_PORT"))
	err = http.ListenAndServe(":"+os.Getenv("APP_PORT"), CORS(mux))
	if err != nil {
		log.Fatal("Failed to start server:", err)
	}
}
