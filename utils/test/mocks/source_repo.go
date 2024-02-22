package mocks

import (
	"time"

	"github.com/AndresCRamos/midas-app-api/models"
	util_models "github.com/AndresCRamos/midas-app-api/utils/api/models"
	error_const "github.com/AndresCRamos/midas-app-api/utils/errors"
)

type SourceRepositoryMock struct{}

func (r SourceRepositoryMock) GetMovementsBySourceAndDate(id string, userID string, page int, date_from time.Time, date_to time.Time) (util_models.PaginatedSearch[models.Movement], error) {
	wrapper := error_const.SourceRepositoryError{}
	switch userID {
	case "0":
		return util_models.PaginatedSearch[models.Movement]{
			CurrentPage: 1,
			TotalData:   1,
			PageSize:    1,
			Data: []models.Movement{
				TestMovement,
			},
		}, nil
	case "1":
		wrapper.Wrap(error_const.FirebaseUnknownError{})
		return util_models.PaginatedSearch[models.Movement]{}, wrapper
	case "2":
		wrapper.Wrap(error_const.SourceNotFound{})
		return util_models.PaginatedSearch[models.Movement]{}, wrapper
	case "3":
		wrapper.Wrap(error_const.SourceDifferentOwner{})
		return util_models.PaginatedSearch[models.Movement]{}, wrapper
	case "4":
		wrapper.Wrap(error_const.MovementNotEnoughData{})
		return util_models.PaginatedSearch[models.Movement]{}, wrapper
	default:
		return util_models.PaginatedSearch[models.Movement]{}, error_const.TestInvalidTestCaseError{Param: userID}
	}
}

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
var TestSourceRetrieve = models.SourceRetrieve{UID: "0", Name: "TEST_SOURCE"}

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

func (r SourceRepositoryMock) GetSourcesByUser(userID string, page int) (util_models.PaginatedSearch[models.Source], error) {
	wrapper := error_const.SourceRepositoryError{}
	switch userID {
	case "0":
		return util_models.PaginatedSearch[models.Source]{
			CurrentPage: page,
			TotalData:   1,
			PageSize:    1,
			Data: []models.Source{
				TestSource,
			},
		}, nil
	case "1":
		wrapper.Wrap(error_const.FirebaseUnknownError{})
		return util_models.PaginatedSearch[models.Source]{}, wrapper
	case "2":
		wrapper.Wrap(error_const.SourceNotEnoughData{})
		return util_models.PaginatedSearch[models.Source]{}, wrapper
	default:
		return util_models.PaginatedSearch[models.Source]{}, error_const.TestInvalidTestCaseError{Param: userID}
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

func (r SourceRepositoryMock) DeleteSource(id string, userID string) error {
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
