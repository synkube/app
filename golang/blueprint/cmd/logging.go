package cmd

import (
	"context"
	"log/slog"

	"github.com/synkube/app/core/common"
)

func InitLogging() *slog.Logger {

	logger := common.NewLogger()
	logger.InfoContext(context.Background(), "Logger initialized")
	slog.SetDefault(logger)
	return logger
}
