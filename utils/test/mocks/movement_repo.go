package mocks

import (
	"time"

	"github.com/AndresCRamos/midas-app-api/models"
	util_models "github.com/AndresCRamos/midas-app-api/utils/api/models"
	error_const "github.com/AndresCRamos/midas-app-api/utils/errors"
)

type MovementRepositoryMock struct{}

func (r MovementRepositoryMock) CreateNewMovement(movement models.Movement) (models.Movement, error) {
	wrapper := error_const.MovementRepositoryError{}
	switch movement.Name {
	case "Success":
		return movement, nil
	case "CantConnect":
		wrapper.Wrap(error_const.FirebaseUnknownError{})
		return models.Movement{}, wrapper
	case "Duplicated":
		wrapper.Wrap(error_const.FirestoreAlreadyExistsError{DocID: movement.UID})
		return models.Movement{}, wrapper
	default:
		return models.Movement{}, error_const.TestInvalidTestCaseError{Param: movement.Name}
	}
}

var TestMovement = models.Movement{UID: "0", Name: "TEST_SOURCE"}

func (r MovementRepositoryMock) GetMovementByID(id string, user string) (models.Movement, error) {
	wrapper := error_const.MovementRepositoryError{}
	switch id {
	case "0":
		return TestMovement, nil
	case "1":
		wrapper.Wrap(error_const.FirebaseUnknownError{})
		return models.Movement{}, wrapper
	case "2":
		wrapper.Wrap(error_const.FirestoreNotFoundError{DocID: id})
		return models.Movement{}, wrapper
	case "3":
		wrapper.Wrap(error_const.FirestoreParsingError{DocID: id, StructName: "movement"})
		return models.Movement{}, wrapper
	default:
		return models.Movement{}, error_const.TestInvalidTestCaseError{Param: id}
	}
}

func (r MovementRepositoryMock) GetMovementsByUserAndDate(userID string, page int, date_from time.Time, date_to time.Time) (util_models.PaginatedSearch[models.Movement], error) {
	wrapper := error_const.MovementRepositoryError{}
	switch userID {
	case "0":
		return util_models.PaginatedSearch[models.Movement]{
			CurrentPage: page,
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
		wrapper.Wrap(error_const.MovementNotEnoughData{})
		return util_models.PaginatedSearch[models.Movement]{}, wrapper
	default:
		return util_models.PaginatedSearch[models.Movement]{}, error_const.TestInvalidTestCaseError{Param: userID}
	}
}

func (r MovementRepositoryMock) UpdateMovement(movement models.Movement) (models.Movement, error) {
	id := movement.UID
	wrapper := error_const.MovementRepositoryError{}
	switch id {
	case "0":
		return movement, nil
	case "1":
		wrapper.Wrap(error_const.FirebaseUnknownError{})
		return models.Movement{}, wrapper
	case "2":
		wrapper.Wrap(error_const.FirestoreNotFoundError{DocID: id})
		return models.Movement{}, wrapper
	case "3":
		wrapper.Wrap(error_const.FirestoreParsingError{DocID: id, StructName: "movement"})
		return models.Movement{}, wrapper
	default:
		return models.Movement{}, error_const.TestInvalidTestCaseError{Param: id}
	}
}

func (r MovementRepositoryMock) DeleteMovement(id string, userID string) error {
	wrapper := error_const.MovementRepositoryError{}
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
		wrapper.Wrap(error_const.FirestoreParsingError{DocID: id, StructName: "movement"})
		return wrapper
	default:
		return error_const.TestInvalidTestCaseError{Param: id}
	}
}
