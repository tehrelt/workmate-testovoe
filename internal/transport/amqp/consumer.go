package amqp

import (
	"context"
	"log/slog"

	"github.com/rabbitmq/amqp091-go"
	"github.com/tehrelt/workmate-testovoe/internal/config"
	"github.com/tehrelt/workmate-testovoe/internal/lib/rmq"
	"github.com/tehrelt/workmate-testovoe/internal/services/taskservice"
	"github.com/tehrelt/workmate-testovoe/pkg/sl"
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
	go func() {
		for {
			select {
			case <-ctx.Done():
				c.manager.Close()
				return
			default:
				if err := c.manager.Consume(ctx, c.cfg.RoutingKey, c.handleProcessedMessage); err != nil {
					slog.Error("failed to consume messages", sl.Err(err))
				}
			}
		}
	}()

	<-ctx.Done()
	return ctx.Err()
}
