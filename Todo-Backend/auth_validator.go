package main

import "strings"

func isBlankPointer(value *string) bool {
	return value == nil || strings.TrimSpace(*value) == ""
}

func validateCreateUser(input RegisterInput) []string {
	var errorsList []string
	if isBlankPointer(input.Username) {
		errorsList = append(errorsList, "username is required")
	}

	if isBlankPointer(input.Email) {
		errorsList = append(errorsList, "email is required")
	}

	if isBlankPointer(input.Password) {
		errorsList = append(errorsList, "password is required")
	}

	if isBlankPointer(input.ConfirmPassword) {
		errorsList = append(errorsList, "confirm password is required")
	}

	if len(errorsList) == 0 && *input.Password != *input.ConfirmPassword {
		errorsList = append(errorsList, "passwords do not match")
	}

	return errorsList
}
