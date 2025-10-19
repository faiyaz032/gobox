package errors

import (
	"errors"
	"fmt"
)

type AppError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Err     error  `json:"-"` // internal error (not sent to client)
}

func (e *AppError) Error() string {
	if e.Err != nil {
		return fmt.Sprintf("%s: %v", e.Message, e.Err)
	}
	return e.Message
}

func (e *AppError) Unwrap() error {
	return e.Err
}

func NewAppError(code int, message string) *AppError {
	return &AppError{Code: code, Message: message}
}

func Wrap(err error, code int, message string) *AppError {
	if err == nil {
		return nil
	}
	return &AppError{
		Code:    code,
		Message: message,
		Err:     err,
	}
}

func IsAppError(err error) bool {
	var appErr *AppError
	return errors.As(err, &appErr)
}
