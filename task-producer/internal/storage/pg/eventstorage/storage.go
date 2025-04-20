package eventstorage

import (
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/tehrelt/workmate-testovoe/task-producer/internal/services/taskservice"
)

var _ taskservice.EventSaver = (*Storage)(nil)

type Storage struct {
	pool *pgxpool.Pool
}

func New(pool *pgxpool.Pool) *Storage {
	return &Storage{
		pool: pool,
	}
}
