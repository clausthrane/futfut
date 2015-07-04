package services

import (
	"fmt"
)

type serviceTimeoutError struct {
	msg string
}

func (e *serviceTimeoutError) Error() string {
	return fmt.Sprintf("Timeout: %s", e.msg)
}

func NewServiceTimeoutError(msg string) error {
	return &serviceTimeoutError{msg}
}
