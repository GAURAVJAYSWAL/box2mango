package mangotools

import (
	"testing"

	"github.com/joho/godotenv"
)

func TestCreateFolderEntry(t *testing.T) {
	err := godotenv.Load("../.env")
	if err != nil {
		t.Errorf("Error loading .env file")
	}
	err1 := CreateUserBoxFolderEntry("siddhartham@mangospring.com", "Siddhartha Mukherjee")
	if err1 != nil {
		t.Errorf(err1.Error())
	}

	err2 := CreateUserBoxFolderEntry("max@mangospring.com", "Max")
	if err2 != nil {
		t.Errorf(err2.Error())
	}
}
