package amqp

import (
	"context"
	"log/slog"

	"github.com/rabbitmq/amqp091-go"
	"github.com/tehrelt/workmate-testovoe/task-producer/internal/config"
	"github.com/tehrelt/workmate-testovoe/task-producer/internal/lib/rmq"
	"github.com/tehrelt/workmate-testovoe/task-producer/internal/services/taskservice"
	"github.com/tehrelt/workmate-testovoe/task-producer/pkg/sl"
)

type Consumer struct {
	cfg         config.QueueConfig
	taskService *taskservice.TaskService
	manager     *rmq.Manager
}

func New(ch *amqp091.Channel, cfg config.QueueConfig, ts *taskservice.TaskService) *Consumer {

	manager := rmq.New(ch)
	slog.Debug("creating consumer", slog.Any("cfg", cfg), slog.Any("manager", manager))
	return &Consumer{
		cfg:         cfg,
		manager:     manager,
		taskService: ts,
	}
}

func (c *Consumer) Run(ctx context.Context) error {

	slog.Info("starting up consumers", slog.Any("config", c.cfg))
	for range c.cfg.PoolConfig.MaxWorkers {
		go func() {
			if err := c.manager.Consume(ctx, c.cfg.RoutingKey, c.handleProcessedMessage); err != nil {
				slog.Error("failed to consume messages", sl.Err(err))
			}
		}()
	}

	<-ctx.Done()
	return ctx.Err()
}
