package errors

import "fmt"

const (
	firebase_auth_error = "Failed to initialize Firebase Authentication: %s"
)

// FirebaseAuthError struct
type FirebaseAuthError struct {
	Err error
}

func (ae FirebaseAuthError) Error() string {
	return fmt.Sprintf(firebase_auth_error, ae.Err.Error())
}

func (ae *FirebaseAuthError) Wrap(err error) {
	ae.Err = err
}

func (ae FirebaseAuthError) Unwrap() error {
	return ae.Err
}

// Can't connect to Firebase Auth
type FirebaseAuthCantConnect struct{}

func (fac FirebaseAuthCantConnect) Error() string {
	return "Can't connect to Firebase Auth"
}

func (fac FirebaseAuthCantConnect) Unwrap() error {
	return nil
}
