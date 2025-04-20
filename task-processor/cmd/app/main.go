package main

import (
	"context"
	"flag"
	"log/slog"
	"os"

	"github.com/joho/godotenv"
	"github.com/tehrelt/workmate-testovoe/task-processor/internal/app"
	"github.com/tehrelt/workmate-testovoe/task-processor/pkg/sl"
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

	instance, cleanup, err := app.New(ctx)
	if err != nil {
		slog.Error("failed to create app instance", sl.Err(err))
		os.Exit(-1)
	}
	defer cleanup()

	if err := instance.Run(ctx); err != nil {
		panic(err)
	}
}
