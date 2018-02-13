package main

import (
	"fmt"
	"runtime"

	"github.com/CrowdSurge/banner"
	"github.com/fatih/color"
	"github.com/joho/godotenv"
	"github.com/siddhartham/box2mango/boxtools"
	"github.com/siddhartham/box2mango/lib"
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
		lib.Err("main", fmt.Errorf("Error loading .env file"))
	}

}

func (b2m *Box2Mango) downloadFolderRecursively(folderBoxID string, userBoxID string, parentFolderBoxID string, folderName string) {
	var folderID int64
	if folderBoxID != "0" {
		folderID, _ = b2m.mango.CreateBoxChildFolderEntry(userBoxID, parentFolderBoxID, folderName, folderBoxID)
	}
	items, _ := b2m.box.GetFolderItems(folderBoxID, userBoxID)
	for _, item := range items.Entries {
		if item.Type == "file" {
			sanPath, err := b2m.box.DownloadFile(item.ID, userBoxID)
			if err == nil {
				b2m.mango.CreateBoxChildFileEntry(folderID, item.Name, sanPath, item.ID)
			}
		} else {
			b2m.downloadFolderRecursively(item.ID, userBoxID, folderBoxID, item.Name)
		}
	}
}

func (b2m *Box2Mango) createAllUserFolders() {
	users, _ := b2m.box.GetEntpUsers()
	for _, user := range users.Entries {
		b2m.mango.CreateUserBoxFolderEntry(user.Login, user.Name, user.ID)
	}
}
