package util

import "github.com/google/uuid"

// Generates a random UUID string
func GenerateStringUUID() string {
	id, _ := uuid.NewRandom()
	return id.String()
}
