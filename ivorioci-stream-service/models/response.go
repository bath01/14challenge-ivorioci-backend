package models

import "time"

type SuccessResponse struct {
	Success   bool        `json:"success"`
	Data      interface{} `json:"data"`
	Timestamp string      `json:"timestamp"`
}

type ErrorDetail struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

type ErrorResponse struct {
	Success   bool        `json:"success"`
	Error     ErrorDetail `json:"error"`
	Timestamp string      `json:"timestamp"`
}

type PaginatedData struct {
	Items      interface{} `json:"items"`
	Total      int         `json:"total"`
	Page       int         `json:"page"`
	Limit      int         `json:"limit"`
	TotalPages int         `json:"totalPages"`
}

func NewSuccess(data interface{}) SuccessResponse {
	return SuccessResponse{
		Success:   true,
		Data:      data,
		Timestamp: time.Now().UTC().Format(time.RFC3339),
	}
}

func NewError(code, message string) ErrorResponse {
	return ErrorResponse{
		Success:   false,
		Error:     ErrorDetail{Code: code, Message: message},
		Timestamp: time.Now().UTC().Format(time.RFC3339),
	}
}
