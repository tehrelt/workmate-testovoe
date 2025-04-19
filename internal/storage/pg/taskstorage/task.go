package taskstorage

import (
	"time"

	"github.com/google/uuid"
	"github.com/tehrelt/workmate-testovoe/internal/models"
)

type task struct {
	id        string
	title     string
	status    string
	createdAt time.Time
	updatedAt *time.Time
}

func (t *task) ToModel() (*models.Task, error) {

	id, err := uuid.Parse(t.id)
	if err != nil {
		return nil, err
	}

	return &models.Task{
		Id:        id,
		Title:     t.title,
		Status:    models.TaskStatus(t.status),
		CreatedAt: t.createdAt,
		UpdatedAt: t.updatedAt,
	}, nil
}
