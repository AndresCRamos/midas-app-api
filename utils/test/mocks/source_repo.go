package mocks

import (
	"github.com/AndresCRamos/midas-app-api/models"
	error_const "github.com/AndresCRamos/midas-app-api/utils/errors"
)

type SourceRepositoryMock struct{}

func (r SourceRepositoryMock) CreateNewSource(user models.Source) error {
	wrapper := error_const.SourceRepositoryError{}
	switch user.Name {
	case "Success":
		return nil
	case "CantConnect":
		wrapper.Wrap(error_const.FirebaseUnknownError{})
		return wrapper
	case "Duplicated":
		wrapper.Wrap(error_const.FirestoreAlreadyExistsError{DocID: user.UID})
		return wrapper
	default:
		return error_const.TestInvalidTestCaseError{Param: user.Name}
	}
}

var TestSource = models.Source{UID: "0", Name: "TEST_SOURCE"}

func (r SourceRepositoryMock) GetSourceByID(id string) (models.Source, error) {
	wrapper := error_const.SourceRepositoryError{}
	switch id {
	case "0":
		return TestSource, nil
	case "1":
		wrapper.Wrap(error_const.FirebaseUnknownError{})
		return models.Source{}, wrapper
	case "2":
		wrapper.Wrap(error_const.FirestoreNotFoundError{DocID: id})
		return models.Source{}, wrapper
	case "3":
		wrapper.Wrap(error_const.FirestoreParsingError{DocID: id, StructName: "user"})
		return models.Source{}, wrapper
	default:
		return models.Source{}, error_const.TestInvalidTestCaseError{Param: id}
	}
}

func (r SourceRepositoryMock) UpdateSource(id string) (models.Source, error) {
	wrapper := error_const.SourceRepositoryError{}
	switch id {
	case "0":
		return TestSource, nil
	case "1":
		wrapper.Wrap(error_const.FirebaseUnknownError{})
		return models.Source{}, wrapper
	case "2":
		wrapper.Wrap(error_const.FirestoreNotFoundError{DocID: id})
		return models.Source{}, wrapper
	case "3":
		wrapper.Wrap(error_const.FirestoreParsingError{DocID: id, StructName: "user"})
		return models.Source{}, wrapper
	default:
		return models.Source{}, error_const.TestInvalidTestCaseError{Param: id}
	}
}

func (r SourceRepositoryMock) DeleteSource(id string) (models.Source, error) {
	wrapper := error_const.SourceRepositoryError{}
	switch id {
	case "0":
		return TestSource, nil
	case "1":
		wrapper.Wrap(error_const.FirebaseUnknownError{})
		return models.Source{}, wrapper
	case "2":
		wrapper.Wrap(error_const.FirestoreNotFoundError{DocID: id})
		return models.Source{}, wrapper
	case "3":
		wrapper.Wrap(error_const.FirestoreParsingError{DocID: id, StructName: "user"})
		return models.Source{}, wrapper
	default:
		return models.Source{}, error_const.TestInvalidTestCaseError{Param: id}
	}
}
