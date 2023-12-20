package errors

import (
	"fmt"

	"github.com/AndresCRamos/midas-app-api/models"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func CheckFirebaseError(err error, id string, user models.User) error {

	statusErrCode := status.Code(err)

	if statusErrCode == codes.NotFound {
		return fmt.Errorf(FIRESTORE_NOT_FOUND, id)
	}
	if statusErrCode == codes.Unauthenticated {
		return UNAUTHENTICATED
	}
	if statusErrCode == codes.Internal {
		return INTERNAL_ERROR
	}
	if statusErrCode == codes.ResourceExhausted {
		return MAX_QUOTA
	}
	if statusErrCode == codes.Unavailable {
		return UNAVAILABLE
	}
	if statusErrCode == codes.AlreadyExists {
		return fmt.Errorf(ALREADY_EXISTS, id)
	}

	return UNKNOWN
}
