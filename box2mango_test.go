package box2mango

import (
	"testing"

	"github.com/joho/godotenv"
)

func TestCreateAllUserFolders(t *testing.T) {
	err := godotenv.Load(".env")
	if err != nil {
		t.Errorf("Error loading .env file")
	}
	createAllUserFolders()
}

func TestDownloadFolders(t *testing.T) {
	err := godotenv.Load(".env")
	if err != nil {
		t.Errorf("Error loading .env file")
	}
}
