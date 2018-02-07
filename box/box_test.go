package box

import (
	"testing"

	"github.com/joho/godotenv"
)

func TestGetEntpToken(t *testing.T) {
	err := godotenv.Load("../.env")
	if err != nil {
		t.Errorf("Error loading .env file")
	}
	tok, errStr, err := getEntpToken()
	if err != nil || tok == nil {
		t.Errorf("Error getting enterprise token: %v : %v", errStr, err)
	}

	getEntpUsers(tok)
}
