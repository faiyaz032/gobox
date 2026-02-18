package domain

import (
	"errors"
	"fmt"
	"net/http"
)

// ErrorType represents the category of error
type ErrorType string

const (
	ErrorTypeNotFound     ErrorType = "NOT_FOUND"
	ErrorTypeValidation   ErrorType = "VALIDATION"
	ErrorTypeInternal     ErrorType = "INTERNAL"
	ErrorTypeDatabase     ErrorType = "DATABASE"
	ErrorTypeDocker       ErrorType = "DOCKER"
	ErrorTypeConflict     ErrorType = "CONFLICT"
	ErrorTypeUnauthorized ErrorType = "UNAUTHORIZED"
)

// AppError represents a structured application error
type AppError struct {
	Type       ErrorType
	Message    string
	StatusCode int
	Err        error // underlying error
}

// Error implements the error interface
func (e *AppError) Error() string {
	if e.Err != nil {
		return fmt.Sprintf("%s: %s (caused by: %v)", e.Type, e.Message, e.Err)
	}
	return fmt.Sprintf("%s: %s", e.Type, e.Message)
}

// Unwrap implements the errors.Unwrap interface
func (e *AppError) Unwrap() error {
	return e.Err
}

// IsType checks if the error is of a specific type
func (e *AppError) IsType(errType ErrorType) bool {
	return e.Type == errType
}

// Constructors for common error types

// NewNotFoundError creates a not found error
func NewNotFoundError(resource string, identifier string) *AppError {
	return &AppError{
		Type:       ErrorTypeNotFound,
		Message:    fmt.Sprintf("%s not found: %s", resource, identifier),
		StatusCode: http.StatusNotFound,
	}
}

// NewValidationError creates a validation error
func NewValidationError(message string) *AppError {
	return &AppError{
		Type:       ErrorTypeValidation,
		Message:    message,
		StatusCode: http.StatusBadRequest,
	}
}

// NewInternalError creates an internal server error
func NewInternalError(message string, err error) *AppError {
	return &AppError{
		Type:       ErrorTypeInternal,
		Message:    message,
		StatusCode: http.StatusInternalServerError,
		Err:        err,
	}
}

// NewDatabaseError creates a database error
func NewDatabaseError(operation string, err error) *AppError {
	return &AppError{
		Type:       ErrorTypeDatabase,
		Message:    fmt.Sprintf("database operation failed: %s", operation),
		StatusCode: http.StatusInternalServerError,
		Err:        err,
	}
}

// NewDockerError creates a docker error
func NewDockerError(operation string, err error) *AppError {
	return &AppError{
		Type:       ErrorTypeDocker,
		Message:    fmt.Sprintf("docker operation failed: %s", operation),
		StatusCode: http.StatusInternalServerError,
		Err:        err,
	}
}

// NewConflictError creates a conflict error
func NewConflictError(message string) *AppError {
	return &AppError{
		Type:       ErrorTypeConflict,
		Message:    message,
		StatusCode: http.StatusConflict,
	}
}

// NewUnauthorizedError creates an unauthorized error
func NewUnauthorizedError(message string) *AppError {
	return &AppError{
		Type:       ErrorTypeUnauthorized,
		Message:    message,
		StatusCode: http.StatusUnauthorized,
	}
}

// Helper functions

// IsAppError checks if an error is an AppError
func IsAppError(err error) (*AppError, bool) {
	var appErr *AppError
	if errors.As(err, &appErr) {
		return appErr, true
	}
	return nil, false
}

// GetStatusCode extracts HTTP status code from error
func GetStatusCode(err error) int {
	if appErr, ok := IsAppError(err); ok {
		return appErr.StatusCode
	}
	return http.StatusInternalServerError
}

// GetErrorType extracts error type from error
func GetErrorType(err error) ErrorType {
	if appErr, ok := IsAppError(err); ok {
		return appErr.Type
	}
	return ErrorTypeInternal
}

// GetErrorMessage extracts user-friendly message from error
func GetErrorMessage(err error) string {
	if appErr, ok := IsAppError(err); ok {
		return appErr.Message
	}
	return "An unexpected error occurred"
}
