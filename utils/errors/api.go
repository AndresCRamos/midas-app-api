package errors

import "fmt"

const (
	API_UNKNOWN       = "An error has ocurred, try again later"
	API_UNAUTHORIZED  = "Failed to authenticate user"
	USER_INVALID_BODY = "Invalid request format"
)

const (
	api_initialize_app_error = "Failed to initialize: %s"
)

type APIError struct {
	Error string `json:"error"`
}

type ErrorWrapper interface {
	Error() string
	Wrap(err error)
}

// InitializeAppError struct
type InitializeAppError struct {
	Err error
}

func (iae InitializeAppError) Error() string {
	return fmt.Sprintf(api_initialize_app_error, iae.Err.Error())
}

func (iae *InitializeAppError) Wrap(err error) {
	iae.Err = err
}

func (iae *InitializeAppError) Unwrap() error {
	return iae.Err
}
