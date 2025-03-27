package main

import (
	"fmt"
	"os"

	"github.com/jcfug8/daylear/server/daylearctl/commands/start"
)

func main() {
	var err error
	if len(os.Args) > 1 {
		switch os.Args[1] {
		case "users":
			// err = users.Users(os.Args[2:]...)
		default:
			fmt.Printf("Unknown command: %s\n", os.Args[1])
			os.Exit(1)
		}
	} else {
		err = start.Start()
	}

	if err != nil {
		fmt.Printf("Error: %s\n", err)
		os.Exit(1)
	}
}
