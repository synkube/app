package main

import (
	"log"
	"os"

	"github.com/synkube/app/blueprint/cmd"
	"github.com/synkube/app/core/common"
)

func main() {
	buildInfo := common.BuildInfo()
	err := cmd.Start(os.Args, buildInfo)
	if err != nil {
		log.Fatalf("Failed to run the application: %v", err)
	}
}
