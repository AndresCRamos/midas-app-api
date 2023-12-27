package errors

import "errors"

const (
	FIREBASE_ERROR        = "Failed to initialize Firebase: %w"
	FIRESTORE_ERROR       = "Failed to initialize FireStore: %w"
	AUTH_ERROR            = "Failed to initialize Firebase Authentication: %w"
	INITIALIZE_APP_ERROR  = "Failed to initialize: %w"
	FIRESTORE_NOT_FOUND   = "Cant find the specified document: %s"
	ALREADY_EXISTS        = "Document %s already exists"
	INVALID_TEST_CASE     = "Parameter %v is not a valid test case"
	PARSING_ERROR         = "Cant parse document %s into struct %s"
	USER_REPOSITORY_ERROR = "UserRepository: %w"
	USER_SERVICE_ERROR    = "UserService: %w"
)

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
)
