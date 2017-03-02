package models

// AppError is a custom error struct to implement error interface
type AppError struct {
	ID            string
	Message       string // Message to be display to the end user without debugging information
	DetailedError string // Internal error string to help the developer
	StatusCode    int    // The http status code
	Where         string // The function where it happened in the form of Struct.Func
}

// Error is a implemented method of AppError
func (er *AppError) Error() string {
	return er.Where + ": " + er.Message + ", " + er.DetailedError
}

// NewAppError new a AppError
func NewAppError(where string, id string, details string, status int) *AppError {
	ap := &AppError{}
	ap.ID = id
	ap.Message = id
	ap.Where = where
	ap.DetailedError = details
	ap.StatusCode = status
	return ap
}
