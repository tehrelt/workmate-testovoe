package taskstorage

import (
	"context"
	"log/slog"

	sq "github.com/Masterminds/squirrel"
	"github.com/tehrelt/workmate-testovoe/task-producer/internal/lib/tracer"
	"github.com/tehrelt/workmate-testovoe/task-producer/internal/lib/tx"
	"github.com/tehrelt/workmate-testovoe/task-producer/internal/models"
	"github.com/tehrelt/workmate-testovoe/task-producer/internal/storage/pg"
	"github.com/tehrelt/workmate-testovoe/task-producer/pkg/sl"
	"go.opentelemetry.io/otel"
)

func (ts *TaskStorage) Save(ctx context.Context, in *models.CreateTask) (*models.Task, error) {
	fn := "taskstorage.Save"
	log := slog.With(slog.String("method", fn))
	ctx, span := otel.
		Tracer(tracer.TracerKey).
		Start(ctx, fn)
	defer span.End()

	ctx, tx, err := tx.GetOrDefault(ctx, ts.pool)
	if err != nil {
		log.Error("failed to get transaction", sl.Err(err))
		return nil, err
	}

	query, args, err := sq.
		Insert(pg.TaskTable).
		Columns("title").
		Values(in.Title).
		Suffix("RETURNING *").
		PlaceholderFormat(sq.Dollar).
		ToSql()
	if err != nil {
		log.Error("failed to build query", sl.Err(err))
		return nil, err
	}

	row := tx.QueryRow(ctx, query, args...)
	task := &task{}
	if err := row.Scan(
		&task.id,
		&task.title,
		&task.status,
		&task.createdAt,
		&task.updatedAt,
	); err != nil {
		log.Error("failed to scan row", sl.Err(err))
		return nil, err
	}

	ret, err := task.ToModel()
	if err != nil {
		log.Error("failed to convert task to model", sl.Err(err))
		return nil, err
	}

	log.Debug("task saved", slog.Any("task", ret))

	return ret, nil
}
