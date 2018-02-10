package boxtools

import (
	"testing"

	"github.com/joho/godotenv"
)

func TestTheWholeThing(t *testing.T) {
	err := godotenv.Load("../.env")
	if err != nil {
		t.Errorf("Error loading .env file")
	}
	bs := BoxService{}
	bs.GetEntpUsers()
}
