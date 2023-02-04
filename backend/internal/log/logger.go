package log

import (
	"log"
	"os"
	"strings"

	"github.com/eko/authz/backend/configs"
	"golang.org/x/exp/slog"
)

var (
	levels = map[string]slog.Level{
		slog.LevelDebug.String(): slog.LevelDebug,
		slog.LevelInfo.String():  slog.LevelInfo,
		slog.LevelWarn.String():  slog.LevelWarn,
		slog.LevelError.String(): slog.LevelError,
	}
)

// New generate a new zap logger instance.
func New(cfg *configs.Logger) *slog.Logger {
	level, exists := levels[strings.ToUpper(cfg.Level)]
	if !exists {
		log.Fatalf("unable to find log level: %v", cfg.Level)
	}

	handler := slog.NewTextHandler(os.Stdout)
	logger := slog.New(NewLevelHandler(level, handler))

	return logger
}
