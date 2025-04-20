package taskstorage

import (
	"context"
	"errors"
	"log/slog"

	sq "github.com/Masterminds/squirrel"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v4"
	"github.com/tehrelt/workmate-testovoe/internal/lib/tracer"
	"github.com/tehrelt/workmate-testovoe/internal/models"
	"github.com/tehrelt/workmate-testovoe/internal/storage"
	"github.com/tehrelt/workmate-testovoe/internal/storage/pg"
	"github.com/tehrelt/workmate-testovoe/pkg/sl"
	"go.opentelemetry.io/otel"
)

func (ts *TaskStorage) Task(ctx context.Context, id uuid.UUID) (*models.Task, error) {
	fn := "taskstorage.Task"
	log := slog.With(slog.String("method", fn))
	ctx, span := otel.
		Tracer(tracer.TracerKey).
		Start(ctx, fn)
	defer span.End()

	query, args, err := sq.
		Select("id", "title", "status", "created_at", "updated_at").
		From(pg.TaskTable).
		Where(sq.Eq{"id": id.String()}).
		PlaceholderFormat(sq.Dollar).
		ToSql()
	if err != nil {
		log.Error("failed to build query", sl.Err(err))
		return nil, err
	}

	row := ts.pool.QueryRow(ctx, query, args...)
	task := &task{}
	if err := row.Scan(
		&task.id,
		&task.title,
		&task.status,
		&task.createdAt,
		&task.updatedAt,
	); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			log.Warn("task not found", slog.String("id", id.String()))
			return nil, storage.ErrTaskNotFound
		}
		log.Error("failed to scan row", sl.Err(err))
		return nil, err
	}

	ret, err := task.ToModel()
	if err != nil {
		log.Error("failed to convert task to model", sl.Err(err))
		return nil, err
	}

	return ret, nil
}
