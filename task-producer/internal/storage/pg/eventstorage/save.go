package eventstorage

import (
	"context"
	"errors"
	"fmt"
	"log/slog"

	sq "github.com/Masterminds/squirrel"
	"github.com/google/uuid"
	"github.com/jackc/pgconn"
	"github.com/tehrelt/workmate-testovoe/task-producer/internal/lib/tracer"
	"github.com/tehrelt/workmate-testovoe/task-producer/internal/lib/tx"
	"github.com/tehrelt/workmate-testovoe/task-producer/internal/storage"
	"github.com/tehrelt/workmate-testovoe/task-producer/internal/storage/pg"
	"github.com/tehrelt/workmate-testovoe/task-producer/pkg/sl"
	"go.opentelemetry.io/otel"
)

func (ts *Storage) Save(ctx context.Context, id uuid.UUID) error {
	fn := "eventstorage.Save"
	log := slog.With(slog.String("method", fn))
	ctx, span := otel.
		Tracer(tracer.TracerKey).
		Start(ctx, fn)
	defer span.End()

	ctx, tx, err := tx.GetOrDefault(ctx, ts.pool)

	query, args, err := sq.
		Insert(pg.EventTable).
		Columns("id").
		Values(id.String()).
		PlaceholderFormat(sq.Dollar).
		ToSql()
	if err != nil {
		log.Error("failed to build query", sl.Err(err))
		return err
	}

	if _, err := tx.Exec(ctx, query, args...); err != nil {
		log.Error("error ocurred", slog.String("type", fmt.Sprintf("%t", err)))
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			log.Error("occured pg error", slog.Any("pgerr", pgErr))
			if pgErr.Code == "23505" {
				slog.Error("event already exists", slog.String("id", id.String()))
				return storage.ErrEventAlreadyExists
			}
		}
		log.Error("failed to execute query", sl.Err(err))
		return err
	}

	return nil
}
