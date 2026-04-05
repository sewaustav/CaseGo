package apperrors

import (
	"errors"
	"fmt"
	"net/http"
)

// AppError — типизированная ошибка приложения с HTTP-кодом.
type AppError struct {
	Code    int    // HTTP status code
	Message string // Сообщение для клиента
	Err     error  // Оригинальная ошибка (для логов)
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

func NewNotFound(msg string, err error) *AppError {
	return &AppError{Code: http.StatusNotFound, Message: msg, Err: err}
}

func NewBadRequest(msg string, err error) *AppError {
	return &AppError{Code: http.StatusBadRequest, Message: msg, Err: err}
}

func NewForbidden(msg string, err error) *AppError {
	return &AppError{Code: http.StatusForbidden, Message: msg, Err: err}
}

func NewConflict(msg string, err error) *AppError {
	return &AppError{Code: http.StatusConflict, Message: msg, Err: err}
}

func NewInternal(msg string, err error) *AppError {
	return &AppError{Code: http.StatusInternalServerError, Message: msg, Err: err}
}

func NewUnauthorized(msg string, err error) *AppError {
	return &AppError{Code: http.StatusUnauthorized, Message: msg, Err: err}
}

// Sentinel errors для репозиторного уровня
var (
	ErrNotFound      = errors.New("resource not found")
	ErrAlreadyExists = errors.New("resource already exists")
)
