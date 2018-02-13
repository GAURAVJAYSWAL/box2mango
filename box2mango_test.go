package main

import (
	"fmt"
	"testing"

	"github.com/joho/godotenv"
	"github.com/siddhartham/box2mango/lib"
)

func TestCreateUserFolder(t *testing.T) {
	err := godotenv.Load()
	if err != nil {
		lib.Err("main", fmt.Errorf("Error loading .env file"))
	}

	// users, _ := b2m.box.GetEntpUsers()
	// for _, user := range users.Entries {
	// 	b2m.mango.CreateUserBoxFolderEntry(user.Login, user.Name, user.ID)
	// 	newFolderName := fmt.Sprintf("%v%v", user.Name, os.Getenv("BOXFOLDERSUFFIX"))
	// 	b2m.downloadFolderRecursively("0", user.ID, "0", newFolderName)
	// }
}
