package boxtools

import (
	"testing"

	"github.com/joho/godotenv"
)

func TestGetEntpToken(t *testing.T) {
	err := godotenv.Load("../.env")
	if err != nil {
		t.Errorf("Error loading .env file")
	}
	tok, errStr, err := GetEntpToken()
	if err != nil || tok == nil {
		t.Errorf("Error getting enterprise token: %v : %v", errStr, err)
	}

	GetEntpUsers(tok)
}
