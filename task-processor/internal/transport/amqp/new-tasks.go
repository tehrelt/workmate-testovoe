package amqp

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/google/uuid"
	"github.com/rabbitmq/amqp091-go"
	"github.com/tehrelt/workmate-testovoe/task-processor/internal/models"
)

func (c *Consumer) handleNewTasks(ctx context.Context, msg amqp091.Delivery) (err error) {

	defer func() {
		err = msg.Ack(false)
		if err != nil {
			fmt.Printf("failed to ack message: %v", err)
		}
	}()

	payload := msg.Body

	event := &models.CreatedTaskEvent{}
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

	if err := c.ts.ProcessTask(ctx, &models.ProcessTask{
		EventId: eventId,
		TaskId:  taskId,
	}); err != nil {
		return err
	}

	return nil
}
