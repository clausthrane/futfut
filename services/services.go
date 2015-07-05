// The service package implements all "business logic" in Futfut
//
// the train subpackage implementes train specific logic
// the station subpackage implements station specific logic

package services

import (
	"fmt"
)

// StationID represents the id of a train station
type StationID string

// TrainID represents the id of a train
type TrainID string

type ServiceTimeoutError struct {
	msg string
}

func (e *ServiceTimeoutError) Error() string {
	return fmt.Sprintf("Timeout: %s", e.msg)
}

func NewServiceTimeoutError(msg string) error {
	return &ServiceTimeoutError{msg}
}

type ServiceValidationError struct {
	msg string
}

func NewServiceValidationError(msg string) error {
	return &ServiceValidationError{msg}
}

func (e *ServiceValidationError) Error() string {
	return fmt.Sprintf("Validation error: %s", e.msg)
}
