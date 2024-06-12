package main

import (
	"log"
	"log/slog"
	"os"

	"github.com/synkube/app/blueprint/cmd"
	"github.com/synkube/app/core/common"
)

var Logger *slog.Logger

func main() {
	buildInfo := common.BuildInfo()
	Logger = cmd.InitLogging()

	err := cmd.Start(os.Args, buildInfo)
	if err != nil {
		log.Fatalf("Failed to run the application: %v", err)
	}
}
