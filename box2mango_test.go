package box2mango

import (
	"testing"

	"github.com/joho/godotenv"
	"github.com/siddhartham/box2mango/boxtools"
	"github.com/siddhartham/box2mango/mangotools"
)

var (
	b2m = Box2Mango{
		box:   boxtools.BoxService{},
		mango: mangotools.MangoService{},
	}
)

func TestCreateAllUserFolders(t *testing.T) {
	err := godotenv.Load(".env")
	if err != nil {
		t.Errorf("Error loading .env file")
	}
	b2m.createAllUserFolders()
}

func TestDownloadFolders(t *testing.T) {
	err := godotenv.Load(".env")
	if err != nil {
		t.Errorf("Error loading .env file")
	}
	b2m.downloadFolderRecursively("0", "272313645")
}
