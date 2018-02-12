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

	ms := MangoService{}

	err1 := ms.CreateUserBoxFolderEntry("siddhartham@mangospring.com", "Test1", "12345678")
	if err1 != nil {
		t.Errorf(err1.Error())
	}

	err2 := ms.CreateUserBoxFolderEntry("max@mangospring.com", "Test2", "12345679")
	if err2 != nil {
		t.Errorf(err2.Error())
	}
}
