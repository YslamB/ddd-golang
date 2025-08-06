package common

import "errors"

type CustomError interface {
	Error() string
	Is(error) bool
}

type NotFoundError struct {
	Message string
}

func (e *NotFoundError) Error() string {
	return e.Message
}

func (e *NotFoundError) Is(target error) bool {
	_, ok := target.(*NotFoundError)
	return ok
}

func NewNotFoundError(message string) *NotFoundError {
	return &NotFoundError{Message: message}
}

type ConflictError struct {
	Message string
}

func (e *ConflictError) Error() string {
	return e.Message
}

func (e *ConflictError) Is(target error) bool {
	_, ok := target.(*ConflictError)
	return ok
}

func NewConflictError(message string) *ConflictError {
	return &ConflictError{Message: message}
}

type ValidationError struct {
	Message string
}

func (e *ValidationError) Error() string {
	return e.Message
}

func (e *ValidationError) Is(target error) bool {
	_, ok := target.(*ValidationError)
	return ok
}

func NewValidationError(message string) *ValidationError {
	return &ValidationError{Message: message}
}

func IsNotFoundError(err error) bool {
	var target *NotFoundError
	return errors.As(err, &target)
}

func IsConflictError(err error) bool {
	var target *ConflictError
	return errors.As(err, &target)
}

func IsValidationError(err error) bool {
	var target *ValidationError
	return errors.As(err, &target)
}
