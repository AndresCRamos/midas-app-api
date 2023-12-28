package errors

import (
	"errors"
	"fmt"
)

type ErrorWrapper interface {
	Error() string
	Wrap(err error)
}

const (
	firebase_error        = "Failed to initialize Firebase: %s"
	firestore_error       = "Failed to initialize FireStore: %s"
	auth_error            = "Failed to initialize Firebase Authentication: %s"
	initialize_app_error  = "Failed to initialize: %s"
	firestore_not_found   = "Cant find the specified document: %s"
	already_exists        = "Document %s already exists"
	invalid_test_case     = "Parameter %v is not a valid test case"
	parsing_error         = "Cant parse document %s into struct %s"
	user_repository_error = "UserRepository: %s"
	user_service_error    = "UserService: %s"
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

// FirestoreError struct
type FirestoreError struct {
	Err error
}

func (fs FirestoreError) Error() string {
	return fmt.Sprintf(firestore_error, fs.Err.Error())
}

func (fs *FirestoreError) Wrap(err error) {
	fs.Err = err
}

func (fs *FirestoreError) Unwrap() error {
	return fs.Err
}

// AuthError struct
type AuthError struct {
	Err error
}

func (ae AuthError) Error() string {
	return fmt.Sprintf(auth_error, ae.Err.Error())
}

func (ae *AuthError) Wrap(err error) {
	ae.Err = err
}

func (ae *AuthError) Unwrap() error {
	return ae.Err
}

// InitializeAppError struct
type InitializeAppError struct {
	Err error
}

func (iae InitializeAppError) Error() string {
	return fmt.Sprintf(initialize_app_error, iae.Err.Error())
}

func (iae *InitializeAppError) Wrap(err error) {
	iae.Err = err
}

func (iae *InitializeAppError) Unwrap() error {
	return iae.Err
}

// FirestoreNotFoundError struct
type FirestoreNotFoundError struct {
	DocID string
}

func (fnf FirestoreNotFoundError) Error() string {
	return fmt.Sprintf(firestore_not_found, fnf.DocID)
}

func (fnf *FirestoreNotFoundError) Unwrap() error {
	return nil
}

// AlreadyExistsError struct
type AlreadyExistsError struct {
	DocID string
}

func (aee AlreadyExistsError) Error() string {
	return fmt.Sprintf(already_exists, aee.DocID)
}

func (aee *AlreadyExistsError) Unwrap() error {
	return nil
}

// InvalidTestCaseError struct
type InvalidTestCaseError struct {
	Param interface{}
}

func (tce InvalidTestCaseError) Error() string {
	return fmt.Sprintf(invalid_test_case, tce.Param)
}

func (tce *InvalidTestCaseError) Unwrap() error {
	return nil
}

// ParsingError struct
type ParsingError struct {
	DocID      string
	StructName string
}

func (pe ParsingError) Error() string {
	return fmt.Sprintf(parsing_error, pe.DocID, pe.StructName)
}

func (pe *ParsingError) Unwrap() error {
	return nil
}

// UserRepositoryError struct
type UserRepositoryError struct {
	Err error
}

func (ure UserRepositoryError) Error() string {
	return fmt.Sprintf(user_repository_error, ure.Err.Error())
}

func (ure *UserRepositoryError) Wrap(err error) {
	ure.Err = err
}

func (ure *UserRepositoryError) Unwrap() error {
	return ure.Err
}

// UserServiceError struct
type UserServiceError struct {
	Err error
}

func (use UserServiceError) Error() string {
	return fmt.Sprintf(user_service_error, use.Err.Error())
}

func (use *UserServiceError) Wrap(err error) {
	use.Err = err
}

func (use *UserServiceError) Unwrap() error {
	return use.Err
}

var (
	EMPTY_PROJECT              = errors.New("Firebase project cannot be empty")
	FIRESTORE_CANT_CONNECT     = errors.New("Can't connect to Firestore")
	FIREBASE_AUTH_CANT_CONNECT = errors.New("Can't connect to Firebase Auth")
	UNAUTHENTICATED            = errors.New("Unauthenticated")
	UNKNOWN                    = errors.New("Unknown Error")
	INTERNAL_ERROR             = errors.New("Internal Server error")
	UNAVAILABLE                = errors.New("Firestore is unavailable")
	DATA_LOSS                  = errors.New("Data got corrupted, try again")
	MAX_QUOTA                  = errors.New("Firebase max quota reached")
	MAP_INTERFACE_NOT_FOUND    = errors.New("Cant find value")
	MAP_INTERFACE_CANT_ASSERT  = errors.New("Cant assert value")
)

type APIError struct {
	Error string `json:"error"`
}
