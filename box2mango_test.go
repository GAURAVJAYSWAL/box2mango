package box2mango

import (
	"testing"

	"github.com/joho/godotenv"
)

func TestGetEntpToken(t *testing.T) {
	err := godotenv.Load()
	if err != nil {
		t.Errorf("Error loading .env file")
	}

	token, errStr, err := getEntpToken()
	if err != nil || token == "" {
		t.Errorf("Error getting enterprise token: %v : %v", errStr, err)
	} else {
		t.Logf("Token: %v", token)
	}
}
