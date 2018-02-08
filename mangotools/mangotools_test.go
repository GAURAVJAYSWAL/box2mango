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
	err1 := CreateUserBoxFolderEntry("siddhartham@mangospring.com")
	if err1 != nil {
		t.Errorf(err1.Error())
	}
}
