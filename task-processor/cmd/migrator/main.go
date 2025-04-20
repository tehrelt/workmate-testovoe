package main

import (
	"errors"
	"flag"
	"fmt"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/joho/godotenv"
	"github.com/tehrelt/workmate-testovoe/task-processor/internal/config"
)

var (
	local bool
)

func init() {
	flag.BoolVar(&local, "local", false, "run in local mode")
}

func main() {
	flag.Parse()

	if local {
		if err := godotenv.Load(); err != nil {
			panic(fmt.Errorf("cannot load env: %w", err))
		}
	}

	cfg := config.New()

	pg := cfg.PG
	m, err := migrate.New(
		"file://migrations",
		pg.ConnectionString(),
	)
	if err != nil {
		panic(err)
	}

	if err := m.Up(); err != nil {
		if !errors.Is(err, migrate.ErrNoChange) {
			panic(err)
		}
	}
}
