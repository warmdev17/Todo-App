package main

import (
	"encoding/json"
	"log"
	"net/http"
)

func setJSONHeader(w http.ResponseWriter) {
	w.Header().Set("Content-Type", "application/json")
}

func writeJSON(w http.ResponseWriter, statusCode int, data any) {
	setJSONHeader(w)

	w.WriteHeader(statusCode)
	err := json.NewEncoder(w).Encode(data)
	if err != nil {
		log.Println("invalid json: ", err)
	}
}

func writeError(w http.ResponseWriter, statusCode int, message string) {
	writeJSON(w, statusCode, map[string]any{
		"success": false,
		"data":    nil,
		"errors":  message,
	})
}

func writeValidationError(w http.ResponseWriter, statusCode int, errors []string) {
	writeJSON(w, statusCode, map[string]any{
		"success": false,
		"data":    nil,
		"errors":  errors,
	})
}
