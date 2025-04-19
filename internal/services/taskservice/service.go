package taskservice

import (
	"context"
	"log/slog"

	"github.com/google/uuid"
	"github.com/tehrelt/workmate-testovoe/internal/config"
	"github.com/tehrelt/workmate-testovoe/internal/models"
	"github.com/tehrelt/workmate-testovoe/pkg/sl"
)

type TaskSaver interface {
	Save(ctx context.Context, task *models.CreateTask) (*models.Task, error)
}

type TaskProvider interface {
	Task(ctx context.Context, id uuid.UUID) (*models.Task, error)
	Tasks(ctx context.Context, filter *models.TaskFilter) (<-chan models.Task, error)
}

type TaskProcessor interface {
	Push(ctx context.Context, task *models.Task) error
}

type TaskService struct {
	cfg *config.Config

	taskSaver     TaskSaver
	taskProvider  TaskProvider
	taskProcessor TaskProcessor
}

func NewTaskService(cfg *config.Config, taskSaver TaskSaver, taskProvider TaskProvider) *TaskService {
	return &TaskService{
		cfg:          cfg,
		taskSaver:    taskSaver,
		taskProvider: taskProvider,
	}
}

func (ts *TaskService) CreateTask(ctx context.Context, in *models.CreateTask) error {

	task, err := ts.taskSaver.Save(ctx, in)
	if err != nil {
		slog.Error("failed to save task", sl.Err(err))
		return err
	}

	// TODO transaction manager
	if err := ts.taskProcessor.Push(ctx, task); err != nil {
		slog.Error("failed to push task", sl.Err(err))
		return err
	}

	return nil
}

func (ts *TaskService) Tasks(ctx context.Context, filter *models.TaskFilter) (<-chan models.Task, error) {
	out, err := ts.taskProvider.Tasks(ctx, filter)
	if err != nil {
		slog.Error("failed to get tasks", sl.Err(err))
		return nil, err
	}

	return out, nil
}

func (ts *TaskService) Task(ctx context.Context, id uuid.UUID) (*models.Task, error) {
	task, err := ts.taskProvider.Task(ctx, id)
	if err != nil {
		slog.Error("failed to get task", sl.Err(err))
		return nil, err
	}

	return task, nil
}
