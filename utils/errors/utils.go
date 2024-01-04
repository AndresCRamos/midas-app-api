package errors

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/AndresCRamos/midas-app-api/models"
	"github.com/go-playground/validator/v10"
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

func parseBindingErrors(errs ...error) []string {
	out := []string{}
	for _, err := range errs {
		switch typedErr := any(err).(type) {
		case validator.ValidationErrors:
			for _, fieldErr := range typedErr {
				out = append(out, parseValidationErr(fieldErr))
			}
		case *json.UnmarshalTypeError:
			out = append(out, parseJsonSyntaxError(typedErr))
		default:
			out = append(out, err.Error())
		}

	}
	return out
}

func parseValidationErr(err validator.FieldError) string {
	fieldName := strings.ToLower(err.Field())
	switch err.Tag() {
	case "required":
		return fmt.Sprintf("field %s is required", fieldName)
	case "required_without":
		return fmt.Sprintf("field %s is required if %s is not supplied", fieldName, strings.ToLower(err.Param()))
	default:
		errMsg := fmt.Sprintf("field %s failed validation %s", fieldName, err.Tag())
		if err.Param() != "" {
			errMsg += fmt.Sprintf(": %s", err.Param())
		}
		return errMsg
	}

}

func parseJsonSyntaxError(err *json.UnmarshalTypeError) string {
	return fmt.Sprintf("field %s should be %s", err.Field, err.Type.String())
}
