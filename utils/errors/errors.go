package errors

import (
	"fmt"

	"github.com/AndresCRamos/midas-app-api/models"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func CheckFirebaseError(err error, id string, user models.User, wrapper string) error {

	statusErrCode := status.Code(err)

	if statusErrCode == codes.NotFound {
		logged_err := fmt.Errorf(FIRESTORE_NOT_FOUND, id)
		return fmt.Errorf(wrapper, logged_err)
	}
	if statusErrCode == codes.Unauthenticated {
		return fmt.Errorf(wrapper, UNAUTHENTICATED)
	}
	if statusErrCode == codes.Internal {
		return fmt.Errorf(wrapper, INTERNAL_ERROR)
	}
	if statusErrCode == codes.ResourceExhausted {
		return fmt.Errorf(wrapper, MAX_QUOTA)
	}
	if statusErrCode == codes.Unavailable {
		return fmt.Errorf(wrapper, UNAVAILABLE)
	}
	if statusErrCode == codes.AlreadyExists {
		logged_err := fmt.Errorf(ALREADY_EXISTS, id)
		return fmt.Errorf(wrapper, logged_err)
	}

	return fmt.Errorf(wrapper, UNKNOWN)
}
