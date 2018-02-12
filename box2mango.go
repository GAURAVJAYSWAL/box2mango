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

func (b2m *Box2Mango) downloadFolderRecursively(folderBoxID string, userBoxID string, parentFolderBoxID string, folderName string) {
	fmt.Printf("\nDownloading folder id %v for user id %v\n", folderBoxID, userBoxID)

	if folderBoxID != "0" {
		folderID, _ := b2m.mango.CreateBoxChildFolderEntry(userBoxID, parentFolderBoxID, folderName, folderBoxID)
		fmt.Printf("Created folder entry %v", folderID)
	}

	items, err1 := b2m.box.GetFolderItems(folderBoxID, userBoxID)
	if err1 != nil {
		fmt.Printf("Error getting folder items: %v", err1)
	}

	for _, item := range items.Entries {
		if item.Type == "file" {
			b2m.box.DownloadFile(item.ID, userBoxID)
		} else {
			b2m.downloadFolderRecursively(item.ID, userBoxID, folderBoxID, item.Name)
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
