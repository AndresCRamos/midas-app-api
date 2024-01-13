package errors

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"strings"

	"github.com/go-playground/validator/v10"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func CheckFirebaseError(err error, id string, wrapper ErrorWrapper) error {

	statusErrCode := status.Code(err)

	if CheckFirebaseNotFound(err) {
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

func CheckFirebaseNotFound(err error) bool {
	return status.Code(err) == codes.NotFound
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
	case "depends_on":
		return fmt.Sprintf("field %s depends on %s, which is not supplied", fieldName, strings.ToLower(err.Param()))
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

func CheckServiceErrors(id string, err error, typeName string) APIError {

	alreadyExists := &FirestoreAlreadyExistsError{}
	unauthorized := &FirebaseUnauthorizedError{}
	notFound := &FirestoreNotFoundError{}

	if errors.As(err, unauthorized) {
		return APIUnauthorized{}
	}
	if errors.As(err, alreadyExists) {
		return getAlreadyExistsByType(typeName, id)
	}
	if errors.As(err, notFound) {
		return getNotFoundByType(typeName, id)
	}
	log.Print(err)
	return APIUnknown{}
}

func getAlreadyExistsByType(typeName string, id string) APIError {
	switch typeName {
	case "user":
		return UserDuplicated{UserID: id}

	case "source":
		return SourceDuplicated{SourceID: id}
	}
	return APIUnknown{}
}

func getNotFoundByType(typeName string, id string) APIError {
	switch typeName {
	case "user":
		return UserNotFound{UserID: id}
	case "source":
		return SourceDuplicated{SourceID: id}
	}
	return APIUnknown{}
}
