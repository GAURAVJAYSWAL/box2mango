package box2mango

import (
	"fmt"
	"log"
	"runtime"

	"github.com/CrowdSurge/banner"
	"github.com/fatih/color"
	"github.com/joho/godotenv"
	"github.com/siddhartham/box2mango/boxtools"
	"github.com/siddhartham/box2mango/mangotools"
)

type Box2Mango struct {
	box   boxtools.BoxService
	mango mangotools.MangoService
}

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())

	banner.Print("  box2mango  ")
	color.Red("\n- by Siddhartha Mukherjee <mukherjee.siddhartha@gmail.com>")
	color.Yellow("\nUsage: box2mango\n")

	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

}

func (b2m *Box2Mango) downloadFolderRecursively(folderID string, userID string) {
	fmt.Printf("\nDownloading folder id %v for user id %v\n", folderID, userID)
	items, err := b2m.box.GetFolderItems(folderID, userID)
	if err != nil {
		fmt.Printf("Error getting folder items: %v", err)
	}

	for _, item := range items.Entries {
		if item.Type == "folder" {
			b2m.downloadFolderRecursively(item.ID, userID)
		}
	}
}

func (b2m *Box2Mango) createAllUserFolders() {
	users, err := b2m.box.GetEntpUsers()
	if err != nil {
		fmt.Printf("Error getting enterprise users: %v", err)
	}
	for _, user := range users.Entries {
		fmt.Printf("\nCreating folder for %v\n", user.Login)
		b2m.mango.CreateUserBoxFolderEntry(user.Login, user.Name, user.ID)
	}
}
