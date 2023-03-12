package log

import (
	"context"

	"golang.org/x/exp/slog"
)

type nopHandler struct{}

// NewNopHandler returns an empty handler.
func NewNopHandler() *nopHandler {
	return &nopHandler{}
}

func (h *nopHandler) Enabled(_ context.Context, level slog.Level) bool {
	return true
}

func (h *nopHandler) Handle(_ context.Context, r slog.Record) error {
	return nil
}

func (h *nopHandler) WithAttrs(attrs []slog.Attr) slog.Handler {
	return nil
}

func (h *nopHandler) WithGroup(name string) slog.Handler {
	return nil
}
