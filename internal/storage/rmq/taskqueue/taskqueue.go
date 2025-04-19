package taskqueue

import (
	"context"
	"log/slog"

	"github.com/tehrelt/workmate-testovoe/internal/models"
	"github.com/tehrelt/workmate-testovoe/internal/services/taskservice"
)

var _ taskservice.TaskProcessor = (*TaskQueue)(nil)

type TaskQueue struct {
	// Add fields here
}

func New() *TaskQueue {
	return &TaskQueue{}
}

// Push implements taskservice.TaskProcessor.
func (t *TaskQueue) Push(ctx context.Context, task *models.Task) error {
	slog.Info("Pushing task", slog.Any("task", task))
	return nil
}
