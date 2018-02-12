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

func downloadFolderRecursively(boxUserId string, boxFolderId string) {
	// bs := boxtools.BoxService{}
}

func createAllUserFolders() {
	bs := boxtools.BoxService{}
	users, err := bs.GetEntpUsers()
	if err != nil {
		fmt.Printf("Error getting enterprise users: %v", err)
	}
	for _, user := range users.Entries {
		fmt.Printf("\nCreating folder for %v\n", user.Login)
		mangotools.CreateUserBoxFolderEntry(user.Login, user.Name, user.ID)
	}
}
