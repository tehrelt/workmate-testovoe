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
	StatusCompleted TaskStatus = "done"
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

type UpdateTask struct {
	Id        uuid.UUID
	EventId   uuid.UUID
	NewStatus TaskStatus
}

type Range struct {
	From  uint64
	Limit uint64
}

type TaskFilter struct {
	Range
	Status TaskStatus
}

func (tf *TaskFilter) SetStatus(status TaskStatus) *TaskFilter {
	tf.Status = status
	return tf
}

func (tf *TaskFilter) SetFrom(from uint64) *TaskFilter {
	tf.Range.From = from
	return tf
}

func (tf *TaskFilter) SetLimit(limit uint64) *TaskFilter {
	tf.Range.Limit = limit
	return tf
}
