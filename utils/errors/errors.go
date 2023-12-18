package errors

import "errors"

const (
	FIREBASE_ERROR       = "Failed to initialize Firebase: %w"
	FIRESTORE_ERROR      = "Failed to initialize FireStore: %w"
	AUTH_ERROR           = "Failed to initialize Firebase Authentication: %w"
	INITIALIZE_APP_ERROR = "Failed to initialize: %w"
)

var (
	EMPTY_PROJECT              = errors.New("Firebase project cannot be empty")
	FIRESTORE_CANT_CONNECT     = errors.New("Can't connect to Firestore")
	FIREBASE_AUTH_CANT_CONNECT = errors.New("Can't connect to Firebase Auth")
)
