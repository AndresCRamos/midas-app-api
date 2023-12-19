package repository

import (
	"fmt"

	"firebase.google.com/go/v4/errorutils"
	"github.com/AndresCRamos/midas-app-api/models"
	error_const "github.com/AndresCRamos/midas-app-api/utils/errors"
)

func checkError(err error, id string, user models.User) error {
	if errorutils.IsNotFound(err) {
		return fmt.Errorf(error_const.FIRESTORE_NOT_FOUND, id)
	}
	if errorutils.IsUnauthenticated(err) {
		return error_const.UNAUTHENTICATED
	}
	if errorutils.IsInternal(err) {
		return error_const.INTERNAL_ERROR
	}

	if errorutils.IsResourceExhausted(err) {
		return error_const.MAX_QUOTA
	}

	return error_const.UNKNOWN
}
