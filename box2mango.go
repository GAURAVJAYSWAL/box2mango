package box2mango

import (
	"log"
	"runtime"

	"github.com/CrowdSurge/banner"
	"github.com/fatih/color"
	"github.com/joho/godotenv"
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
