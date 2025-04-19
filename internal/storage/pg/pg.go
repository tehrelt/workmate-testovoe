package pg

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/tehrelt/workmate-testovoe/internal/config"
	"github.com/tehrelt/workmate-testovoe/pkg/sl"
)

func New(ctx context.Context, cfg *config.Config) (*pgxpool.Pool, func(), error) {
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
