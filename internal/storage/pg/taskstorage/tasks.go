package taskstorage

import (
	"context"
	"log/slog"

	sq "github.com/Masterminds/squirrel"
	"github.com/tehrelt/workmate-testovoe/internal/lib/tracer"
	"github.com/tehrelt/workmate-testovoe/internal/models"
	"github.com/tehrelt/workmate-testovoe/internal/storage/pg"
	"github.com/tehrelt/workmate-testovoe/pkg/sl"
	"go.opentelemetry.io/otel"
)

// Tasks implements taskservice.TaskProvider.
func (ts *TaskStorage) Tasks(ctx context.Context, filter *models.TaskFilter) (<-chan *models.Task, error) {

	out := make(chan *models.Task, 10)

	fn := "taskstorage.Tasks"
	log := slog.With(slog.String("method", fn))
	ctx, span := otel.
		Tracer(tracer.TracerKey).
		Start(ctx, fn)
	defer span.End()

	builder := sq.
		Select("id", "title", "status", "created_at", "updated_at").
		From(pg.TaskTable).
		OrderBy("created_at ASC").
		PlaceholderFormat(sq.Dollar)

	if filter != nil {
		if filter.Status != "" {
			builder = builder.Where(sq.Eq{"status": filter.Status})
		}

		if filter.From != 0 {
			builder = builder.Offset(filter.From)
		}

		if filter.Limit != 0 {
			builder = builder.Limit(filter.Limit)
		}
	}

	query, args, err := builder.ToSql()

	if err != nil {
		log.Error("failed to build query", sl.Err(err))
		return nil, err
	}

	rows, err := ts.pool.Query(ctx, query, args...)
	if err != nil {
		log.Error("failed to query", sl.Err(err))
		return nil, err
	}

	go func() {
		defer close(out)
		defer rows.Close()

		for rows.Next() {
			task := &task{}
			if err := rows.Scan(
				&task.id,
				&task.title,
				&task.status,
				&task.createdAt,
				&task.updatedAt,
			); err != nil {
				log.Error("failed to scan row", sl.Err(err))
				continue
			}

			m, err := task.ToModel()
			if err != nil {
				log.Error("failed to convert task to model", sl.Err(err))
				continue
			}

			out <- m
		}
	}()

	return out, nil
}
