package taskqueue

import (
	"context"
	"encoding/json"
	"log/slog"

	"github.com/rabbitmq/amqp091-go"
	"github.com/tehrelt/workmate-testovoe/task-producer/internal/config"
	"github.com/tehrelt/workmate-testovoe/task-producer/internal/lib/rmq"
	"github.com/tehrelt/workmate-testovoe/task-producer/internal/models"
	"github.com/tehrelt/workmate-testovoe/task-producer/internal/services/taskservice"
)

var _ taskservice.TaskProcessor = (*TaskQueue)(nil)

type TaskQueue struct {
	manager *rmq.Manager
	cfg     config.QueueConfig
}

func New(cfg config.QueueConfig, ch *amqp091.Channel) *TaskQueue {
	return &TaskQueue{
		cfg:     cfg,
		manager: rmq.New(ch),
	}
}

func (t *TaskQueue) Push(ctx context.Context, event *models.CreatedTaskEvent) error {
	slog.Info("Pushing task", slog.Any("event", event))

	j, err := json.Marshal(event)
	if err != nil {
		return err
	}

	return t.manager.Publish(ctx, t.cfg.Exchange, t.cfg.RoutingKey, j)
}
