package errors

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

const (
	SOURCE_ALREADY_EXISTS = "A source with id %s already exists"
	SOURCE_NOT_FOUND      = "A source with id %s doesn't exists"
)

const (
	source_repository_error = "SourceRepository: %s"
	source_service_error    = "SourceService: %s: %s"
)

// SourceRepositoryError struct
type SourceRepositoryError struct {
	Err error
}

func (ure SourceRepositoryError) Error() string {
	return fmt.Sprintf(source_repository_error, ure.Err.Error())
}

func (ure *SourceRepositoryError) Wrap(err error) {
	ure.Err = err
}

func (ure SourceRepositoryError) Unwrap() error {
	return ure.Err
}

// SourceServiceError struct
type SourceServiceError struct {
	Method string
	Err    error
}

func (use SourceServiceError) Error() string {
	MethodMsg := ""
	switch use.Method {
	case "Create":
		MethodMsg = "Cant create"
	case "Retrieve":
		MethodMsg = "Cant retrieve"
	default:
		MethodMsg = "Unknown method"
	}
	return fmt.Sprintf(source_service_error, MethodMsg, use.Err.Error())
}

func (use *SourceServiceError) Wrap(err error) {
	use.Err = err
}

func (use SourceServiceError) Unwrap() error {
	return use.Err
}

type SourceDuplicated struct {
	SourceID string
}

func (ud SourceDuplicated) GetAPIError() (int, gin.H) {
	return http.StatusBadRequest, gin.H{
		"error": fmt.Sprintf(SOURCE_ALREADY_EXISTS, ud.SourceID),
	}
}

func (ud SourceDuplicated) Error() string {
	return fmt.Sprintf(SOURCE_ALREADY_EXISTS, ud.SourceID)
}

type SourceNotFound struct {
	SourceID string
}

func (ud SourceNotFound) GetAPIError() (int, gin.H) {
	return http.StatusNotFound, gin.H{
		"error": fmt.Sprintf(SOURCE_NOT_FOUND, ud.SourceID),
	}
}

func (ud SourceNotFound) Error() string {
	return fmt.Sprintf(SOURCE_NOT_FOUND, ud.SourceID)
}
