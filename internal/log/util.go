package log

import (
	"fmt"
	"log/slog"
	"os"
)

func Fatal(msg string, err error, args ...any) {
	slog.Error(fmt.Sprintf(msg, args), "error", err)
	os.Exit(1)
}