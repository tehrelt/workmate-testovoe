package app

import (
	"context"

	"github.com/tehrelt/workmate-testovoe/internal/config"
	"github.com/tehrelt/workmate-testovoe/internal/transport/http"
	"go.opentelemetry.io/otel/trace"
)

type App struct {
	cfg    *config.Config
	server *http.Server
	tracer trace.Tracer
}

func build(cfg *config.Config, server *http.Server, tracer trace.Tracer) *App {
	return &App{
		cfg:    cfg,
		server: server,
		tracer: tracer,
	}
}

func (a *App) Run(ctx context.Context) error {
	return a.server.Run(ctx)
}
