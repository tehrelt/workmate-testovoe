package main

import (
	"context"
	"flag"
	"log/slog"
	"os"

	"github.com/joho/godotenv"
	"github.com/tehrelt/workmate-testovoe/internal/config"
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
	server := http.New(cfg)

	if err := server.Run(ctx); err != nil {
		panic(err)
	}
}
