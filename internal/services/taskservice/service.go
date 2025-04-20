package taskservice

import (
	"context"
	"log/slog"

	"github.com/google/uuid"
	"github.com/tehrelt/workmate-testovoe/internal/models"
	"github.com/tehrelt/workmate-testovoe/pkg/sl"
)

//go:generate go run github.com/vektra/mockery/v2 --name=TaskSaver
type TaskSaver interface {
	Save(ctx context.Context, task *models.CreateTask) (*models.Task, error)
	Update(ctx context.Context, task *models.UpdateTask) error
}

//go:generate go run github.com/vektra/mockery/v2 --name=TaskProvider
type TaskProvider interface {
	Task(ctx context.Context, id uuid.UUID) (*models.Task, error)
	Tasks(ctx context.Context, filter *models.TaskFilter) (<-chan *models.Task, error)
	Total(ctx context.Context) (uint64, error)
}

type EventSaver interface {
	Save(ctx context.Context, eventId uuid.UUID) error
}

//go:generate go run github.com/vektra/mockery/v2 --name=TaskProcessor
type TaskProcessor interface {
	Push(ctx context.Context, task *models.CreatedTaskEvent) error
}

type TaskService struct {
	taskSaver     TaskSaver
	taskProvider  TaskProvider
	taskProcessor TaskProcessor

	eventSaver EventSaver
}

func New(
	taskSaver TaskSaver,
	taskProvider TaskProvider,
	taskProcessor TaskProcessor,
	eventSaver EventSaver,
) *TaskService {
	return &TaskService{
		taskSaver:     taskSaver,
		taskProvider:  taskProvider,
		taskProcessor: taskProcessor,
		eventSaver:    eventSaver,
	}
}

func (ts *TaskService) CreateTask(ctx context.Context, in *models.CreateTask) (*models.Task, error) {
	slog.Debug("saving task", slog.Any("in", in))
	task, err := ts.taskSaver.Save(ctx, in)
	if err != nil {
		slog.Error("failed to save task", sl.Err(err))
		return nil, err
	}

	event := &models.CreatedTaskEvent{
		TaskId:  task.Id.String(),
		EventId: uuid.NewString(),
	}

	// TODO transaction manager
	if err := ts.taskProcessor.Push(ctx, event); err != nil {
		slog.Error("failed to push task", sl.Err(err))
		return nil, err
	}

	return task, nil
}

func (ts *TaskService) Tasks(ctx context.Context, filter *models.TaskFilter) (<-chan *models.Task, uint64, error) {
	out, err := ts.taskProvider.Tasks(ctx, filter)
	if err != nil {
		slog.Error("failed to get tasks", sl.Err(err))
		return nil, 0, err
	}

	total, err := ts.taskProvider.Total(ctx)
	if err != nil {
		slog.Error("failed to count tasks", sl.Err(err))
		return nil, 0, err
	}

	return out, total, nil
}

func (ts *TaskService) Task(ctx context.Context, id uuid.UUID) (*models.Task, error) {
	task, err := ts.taskProvider.Task(ctx, id)
	if err != nil {
		slog.Error("failed to get task", sl.Err(err))
		return nil, err
	}

	return task, nil
}

func (ts *TaskService) UpdateTask(ctx context.Context, in *models.UpdateTask) error {

	// TODO transaction manager
	slog.Info("saving event", slog.String("eventId", in.EventId.String()))
	err := ts.eventSaver.Save(ctx, in.EventId)
	if err != nil {
		return err
	}

	slog.Info(
		"updating task status",
		slog.String("taskId", in.Id.String()),
		slog.String("newStatus", string(in.NewStatus)),
	)
	if err := ts.taskSaver.Update(ctx, in); err != nil {
		slog.Error("failed to update task", sl.Err(err))
		return err
	}

	return nil
}
