package server

import (
	"context"

	"github.com/go-oauth2/oauth2/v4"
	"github.com/go-oauth2/oauth2/v4/errors"
	"github.com/go-oauth2/oauth2/v4/server"
	"golang.org/x/exp/slog"
)

func NewServer(
	logger *slog.Logger,
	manager oauth2.Manager,
) *server.Server {
	srv := server.NewDefaultServer(manager)
	srv.SetAllowGetAccessRequest(true)
	srv.SetClientInfoHandler(server.ClientFormHandler)
	srv.SetResponseErrorHandler(func(_ context.Context, re *errors.Response) {
		logger.Error(
			"OAuth response error",
			slog.Any("error", re.Error),
			slog.Int("status_code", re.StatusCode),
			slog.String("description", re.Description),
		)
	})

	return srv
}
