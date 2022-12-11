package configs

import (
	"context"
	"log"

	"github.com/heetch/confita"
	"github.com/heetch/confita/backend/env"
)

type Base struct {
	Database   *Database
	Logger     *Logger
	HTTPServer *HTTPServer
}

func Load(ctx context.Context) *Base {
	var cfg = &Base{
		Database:   newDatabase(),
		Logger:     newLogger(),
		HTTPServer: newHTTPServer(),
	}

	loader := confita.NewLoader(
		env.NewBackend(),
	)

	if err := loader.Load(ctx, cfg); err != nil {
		log.Fatalf("cannot load configuration: %v", err)
	}

	return cfg
}
