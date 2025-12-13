package geckocoin

import (
	"errors"
	"fmt"
)

var ErrCoinNotFound = errors.New("coin_id_not_found")

type APIError struct {
	Status  string
	Code    int
	Message string
}

func NewAPIError(status string, code int, message string) *APIError {
	return &APIError{Status: status, Code: code, Message: message}
}

func (e *APIError) Error() string {
	return fmt.Sprintf("HTTP %s: %d\nmessage:%s", e.Status, e.Code, e.Message)
}
