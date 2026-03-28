package models

import "net/http"

type AppError struct {
	Code       string
	Message    string
	StatusCode int
}

func (e *AppError) Error() string {
	return e.Message
}

var (
	ErrNotFound     = &AppError{Code: "NOT_FOUND", Message: "Resource not found", StatusCode: http.StatusNotFound}
	ErrUnauthorized = &AppError{Code: "UNAUTHORIZED", Message: "Authentication required", StatusCode: http.StatusUnauthorized}
	ErrForbidden    = &AppError{Code: "FORBIDDEN", Message: "Access denied", StatusCode: http.StatusForbidden}
	ErrBadRequest   = &AppError{Code: "BAD_REQUEST", Message: "Bad request", StatusCode: http.StatusBadRequest}
	ErrConflict     = &AppError{Code: "CONFLICT", Message: "Resource already exists", StatusCode: http.StatusConflict}
	ErrInternal     = &AppError{Code: "INTERNAL_ERROR", Message: "Internal server error", StatusCode: http.StatusInternalServerError}
	ErrTokenExpired = &AppError{Code: "TOKEN_EXPIRED", Message: "Access token expired", StatusCode: http.StatusUnauthorized}
	ErrTokenInvalid = &AppError{Code: "TOKEN_INVALID", Message: "Invalid token", StatusCode: http.StatusUnauthorized}
)
