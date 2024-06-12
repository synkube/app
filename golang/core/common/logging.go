package common

import (
	"fmt"
	"log/slog"
	"os"
	"strings"
	"time"
)

var Logger *slog.Logger

func NewLogger() *slog.Logger {
	// Define the handler with custom options
	baseHandler := slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		// AddSource: true,
		ReplaceAttr: func(groups []string, a slog.Attr) slog.Attr {

			switch a.Key {
			case slog.SourceKey:
				src := a.Value.Any().(*slog.Source)
				shortPath := ""
				fullPath := src.File
				seps := strings.Split(fullPath, "/")
				shortPath += seps[len(seps)-1]
				shortPath += fmt.Sprintf(":%d", src.Line)
				a.Value = slog.StringValue(shortPath)
			case slog.TimeKey:
				return slog.String("time", time.Now().Format("2006-01-02 15:04:05"))
			default:
				return a
			}
			return a
		},
	})
	// Wrap the base handler with the custom handler
	// NOTE !!!! Not currently used
	// customHandler := NewCustomHandler(baseHandler)

	// Create a new logger with the custom handler and assign it to the global Logger variable
	Logger = slog.New(baseHandler)
	return Logger
}
