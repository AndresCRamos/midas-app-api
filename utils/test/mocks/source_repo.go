package mocks

import (
	"github.com/AndresCRamos/midas-app-api/models"
	error_const "github.com/AndresCRamos/midas-app-api/utils/errors"
)

type SourceRepositoryMock struct{}

func (r SourceRepositoryMock) CreateNewSource(source models.Source) (models.Source, error) {
	wrapper := error_const.SourceRepositoryError{}
	switch source.Name {
	case "Success":
		return source, nil
	case "CantConnect":
		wrapper.Wrap(error_const.FirebaseUnknownError{})
		return models.Source{}, wrapper
	case "Duplicated":
		wrapper.Wrap(error_const.FirestoreAlreadyExistsError{DocID: source.UID})
		return models.Source{}, wrapper
	default:
		return models.Source{}, error_const.TestInvalidTestCaseError{Param: source.Name}
	}
}

var TestSource = models.Source{UID: "0", Name: "TEST_SOURCE"}

func (r SourceRepositoryMock) GetSourceByID(id string, user string) (models.Source, error) {
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
		wrapper.Wrap(error_const.FirestoreParsingError{DocID: id, StructName: "source"})
		return models.Source{}, wrapper
	default:
		return models.Source{}, error_const.TestInvalidTestCaseError{Param: id}
	}
}

func (r SourceRepositoryMock) UpdateSource(source models.Source) (models.Source, error) {
	id := source.UID
	wrapper := error_const.SourceRepositoryError{}
	switch id {
	case "0":
		return source, nil
	case "1":
		wrapper.Wrap(error_const.FirebaseUnknownError{})
		return models.Source{}, wrapper
	case "2":
		wrapper.Wrap(error_const.FirestoreNotFoundError{DocID: id})
		return models.Source{}, wrapper
	case "3":
		wrapper.Wrap(error_const.FirestoreParsingError{DocID: id, StructName: "source"})
		return models.Source{}, wrapper
	default:
		return models.Source{}, error_const.TestInvalidTestCaseError{Param: id}
	}
}

func (r SourceRepositoryMock) DeleteSource(id string) error {
	wrapper := error_const.SourceRepositoryError{}
	switch id {
	case "0":
		return nil
	case "1":
		wrapper.Wrap(error_const.FirebaseUnknownError{})
		return wrapper
	case "2":
		wrapper.Wrap(error_const.FirestoreNotFoundError{DocID: id})
		return wrapper
	case "3":
		wrapper.Wrap(error_const.FirestoreParsingError{DocID: id, StructName: "source"})
		return wrapper
	default:
		return error_const.TestInvalidTestCaseError{Param: id}
	}
}
