package taskstorage

import "github.com/jackc/pgx/v4/pgxpool"

type TaskStorage struct {
	pool *pgxpool.Pool
}

func NewTaskStorage(pool *pgxpool.Pool) *TaskStorage {
	return &TaskStorage{
		pool: pool,
	}
}
