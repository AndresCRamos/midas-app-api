package errors

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

const (
	api_unknown              = "An error has ocurred, try again later"
	api_unauthorized         = "Failed to authenticate user"
	api_initialize_app_error = "Failed to initialize: %s"
	request_invalid_body     = "Invalid request format"
)

type APIError interface {
	GetAPIError() (int, gin.H)
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

type APIUnknown struct{}

func (au APIUnknown) GetAPIError() (int, gin.H) {
	return http.StatusBadRequest, gin.H{
		"error": api_unknown,
	}
}

func (au APIUnknown) Error() string {
	return api_unknown
}

type APIUnauthorized struct{}

func (au APIUnauthorized) GetAPIError() (int, gin.H) {
	return http.StatusBadRequest, gin.H{
		"error": api_unauthorized,
	}
}

type APIInvalidRequestBody struct{}

func (irb APIInvalidRequestBody) GetAPIError() (int, gin.H) {
	return http.StatusBadRequest, gin.H{
		"error": request_invalid_body,
	}
}

func (irb APIInvalidRequestBody) Error() string {
	return request_invalid_body
}
