package services

import "errors"

var (
	ErrEventAlreadyProcessed = errors.New("event already processed")
	ErrEventOnTimeout        = errors.New("event on timeout")
)
