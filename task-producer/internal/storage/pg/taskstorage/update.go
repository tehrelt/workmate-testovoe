package taskstorage

import (
	"context"
	"log/slog"
	"time"

	sq "github.com/Masterminds/squirrel"
	"github.com/tehrelt/workmate-testovoe/task-producer/internal/lib/tracer"
	"github.com/tehrelt/workmate-testovoe/task-producer/internal/models"
	"github.com/tehrelt/workmate-testovoe/task-producer/internal/storage/pg"
	"github.com/tehrelt/workmate-testovoe/task-producer/pkg/sl"
	"go.opentelemetry.io/otel"
)

func (ts *TaskStorage) Update(ctx context.Context, in *models.UpdateTask) error {
	fn := "taskstorage.Update"
	log := slog.With(slog.String("method", fn))
	ctx, span := otel.
		Tracer(tracer.TracerKey).
		Start(ctx, fn)
	defer span.End()

	query, args, err := sq.
		Update(pg.TaskTable).
		Set("status", in.NewStatus).
		Set("updated_at", time.Now()).
		Where(sq.Eq{"id": in.Id}).
		PlaceholderFormat(sq.Dollar).
		ToSql()
	if err != nil {
		log.Error("failed to build query", sl.Err(err))
		return err
	}

	if _, err := ts.pool.Exec(ctx, query, args...); err != nil {
		log.Error("failed to update task", sl.Err(err))
		return err
	}

	return nil
}
