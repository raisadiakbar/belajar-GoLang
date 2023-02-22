package errors

import (
	"errors"
	"fmt"
)

type APIError struct {
	Code    int    `json:"-"`
	Message string `json:"message"`
}

func (e *APIError) Error() string {
	return fmt.Sprintf("%d: %s", e.Code, e.Message)
}

var (
	ErrNotFound     = &APIError{404, "Not Found"}
	ErrBadRequest   = &APIError{400, "Bad Request"}
	ErrUnauthorized = &APIError{401, "Unauthorized"}
)

func NewAPIError(code int, message string) error {
	return &APIError{Code: code, Message: message}
}

func Wrap(err error, message string) error {
	if err == nil {
		return nil
	}
	return fmt.Errorf("%s: %w", message, err)
}

func Is(err, target error) bool {
	if e, ok := err.(*APIError); ok {
		if t, ok := target.(*APIError); ok {
			return e.Code == t.Code
		}
	}
	return errors.Is(err, target)
}
