//go:build wireinject
// +build wireinject

package app

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"github.com/google/wire"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/tehrelt/workmate-testovoe/internal/config"
	"github.com/tehrelt/workmate-testovoe/internal/lib/tracer"
	"github.com/tehrelt/workmate-testovoe/internal/services/taskservice"
	"github.com/tehrelt/workmate-testovoe/internal/storage/pg/taskstorage"
	"github.com/tehrelt/workmate-testovoe/internal/storage/rmq/taskqueue"
	"github.com/tehrelt/workmate-testovoe/internal/transport/http"
	"github.com/tehrelt/workmate-testovoe/pkg/sl"
	"go.opentelemetry.io/otel/trace"
)

//go:generate wire .
func New(ctx context.Context) (*App, func(), error) {
	panic(
		wire.Build(
			build,

			http.New,

			taskservice.New,
			wire.Bind(new(taskservice.TaskSaver), new(*taskstorage.TaskStorage)),
			wire.Bind(new(taskservice.TaskProvider), new(*taskstorage.TaskStorage)),
			wire.Bind(new(taskservice.TaskProcessor), new(*taskqueue.TaskQueue)),

			taskstorage.New,
			taskqueue.New,

			_pg,
			_tracer,
			config.New,
		),
	)
}

func _tracer(ctx context.Context, cfg *config.Config) (trace.Tracer, error) {
	jaeger := cfg.JaegerEndpoint
	appname := cfg.Name

	slog.Debug("connecting to jaeger", slog.String("jaeger", jaeger), slog.String("appname", appname))

	return tracer.SetupTracer(ctx, jaeger, appname)
}

func _pg(ctx context.Context, cfg *config.Config) (*pgxpool.Pool, func(), error) {
	pool, err := pgxpool.Connect(ctx, cfg.PG.ConnectionString())
	if err != nil {
		return nil, nil, err
	}

	log := slog.With(slog.String("connection", cfg.PG.ConnectionString()))
	log.Debug("connecting to database")
	t := time.Now()
	if err := pool.Ping(ctx); err != nil {
		log.Error("failed to connect to database", sl.Err(err))
		return nil, func() { pool.Close() }, err
	}
	log.Info("connected to database", slog.String("ping", fmt.Sprintf("%2.fs", time.Since(t).Seconds())))

	return pool, func() { pool.Close() }, nil
}
