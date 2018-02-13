package main

import (
	"fmt"
	"os"
	"runtime"

	"github.com/CrowdSurge/banner"
	"github.com/fatih/color"
	"github.com/joho/godotenv"
	"github.com/siddhartham/box"
	"github.com/siddhartham/box2mango/boxtools"
	"github.com/siddhartham/box2mango/lib"
	"github.com/siddhartham/box2mango/mangotools"
)

type Box2Mango struct {
	box   boxtools.BoxService
	mango mangotools.MangoService
}

var (
	b2m = Box2Mango{
		box:   boxtools.BoxService{},
		mango: mangotools.MangoService{},
	}
)

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())

	banner.Print("  box2mango  ")
	color.Red("\n- by Siddhartha Mukherjee <mukherjee.siddhartha@gmail.com>")
	color.Yellow("\nUsage: box2mango\n")

	err := godotenv.Load()
	if err != nil {
		lib.Err("main", fmt.Errorf("Error loading .env file"))
	}

	users, _ := b2m.box.GetEntpUsers()
	for _, user := range users.Entries {
		b2m.mango.CreateUserBoxFolderEntry(user.Login, user.Name, user.ID)
		newFolderName := fmt.Sprintf("%v%v", user.Name, os.Getenv("BOXFOLDERSUFFIX"))
		b2m.downloadFolderRecursively("0", user.ID, "0", newFolderName)
	}

}

func (b2m *Box2Mango) downloadFolderRecursively(folderBoxID string, userBoxID string, parentFolderBoxID string, folderName string) {
	var folderID int64
	var roleID int
	var collabs *box.Collaborations
	if folderBoxID != "0" {
		folderID, _ = b2m.mango.CreateBoxChildFolderEntry(userBoxID, parentFolderBoxID, folderName, folderBoxID)
		collabs, _ = b2m.box.GetFolderCollaborations(folderBoxID, userBoxID)
		for _, clb := range collabs.Entries {
			if clb.AccessibleBy.Type == "user" {
				switch role := clb.Role; role {
				case "co-owner":
					roleID = 5
				case "editor":
					roleID = 6
				default:
					roleID = 7
				}
				b2m.mango.CreateFollowListEntry(folderID, clb.AccessibleBy.Login, roleID)
			}
		}
	}
	items, _ := b2m.box.GetFolderItems(folderBoxID, userBoxID)
	for _, item := range items.Entries {
		if item.Type == "file" {
			sanPath, err := b2m.box.DownloadFile(item.ID, userBoxID)
			if err == nil {
				fileID, _ := b2m.mango.CreateBoxChildFileEntry(folderID, item.Name, sanPath, item.ID)
				for _, clb := range collabs.Entries {
					if clb.AccessibleBy.Type == "user" {
						switch role := clb.Role; role {
						case "co-owner":
							roleID = 1
						case "editor":
							roleID = 2
						default:
							roleID = 3
						}
						b2m.mango.CreateFollowListEntry(fileID, clb.AccessibleBy.Login, roleID)
					}
				}
			}
		} else {
			b2m.downloadFolderRecursively(item.ID, userBoxID, folderBoxID, item.Name)
		}
	}
}
