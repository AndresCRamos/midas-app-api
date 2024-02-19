package errors

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

const (
	MOVEMENT_ALREADY_EXISTS    = "A movement with id %s already exists"
	MOVEMENT_NOT_FOUND         = "A movement with id %s doesn't exists"
	MOVEMENT_OWNER_NOT_FOUND   = "The movement %s cant be created, because owner %s doesn't exists"
	MOVEMENT_SOURCE_NOT_FOUND  = "The movement %s cant be created, because source %s doesn't exists"
	MOVEMENT_OWNER_CANT_CHANGE = "The provided owner %s of movement %s is not the current one"
	MOVEMENT_NOT_OWNER         = "The user %s is not the owner of movement %s"
	MOVEMENT_NOT_ENOUGH_DATA   = "Requested movement page does not exists"
	MOVEMENT_BAD_DATES         = "Date from must be less that date to"
)

const (
	movement_repository_error = "MovementRepository: %s"
	movement_service_error    = "MovementService: %s: %s"
)

// MovementRepositoryError struct
type MovementRepositoryError struct {
	Err error
}

func (sre MovementRepositoryError) Error() string {
	return fmt.Sprintf(movement_repository_error, sre.Err.Error())
}

func (sre *MovementRepositoryError) Wrap(err error) {
	sre.Err = err
}

func (sre MovementRepositoryError) Unwrap() error {
	return sre.Err
}

// MovementServiceError struct
type MovementServiceError struct {
	Method string
	Err    error
}

func (sse MovementServiceError) Error() string {
	MethodMsg := ""
	switch sse.Method {
	case "Create":
		MethodMsg = "Cant create"
	case "Retrieve":
		MethodMsg = "Cant retrieve"
	case "List":
		MethodMsg = "Cant get list"
	case "Update":
		MethodMsg = "Cant update"
	case "Delete":
		MethodMsg = "Cant delete"
	default:
		MethodMsg = "Unknown method"
	}
	return fmt.Sprintf(movement_service_error, MethodMsg, sse.Err.Error())
}

func (sse *MovementServiceError) Wrap(err error) {
	sse.Err = err
}

func (sse MovementServiceError) Unwrap() error {
	return sse.Err
}

type MovementDuplicated struct {
	MovementID string
}

func (sd MovementDuplicated) GetAPIError() (int, gin.H) {
	return http.StatusBadRequest, gin.H{
		"error": fmt.Sprintf(MOVEMENT_ALREADY_EXISTS, sd.MovementID),
	}
}

func (sd MovementDuplicated) Error() string {
	return fmt.Sprintf(MOVEMENT_ALREADY_EXISTS, sd.MovementID)
}

type MovementNotFound struct {
	MovementID string
}

func (snf MovementNotFound) GetAPIError() (int, gin.H) {
	return http.StatusNotFound, gin.H{
		"error": fmt.Sprintf(MOVEMENT_NOT_FOUND, snf.MovementID),
	}
}

func (snf MovementNotFound) Error() string {
	return fmt.Sprintf(MOVEMENT_NOT_FOUND, snf.MovementID)
}

type MovementOwnerNotFound struct {
	MovementID string
	OwnerId    string
}

func (onf MovementOwnerNotFound) GetAPIError() (int, gin.H) {
	return http.StatusNotFound, gin.H{
		"error": fmt.Sprintf(OWNER_NOT_FOUND, onf.MovementID, onf.OwnerId),
	}
}

func (onf MovementOwnerNotFound) Error() string {
	return fmt.Sprintf(OWNER_NOT_FOUND, onf.MovementID, onf.OwnerId)
}

type MovementCantChangeOwner struct {
	MovementID string
	OwnerID    string
}

func (sco MovementCantChangeOwner) GetAPIError() (int, gin.H) {
	return http.StatusNotFound, gin.H{
		"error": fmt.Sprintf(OWNER_CANT_CHANGE, sco.OwnerID, sco.MovementID),
	}
}

func (sco MovementCantChangeOwner) Error() string {
	return fmt.Sprintf(OWNER_CANT_CHANGE, sco.OwnerID, sco.MovementID)
}

type MovementDifferentOwner struct {
	MovementID string
	OwnerID    string
}

func (sdo MovementDifferentOwner) GetAPIError() (int, gin.H) {
	return http.StatusNotFound, gin.H{
		"error": fmt.Sprintf(MOVEMENT_NOT_OWNER, sdo.OwnerID, sdo.MovementID),
	}
}

func (sdo MovementDifferentOwner) Error() string {
	return fmt.Sprintf(MOVEMENT_NOT_OWNER, sdo.OwnerID, sdo.MovementID)
}

type MovementSourceNotFound struct {
	MovementID string
	SourceID   string
}

func (smn MovementSourceNotFound) GetAPIError() (int, gin.H) {
	return http.StatusNotFound, gin.H{
		"error": fmt.Sprintf(MOVEMENT_SOURCE_NOT_FOUND, smn.MovementID, smn.SourceID),
	}
}

func (smn MovementSourceNotFound) Error() string {
	return fmt.Sprintf(MOVEMENT_SOURCE_NOT_FOUND, smn.MovementID, smn.SourceID)
}

type MovementNotEnoughData struct{}

func (sne MovementNotEnoughData) GetAPIError() (int, gin.H) {
	return http.StatusNotFound, gin.H{
		"error": MOVEMENT_NOT_ENOUGH_DATA,
	}
}

func (sne MovementNotEnoughData) Error() string {
	return MOVEMENT_NOT_ENOUGH_DATA
}

type MovementBadDates struct{}

func (sne MovementBadDates) GetAPIError() (int, gin.H) {
	return http.StatusNotFound, gin.H{
		"error": MOVEMENT_BAD_DATES,
	}
}

func (sne MovementBadDates) Error() string {
	return MOVEMENT_BAD_DATES
}
