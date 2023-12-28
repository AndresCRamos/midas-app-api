package errors

import (
	"github.com/AndresCRamos/midas-app-api/models"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func CheckFirebaseError(err error, id string, user models.User, wrapper ErrorWrapper) error {

	statusErrCode := status.Code(err)

	if statusErrCode == codes.NotFound {
		logged_err := FirestoreNotFoundError{DocID: id}
		wrapper.Wrap(&logged_err)
		return wrapper
	}
	if statusErrCode == codes.Unauthenticated {
		wrapper.Wrap(UNAUTHENTICATED)
		return wrapper
	}
	if statusErrCode == codes.Internal {
		wrapper.Wrap(INTERNAL_ERROR)
		return wrapper
	}
	if statusErrCode == codes.ResourceExhausted {
		wrapper.Wrap(MAX_QUOTA)
		return wrapper
	}
	if statusErrCode == codes.Unavailable {
		wrapper.Wrap(UNAVAILABLE)
		return wrapper
	}
	if statusErrCode == codes.AlreadyExists {
		logged_err := AlreadyExistsError{DocID: id}
		wrapper.Wrap(&logged_err)
		return wrapper
	}

	wrapper.Wrap(UNKNOWN)
	return wrapper
}
