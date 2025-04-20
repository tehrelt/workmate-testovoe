package amqp

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/rabbitmq/amqp091-go"
	"github.com/tehrelt/workmate-testovoe/internal/models"
	"github.com/tehrelt/workmate-testovoe/internal/storage"
)

func (c *Consumer) handleProcessedMessage(ctx context.Context, msg amqp091.Delivery) (err error) {

	defer func() {
		if err != nil && !errors.Is(err, storage.ErrEventAlreadyExists) {
			err = msg.Reject(false)
		}

		err = msg.Ack(false)
		if err != nil {
			fmt.Printf("failed to ack message: %v", err)
		}
	}()

	payload := msg.Body

	event := &models.ProcessedTaskEvent{}
	if err := json.Unmarshal(payload, event); err != nil {
		return err
	}

	taskId, err := uuid.Parse(event.TaskId)
	if err != nil {
		return err
	}

	eventId, err := uuid.Parse(event.EventId)
	if err != nil {
		return err
	}

	newStatus := models.StatusCompleted
	if event.Error != "" {
		newStatus = models.StatusError
	}

	if err := c.taskService.UpdateTask(ctx, &models.UpdateTask{
		Id:        taskId,
		EventId:   eventId,
		NewStatus: newStatus,
	}); err != nil {
		return err
	}

	return nil

}
