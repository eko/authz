package configs

import "golang.org/x/exp/slog"

type Logger struct {
	Level string `config:"logger_level"`
}

func newLogger() *Logger {
	return &Logger{
		Level: slog.InfoLevel.String(),
	}
}
