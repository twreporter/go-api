package models

type AppError struct {
	Id            string
	Message       string // Message to be display to the end user without debugging information
	DetailedError string // Internal error string to help the developer
	StatusCode    int    // The http status code
	Where         string // The function where it happened in the form of Struct.Func
}

func (er AppError) Error() string {
	return er.Where + ": " + er.Message + ", " + er.DetailedError
}

func NewAppError(where string, id string, details string, status int) AppError {
	ap := AppError{}
	ap.Id = id
	ap.Message = id
	ap.Where = where
	ap.DetailedError = details
	ap.StatusCode = status
	return ap
}
