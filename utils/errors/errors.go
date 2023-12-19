package errors

import (
	"fmt"

	"firebase.google.com/go/v4/errorutils"
	"github.com/AndresCRamos/midas-app-api/models"
)

func CheckFirebaseError(err error, id string, user models.User) error {
	if errorutils.IsNotFound(err) {
		return fmt.Errorf(FIRESTORE_NOT_FOUND, id)
	}
	if errorutils.IsUnauthenticated(err) {
		return UNAUTHENTICATED
	}
	if errorutils.IsInternal(err) {
		return INTERNAL_ERROR
	}

	if errorutils.IsResourceExhausted(err) {
		return MAX_QUOTA
	}

	return UNKNOWN
}
