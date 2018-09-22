package models

import (
	"fmt"
)

// AppError is a custom error struct to implement error interface
type AppError struct {
	Message       string `json:"message"`        // Message to be display to the end user without debugging information
	DetailedError string `json:"detailed_error"` // Internal error string to help the developer
	StatusCode    int    `json:"status_code"`    // The http status code
	Where         string `json:"where"`          // The function where it happened in the form of Struct.Func
}

// Error is a implemented method of AppError
func (ae *AppError) Error() string {
	return fmt.Sprintf("%s():\n %s", ae.Where, ae.DetailedError)
}

// NewAppError new a AppError
func NewAppError(where string, message string, details string, status int) *AppError {
	appError := &AppError{Message: message, Where: where, DetailedError: details, StatusCode: status}
	return appError
}
