package configs

import (
	"context"
	"log"

	"github.com/heetch/confita"
	"github.com/heetch/confita/backend/env"
)

const (
	ApplicationName = "authz"
)

type Base struct {
	Auth       *Auth
	Database   *Database
	Logger     *Logger
	GRPCServer *GRPCServer
	HTTPServer *HTTPServer
	User       *User
}

func Load(ctx context.Context) *Base {
	var cfg = &Base{
		Auth:       newAuth(),
		Database:   newDatabase(),
		Logger:     newLogger(),
		GRPCServer: newGRPCServer(),
		HTTPServer: newHTTPServer(),
		User:       newUser(),
	}

	loader := confita.NewLoader(
		env.NewBackend(),
	)

	if err := loader.Load(ctx, cfg); err != nil {
		log.Fatalf("cannot load configuration: %v", err)
	}

	return cfg
}
