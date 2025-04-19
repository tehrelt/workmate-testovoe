package main

import (
	"context"
	"flag"
	"log/slog"
	"os"

	"github.com/joho/godotenv"
	"github.com/tehrelt/workmate-testovoe/internal/config"
	"github.com/tehrelt/workmate-testovoe/internal/lib/tracer"
	"github.com/tehrelt/workmate-testovoe/internal/services/taskservice"
	"github.com/tehrelt/workmate-testovoe/internal/storage/pg"
	"github.com/tehrelt/workmate-testovoe/internal/storage/pg/taskstorage"
	"github.com/tehrelt/workmate-testovoe/internal/storage/rmq/taskqueue"
	"github.com/tehrelt/workmate-testovoe/internal/transport/http"
	"github.com/tehrelt/workmate-testovoe/pkg/sl"
)

var (
	envPath string
)

func init() {
	flag.StringVar(&envPath, "env", "", "Path to the environment file")
}

func main() {
	flag.Parse()

	if envPath != "" {
		if err := godotenv.Load(envPath); err != nil {
			slog.Error("failed to parse env file", sl.Err(err))
			os.Exit(-1)
		}
	}

	ctx := context.Background()
	cfg := config.New()
	pool, closePool, err := pg.New(ctx, cfg)
	if err != nil {
		slog.Error("failed to create postgres pool", sl.Err(err))
		os.Exit(-1)
	}
	defer closePool()

	taskStorage := taskstorage.New(pool)
	taskqueue := taskqueue.New()

	taskService := taskservice.New(taskStorage, taskStorage, taskqueue)

	tracer.SetupTracer(ctx, cfg.JaegerEndpoint, cfg.Name)
	server := http.New(cfg, taskService)

	if err := server.Run(ctx); err != nil {
		panic(err)
	}
}
