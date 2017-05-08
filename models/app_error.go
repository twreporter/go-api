package models

// AppError is a custom error struct to implement error interface
type AppError struct {
	Message       string // Message to be display to the end user without debugging information
	DetailedError string // Internal error string to help the developer
	StatusCode    int    // The http status code
	Where         string // The function where it happened in the form of Struct.Func
}

// Error is a implemented method of AppError
func (er AppError) Error() string {
	return er.Where + ": " + er.Message + ", " + er.DetailedError
}

// NewAppError new a AppError
func NewAppError(where string, message string, details string, status int) AppError {
	ap := AppError{Message: message, Where: where, DetailedError: details, StatusCode: status}
	return ap
}
