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
	App        *App
	Auth       *Auth
	Database   *Database
	Logger     *Logger
	GRPCServer *GRPCServer
	HTTPServer *HTTPServer
	OAuth      *OAuth
	User       *User
}

func Load(ctx context.Context) *Base {
	var cfg = &Base{
		App:        newApp(),
		Auth:       newAuth(),
		Database:   newDatabase(),
		Logger:     newLogger(),
		GRPCServer: newGRPCServer(),
		HTTPServer: newHTTPServer(),
		OAuth:      newOAuth(),
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
