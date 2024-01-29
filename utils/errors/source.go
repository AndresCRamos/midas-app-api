package errors

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

const (
	SOURCE_ALREADY_EXISTS = "A source with id %s already exists"
	SOURCE_NOT_FOUND      = "A source with id %s doesn't exists"
	OWNER_NOT_FOUND       = "The source %s cant be created, because owner %s doesn't exists"
	OWNER_CANT_CHANGE     = "The provided owner %s of source %s is not the current one"
	SOURCE_NOT_OWNER      = "The user %s is not the owner of source %s"
)

const (
	source_repository_error = "SourceRepository: %s"
	source_service_error    = "SourceService: %s: %s"
)

// SourceRepositoryError struct
type SourceRepositoryError struct {
	Err error
}

func (sre SourceRepositoryError) Error() string {
	return fmt.Sprintf(source_repository_error, sre.Err.Error())
}

func (sre *SourceRepositoryError) Wrap(err error) {
	sre.Err = err
}

func (sre SourceRepositoryError) Unwrap() error {
	return sre.Err
}

// SourceServiceError struct
type SourceServiceError struct {
	Method string
	Err    error
}

func (sse SourceServiceError) Error() string {
	MethodMsg := ""
	switch sse.Method {
	case "Create":
		MethodMsg = "Cant create"
	case "Retrieve":
		MethodMsg = "Cant retrieve"
	case "Update":
		MethodMsg = "Cant update"
	case "Delete":
		MethodMsg = "Cant delete"
	default:
		MethodMsg = "Unknown method"
	}
	return fmt.Sprintf(source_service_error, MethodMsg, sse.Err.Error())
}

func (sse *SourceServiceError) Wrap(err error) {
	sse.Err = err
}

func (sse SourceServiceError) Unwrap() error {
	return sse.Err
}

type SourceDuplicated struct {
	SourceID string
}

func (sd SourceDuplicated) GetAPIError() (int, gin.H) {
	return http.StatusBadRequest, gin.H{
		"error": fmt.Sprintf(SOURCE_ALREADY_EXISTS, sd.SourceID),
	}
}

func (sd SourceDuplicated) Error() string {
	return fmt.Sprintf(SOURCE_ALREADY_EXISTS, sd.SourceID)
}

type SourceNotFound struct {
	SourceID string
}

func (snf SourceNotFound) GetAPIError() (int, gin.H) {
	return http.StatusNotFound, gin.H{
		"error": fmt.Sprintf(SOURCE_NOT_FOUND, snf.SourceID),
	}
}

func (snf SourceNotFound) Error() string {
	return fmt.Sprintf(SOURCE_NOT_FOUND, snf.SourceID)
}

type SourceOwnerNotFound struct {
	SourceID string
	OwnerId  string
}

func (onf SourceOwnerNotFound) GetAPIError() (int, gin.H) {
	return http.StatusNotFound, gin.H{
		"error": fmt.Sprintf(OWNER_NOT_FOUND, onf.SourceID, onf.OwnerId),
	}
}

func (onf SourceOwnerNotFound) Error() string {
	return fmt.Sprintf(OWNER_NOT_FOUND, onf.SourceID, onf.OwnerId)
}

type SourceCantChangeOwner struct {
	SourceID string
	OwnerID  string
}

func (sco SourceCantChangeOwner) GetAPIError() (int, gin.H) {
	return http.StatusNotFound, gin.H{
		"error": fmt.Sprintf(OWNER_CANT_CHANGE, sco.OwnerID, sco.SourceID),
	}
}

func (sco SourceCantChangeOwner) Error() string {
	return fmt.Sprintf(OWNER_CANT_CHANGE, sco.OwnerID, sco.SourceID)
}

type SourceDifferentOwner struct {
	SourceID string
	OwnerID  string
}

func (sdo SourceDifferentOwner) GetAPIError() (int, gin.H) {
	return http.StatusNotFound, gin.H{
		"error": fmt.Sprintf(SOURCE_NOT_OWNER, sdo.OwnerID, sdo.SourceID),
	}
}

func (sdo SourceDifferentOwner) Error() string {
	return fmt.Sprintf(SOURCE_NOT_OWNER, sdo.OwnerID, sdo.SourceID)
}
