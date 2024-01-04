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
		wrapper.Wrap(logged_err)
		return wrapper
	}
	if statusErrCode == codes.Unauthenticated {
		wrapper.Wrap(FirebaseUnauthorizedError{})
		return wrapper
	}
	if statusErrCode == codes.Internal {
		wrapper.Wrap(FirebaseInternalServerError{})
		return wrapper
	}
	if statusErrCode == codes.ResourceExhausted {
		wrapper.Wrap(FirebaseMaxQuotaError{})
		return wrapper
	}
	if statusErrCode == codes.Unavailable {
		wrapper.Wrap(FirebaseUnavailableError{})
		return wrapper
	}
	if statusErrCode == codes.AlreadyExists {
		logged_err := FirestoreAlreadyExistsError{DocID: id}
		wrapper.Wrap(logged_err)
		return wrapper
	}

	wrapper.Wrap(FirebaseUnknownError{})
	return wrapper
}
