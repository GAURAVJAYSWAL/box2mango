package mango

import (
	"testing"

	"github.com/joho/godotenv"
)

func TestCreateFolderEntry(t *testing.T) {
	err := godotenv.Load("../.env")
	if err != nil {
		t.Errorf("Error loading .env file")
	}
	err1 := createFolderEntry("Test Folder")
	if err1 != nil {
		t.Errorf(err1.Error())
	}
}
