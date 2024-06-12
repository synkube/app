package main

import (
	"log"
	"log/slog"
	"os"

	"github.com/synkube/app/blueprint/cmd"
	"github.com/synkube/app/core/common"
)

var Logger *slog.Logger

var (
	version = common.DefaultVersion
	date    = common.DefaultDate
)

func main() {
	buildInfo := common.NewBuildInfo(version, "blueprint", date)
	Logger = cmd.InitLogging()

	err := cmd.Start(os.Args, buildInfo)
	if err != nil {
		log.Fatalf("Failed to run the application: %v", err)
	}
}
