package taskstorage

import (
	"context"
	"log/slog"

	sq "github.com/Masterminds/squirrel"
	"github.com/tehrelt/workmate-testovoe/task-producer/internal/lib/tracer"
	"github.com/tehrelt/workmate-testovoe/task-producer/internal/storage/pg"
	"github.com/tehrelt/workmate-testovoe/task-producer/pkg/sl"
	"go.opentelemetry.io/otel"
)

func (ts *TaskStorage) Total(ctx context.Context) (uint64, error) {
	fn := "taskstorage.Total"
	log := slog.With(slog.String("method", fn))
	ctx, span := otel.
		Tracer(tracer.TracerKey).
		Start(ctx, fn)
	defer span.End()

	query, args, err := sq.
		Select("count(id)").
		From(pg.TaskTable).
		PlaceholderFormat(sq.Dollar).
		ToSql()
	if err != nil {
		log.Error("failed to build query", sl.Err(err))
		return 0, err
	}

	row := ts.pool.QueryRow(ctx, query, args...)
	var total uint64
	if err := row.Scan(
		&total,
	); err != nil {
		log.Error("failed to scan row", sl.Err(err))
		return 0, err
	}

	return total, nil
}
