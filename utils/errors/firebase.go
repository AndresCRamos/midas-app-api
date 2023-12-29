package errors

import "fmt"

const (
	firebase_error = "Failed to initialize Firebase: %s"
)

// FirebaseError struct
type FirebaseError struct {
	Err error
}

func (fb FirebaseError) Error() string {
	return fmt.Sprintf(firebase_error, fb.Err.Error())
}

func (fb *FirebaseError) Wrap(err error) {
	fb.Err = err
}

func (fb *FirebaseError) Unwrap() error {
	return fb.Err
}

// Firebase project can not be empty
type FirebaseEmptyProject struct{}

func (ep FirebaseEmptyProject) Error() string {
	return "Firebase project cannot be empty"
}

func (ep FirebaseEmptyProject) Unwrap() error {
	return nil
}

// Unauthenticated
type FirebaseUnauthorizedError struct{}

func (ue FirebaseUnauthorizedError) Error() string {
	return "Unauthenticated"
}

func (ue FirebaseUnauthorizedError) Unwrap() error {
	return nil
}

// Unknown Error
type FirebaseUnknownError struct{}

func (ue FirebaseUnknownError) Error() string {
	return "Unknown Error"
}

func (ue FirebaseUnknownError) Unwrap() error {
	return nil
}

// Internal Server error
type FirebaseInternalServerError struct{}

func (ise FirebaseInternalServerError) Error() string {
	return "Internal Server error"
}

func (ise FirebaseInternalServerError) Unwrap() error {
	return nil
}

// Firestore is unavailable
type FirebaseUnavailableError struct{}

func (ue FirebaseUnavailableError) Error() string {
	return "Firestore is unavailable"
}

func (ue FirebaseUnavailableError) Unwrap() error {
	return nil
}

// Data got corrupted, try again
type FirebaseDataLossError struct{}

func (dle FirebaseDataLossError) Error() string {
	return "Data got corrupted, try again"
}

func (dle FirebaseDataLossError) Unwrap() error {
	return nil
}

// Firebase max quota reached
type FirebaseMaxQuotaError struct{}

func (mqe FirebaseMaxQuotaError) Error() string {
	return "Firebase max quota reached"
}

func (mqe FirebaseMaxQuotaError) Unwrap() error {
	return nil
}
