package amqp

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"

	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"github.com/rabbitmq/amqp091-go"
	"github.com/tehrelt/workmate-testovoe/task-processor/internal/models"
	"github.com/tehrelt/workmate-testovoe/task-processor/pkg/sl"
)

func (c *Consumer) handleNewTasks(ctx context.Context, msg amqp091.Delivery) (err error) {

	fn := "handleNewTasks"
	log := slog.With(slog.String("fn", fn))

	defer func() {
		if err != nil {
			log.Error("failed to handle new tasks", sl.Err(err))
			msg.Reject(false)
			return
		}

		err = msg.Ack(false)
		if err != nil {
			fmt.Printf("failed to ack message: %v", err)
		}
	}()

	payload := msg.Body

	log.Info("incoming tasks.new event")
	event := &models.CreatedTaskEvent{}
	if err := json.Unmarshal(payload, event); err != nil {
		return err
	}

	v := validator.New()
	if err := v.Struct(event); err != nil {
		log.Warn("invalid event message, rejecting message", slog.String("payload", string(payload)), sl.Err(err))
		if err := msg.Reject(false); err != nil {
			log.Error("failed to reject message", sl.Err(err))
			return err
		}

		return err
	}
	log.Info("", args ...any)

	taskId, err := uuid.Parse(event.TaskId)
	if err != nil {
		return err
	}

	eventId, err := uuid.Parse(event.EventId)
	if err != nil {
		return err
	}

	if err := c.ts.ProcessTask(ctx, &models.ProcessTask{
		EventId: eventId,
		TaskId:  taskId,
	}); err != nil {
		return err
	}

	return nil
}
