package oauth

import (
	"log"

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
	srv.SetResponseErrorHandler(func(re *errors.Response) {
		log.Printf("Response Error: %#+v\n", re)
	})

	return srv
}
