package main

import (
	"log"

	"github.com/synkube/app/blueprint/cmd"
)

func main() {
	err := cmd.RunApp()
	if err != nil {
		log.Fatalf("Failed to run the application: %v", err)
	}
}
