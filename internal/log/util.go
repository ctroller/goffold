package log

import (
	"log/slog"
	"os"
)

func Fatal(msg string, err error) {
	slog.Error(msg, "error", err)
	os.Exit(1)
}