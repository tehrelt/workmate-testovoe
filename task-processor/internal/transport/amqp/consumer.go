package amqp

import (
	"context"
	"log/slog"

	"github.com/rabbitmq/amqp091-go"
	"github.com/tehrelt/workmate-testovoe/task-processor/internal/config"
	"github.com/tehrelt/workmate-testovoe/task-processor/internal/lib/rmq"
	"github.com/tehrelt/workmate-testovoe/task-processor/internal/services/taskservice"
	"github.com/tehrelt/workmate-testovoe/task-processor/pkg/sl"
)

type Consumer struct {
	cfg     config.QueueConfig
	manager *rmq.Manager
	ts      *taskservice.Service
}

func New(ch *amqp091.Channel, cfg config.QueueConfig, ts *taskservice.Service) *Consumer {
	return &Consumer{
		cfg:     cfg,
		manager: rmq.New(ch),
		ts:      ts,
	}
}

func (c *Consumer) Run(ctx context.Context) error {
	slog.Info("starting up workers for consumer", slog.Any("config", c.cfg))
	for range c.cfg.WorkerPoolConfig.MaxWorkers {
		go func() {
			if err := c.manager.Consume(ctx, c.cfg.RoutingKey, c.handleNewTasks); err != nil {
				slog.Error("failed to consume messages", sl.Err(err))
				slog.Debug("restaring worker")
			}
		}()
	}

	<-ctx.Done()
	return ctx.Err()
}
