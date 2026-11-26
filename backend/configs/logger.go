package configs

import "log/slog"

type Logger struct {
	Level string `config:"logger_level"`
}

func newLogger() *Logger {
	return &Logger{
		Level: slog.LevelInfo.String(),
	}
}
