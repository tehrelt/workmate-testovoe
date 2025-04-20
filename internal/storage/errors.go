package storage

import "errors"

var (
	ErrTaskNotFound       = errors.New("task not found")
	ErrEventAlreadyExists = errors.New("event already exists")
)
