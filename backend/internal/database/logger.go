package database

import (
	"context"
	"errors"
	"fmt"
	"time"

	"golang.org/x/exp/slog"
	gorml "gorm.io/gorm/logger"
)

// Logger is used to replace the default Gorm logger.
type Logger struct {
	logger                    *slog.Logger
	slowThreshold             time.Duration
	ignoreRecordNotFoundError bool
}

func logger(logger *slog.Logger) Logger {
	return Logger{
		logger: logger,
	}
}

// LogMode allows to set the current log level.
func (l Logger) LogMode(lvl gorml.LogLevel) gorml.Interface {
	return l
}

// Info allows to print a new informational message.
func (l Logger) Info(ctx context.Context, msg string, args ...any) {
	l.logger.Info(msg)
	l.call().Info(msg, args...)
}

// Warn allows to print a new warning message.
func (l Logger) Warn(ctx context.Context, msg string, args ...any) {
	l.logger.Info(msg)
	l.call().Warn(msg, args...)
}

// Error allows to print a new error message.
func (l Logger) Error(ctx context.Context, msg string, args ...any) {
	l.logger.Info(msg)
	l.call().Error(msg, args...)
}

// Trace allows to print a new trace message.
func (l Logger) Trace(ctx context.Context, begin time.Time, fc func() (string, int64), err error) {
	elapsed := time.Since(begin)
	sql, rows := fc()

	switch {
	case err != nil && (!errors.Is(err, gorml.ErrRecordNotFound) || !l.ignoreRecordNotFoundError):
		l.call().Error("trace", err, slog.Duration("elapsed", elapsed), slog.Int64("rows", rows), slog.String("sql", sql))
	case elapsed > l.slowThreshold && l.slowThreshold != 0:
		slowLog := fmt.Sprintf("SLOW SQL >= %v", l.slowThreshold)
		l.call().Warn("trace", slog.String("slow", slowLog), slog.Duration("elapsed", elapsed), slog.Int64("rows", rows), slog.String("sql", sql))
	default:
		l.call().Debug("trace", slog.Duration("elapsed", elapsed), slog.Int64("rows", rows), slog.String("sql", sql))
	}
}

func (l Logger) call() *slog.Logger {
	return l.logger
}
