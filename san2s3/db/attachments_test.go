package db

import (
	"fmt"
	"testing"

	"github.com/joho/godotenv"
)

func TestGetFileFromDb(t *testing.T) {
	err := godotenv.Load()
	if err != nil {
		t.Errorf("Error loading .env file")
	}
	files, err := getFileFromDb()
	if err == nil {
		for _, element := range files {
			fmt.Printf("responses %s ", element.SanStorageUrl)
		}
	}
}
