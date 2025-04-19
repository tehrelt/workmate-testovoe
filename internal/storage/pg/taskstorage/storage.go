package taskstorage

import (
	"context"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/tehrelt/workmate-testovoe/internal/models"
	"github.com/tehrelt/workmate-testovoe/internal/services/taskservice"
)

var _ taskservice.TaskSaver = (*TaskStorage)(nil)
var _ taskservice.TaskProvider = (*TaskStorage)(nil)

type TaskStorage struct {
	pool *pgxpool.Pool
}

// Task implements taskservice.TaskProvider.
func (ts *TaskStorage) Task(ctx context.Context, id uuid.UUID) (*models.Task, error) {
	panic("unimplemented")
}

// Tasks implements taskservice.TaskProvider.
func (ts *TaskStorage) Tasks(ctx context.Context, filter *models.TaskFilter) (<-chan models.Task, error) {
	panic("unimplemented")
}

func New(pool *pgxpool.Pool) *TaskStorage {
	return &TaskStorage{
		pool: pool,
	}
}
