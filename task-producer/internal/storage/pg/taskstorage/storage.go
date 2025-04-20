package taskstorage

import (
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/tehrelt/workmate-testovoe/task-producer/internal/services/taskservice"
)

var _ taskservice.TaskSaver = (*TaskStorage)(nil)
var _ taskservice.TaskProvider = (*TaskStorage)(nil)

type TaskStorage struct {
	pool *pgxpool.Pool
}

func New(pool *pgxpool.Pool) *TaskStorage {
	return &TaskStorage{
		pool: pool,
	}
}
