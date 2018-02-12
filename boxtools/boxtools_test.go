package boxtools

import (
	"fmt"
	"testing"

	"github.com/joho/godotenv"
)

func TestTheWholeThing(t *testing.T) {
	err := godotenv.Load("../.env")
	if err != nil {
		t.Errorf("Error loading .env file")
	}
	bs := BoxService{}

	users, err1 := bs.GetEntpUsers()
	fmt.Printf("\nGot users %v", users.TotalCount)
	if err1 != nil {
		t.Errorf("Error getting enterprise users: %v", err1)
	}

	items, err2 := bs.GetFolderItems("0", users.Entries[2].ID)
	fmt.Printf("\nGot items %v\n", items.TotalCount)
	if err2 != nil {
		t.Errorf("Error getting folder items: %v", err2)
	}
}
