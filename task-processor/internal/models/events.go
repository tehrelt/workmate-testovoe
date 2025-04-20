package models

import (
	"time"

	"github.com/google/uuid"
)

type Event struct {
	Id              uuid.UUID `json:"id"`
	Processed       bool      `json:"processed"`
	CreatedAt       time.Time `json:"createdAt"`
	LastProcessedAt time.Time `json:"lastProcessedAt"`
}

type ProcessedTaskEvent struct {
	EventId     string    `json:"eventId"`
	TaskId      string    `json:"taskId"`
	ProcessedAt time.Time `json:"processedAt"`
	Error       string    `json:"error"`
	Result      string    `json:"result"`
}

type CreatedTaskEvent struct {
	EventId   string    `json:"eventId" validate:"required"`
	TaskId    string    `json:"taskId" validate:"required"`
	CreatedAt time.Time `json:"createdAt"`
}

type ProcessTask struct {
	EventId uuid.UUID `json:"eventId"`
	TaskId  uuid.UUID `json:"taskId"`
}
