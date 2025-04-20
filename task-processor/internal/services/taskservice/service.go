package taskservice

import (
	"context"
	"log/slog"
	"math/rand"
	"time"

	"github.com/google/uuid"
	"github.com/tehrelt/workmate-testovoe/task-processor/internal/lib/tracer"
	"github.com/tehrelt/workmate-testovoe/task-processor/internal/models"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
)

type EventSaver interface {
	Save(ctx context.Context, id uuid.UUID) error
}

type TaskNotifier interface {
	Push(ctx context.Context, in *models.ProcessedTaskEvent) error
}

type Service struct {
	eventSaver   EventSaver
	taskNotifier TaskNotifier
}

func New(eventSaver EventSaver, taskNotifier TaskNotifier) *Service {
	return &Service{
		eventSaver:   eventSaver,
		taskNotifier: taskNotifier,
	}
}

func (s *Service) ProcessTask(ctx context.Context, in *models.ProcessTask) error {
	fn := "taskservice.ProcessTask"
	log := slog.With(slog.String("fn", fn))
	ctx, span := otel.Tracer(tracer.TracerKey).Start(ctx, fn)
	defer span.End()

	if err := s.eventSaver.Save(ctx, in.EventId); err != nil {
		return err
	}

	seconds := 10 + rand.Intn(20)
	span.SetAttributes(attribute.Int("desired_process_time", seconds))
	log.Info("start process task", slog.String("taskId", in.TaskId.String()), slog.Int("seconds", seconds))
	time.Sleep(time.Duration(seconds) * time.Second)
	log.Info("ended process task", slog.String("taskId", in.TaskId.String()))

	event := &models.ProcessedTaskEvent{
		EventId: in.EventId.String(),
		TaskId:  in.TaskId.String(),
	}

	if err := s.taskNotifier.Push(ctx, event); err != nil {
		return err
	}

	return nil
}
