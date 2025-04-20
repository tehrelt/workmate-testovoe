package models

import (
	"time"

	"github.com/google/uuid"
)

type TaskStatus string

func (ts TaskStatus) String() string {
	return string(ts)
}

const (
	StatusError     TaskStatus = "error"
	StatusPending   TaskStatus = "pending"
	StatusCompleted TaskStatus = "completed"
)

type Task struct {
	Id        uuid.UUID
	Title     string
	Status    TaskStatus
	CreatedAt time.Time
	UpdatedAt *time.Time
}

type CreateTask struct {
	Title string
}

type TaskFilter struct {
	Status *TaskStatus
}
