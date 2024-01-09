package errors

import "fmt"

const (
	firestore_error          = "Failed to initialize FireStore: %s"
	firestore_not_found      = "Cant find the specified document: %s"
	firestore_already_exists = "Document %s already exists"
	firestore_parsing_error  = "Cant parse document %s into struct %s"
)

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

func (fs FirestoreError) Unwrap() error {
	return fs.Err
}

// Can't connect to Firestore
type FirestoreCantConnect struct{}

func (fcc FirestoreCantConnect) Error() string {
	return "Can't connect to Firestore"
}

func (fcc FirestoreCantConnect) Unwrap() error {
	return nil
}

// FirestoreNotFoundError struct
type FirestoreNotFoundError struct {
	DocID string
}

func (fnf FirestoreNotFoundError) Error() string {
	return fmt.Sprintf(firestore_not_found, fnf.DocID)
}

func (fnf FirestoreNotFoundError) Unwrap() error {
	return nil
}

// FirestoreAlreadyExistsError struct
type FirestoreAlreadyExistsError struct {
	DocID string
}

func (aee FirestoreAlreadyExistsError) Error() string {
	return fmt.Sprintf(firestore_already_exists, aee.DocID)
}

func (aee FirestoreAlreadyExistsError) Unwrap() error {
	return nil
}

// FirestoreParsingError struct
type FirestoreParsingError struct {
	DocID      string
	StructName string
}

func (pe FirestoreParsingError) Error() string {
	return fmt.Sprintf(firestore_parsing_error, pe.DocID, pe.StructName)
}

func (pe FirestoreParsingError) Unwrap() error {
	return nil
}
