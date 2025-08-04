package common

import "errors"

// CustomError is a base interface for all custom application errors.
type CustomError interface {
	Error() string
	Is(error) bool // For errors.Is
}

// NotFoundError represents an error where a requested resource was not found.
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

// NewNotFoundError creates a new NotFoundError.
func NewNotFoundError(message string) *NotFoundError {
	return &NotFoundError{Message: message}
}

// ConflictError represents an error where a request conflicts with the current state of the resource.
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

// NewConflictError creates a new ConflictError.
func NewConflictError(message string) *ConflictError {
	return &ConflictError{Message: message}
}

// ValidationError represents an error due to invalid input data.
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

// NewValidationError creates a new ValidationError.
func NewValidationError(message string) *ValidationError {
	return &ValidationError{Message: message}
}

// IsNotFoundError checks if an error is a NotFoundError.
func IsNotFoundError(err error) bool {
	var target *NotFoundError
	return errors.As(err, &target)
}

// IsConflictError checks if an error is a ConflictError.
func IsConflictError(err error) bool {
	var target *ConflictError
	return errors.As(err, &target)
}

// IsValidationError checks if an error is a ValidationError.
func IsValidationError(err error) bool {
	var target *ValidationError
	return errors.As(err, &target)
}
