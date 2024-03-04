package mocks

import (
	"time"

	"github.com/AndresCRamos/midas-app-api/models"
	util_models "github.com/AndresCRamos/midas-app-api/utils/api/models"
	error_const "github.com/AndresCRamos/midas-app-api/utils/errors"
)

var (
	TestMovementRetrieve = models.MovementRetrieve{UID: "0", Name: "TEST_SOURCE"}
)

type MovementServiceMock struct{}

func (r MovementServiceMock) CreateNewMovement(movement models.Movement) (models.Movement, error) {
	wrapper := &error_const.MovementServiceError{}
	repoWrapper := &error_const.MovementRepositoryError{}
	wrapper.Wrap(repoWrapper)
	switch movement.Name {
	case "Success":
		return movement, nil
	case "CantConnect":
		repoWrapper.Wrap(error_const.FirebaseUnknownError{})
		return models.Movement{}, wrapper
	case "NoSource":
		repoWrapper.Wrap(error_const.MovementSourceNotFound{MovementID: movement.UID, SourceID: movement.SourceID})
		return models.Movement{}, wrapper
	case "Duplicated":
		repoWrapper.Wrap(error_const.FirestoreAlreadyExistsError{DocID: movement.UID})
		return models.Movement{}, wrapper
	default:
		return models.Movement{}, error_const.TestInvalidTestCaseError{Param: movement.Name}
	}
}

func (r MovementServiceMock) GetMovementByID(id string, user string) (models.Movement, error) {
	wrapper := &error_const.MovementRepositoryError{}
	repoWrapper := &error_const.MovementRepositoryError{}
	wrapper.Wrap(repoWrapper)
	switch id {
	case "0":
		return TestMovement, nil
	case "1":
		repoWrapper.Wrap(error_const.FirebaseUnknownError{})
		return models.Movement{}, wrapper
	case "2":
		repoWrapper.Wrap(error_const.FirestoreNotFoundError{DocID: id})
		return models.Movement{}, wrapper
	case "3":
		repoWrapper.Wrap(error_const.FirestoreParsingError{DocID: id, StructName: "movement"})
		return models.Movement{}, wrapper
	case "4":
		repoWrapper.Wrap(error_const.MovementDifferentOwner{MovementID: id, OwnerID: user})
		return models.Movement{}, wrapper
	default:
		return models.Movement{}, error_const.TestInvalidTestCaseError{Param: id}
	}
}

func (r MovementServiceMock) GetMovementsByUserAndDate(userID string, page int, date_from time.Time, date_to time.Time) (util_models.PaginatedSearch[models.Movement], error) {
	wrapper := &error_const.MovementRepositoryError{}
	repoWrapper := &error_const.MovementRepositoryError{}
	wrapper.Wrap(repoWrapper)
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
		repoWrapper.Wrap(error_const.FirebaseUnknownError{})
		return util_models.PaginatedSearch[models.Movement]{}, wrapper
	case "2":
		repoWrapper.Wrap(error_const.MovementNotEnoughData{})
		return util_models.PaginatedSearch[models.Movement]{}, wrapper
	default:
		return util_models.PaginatedSearch[models.Movement]{}, error_const.TestInvalidTestCaseError{Param: userID}
	}
}

func (r MovementServiceMock) UpdateMovement(movement models.Movement) (models.Movement, error) {
	id := movement.UID
	wrapper := &error_const.MovementRepositoryError{}
	repoWrapper := &error_const.MovementRepositoryError{}
	wrapper.Wrap(repoWrapper)
	switch id {
	case "0":
		return movement, nil
	case "1":
		repoWrapper.Wrap(error_const.FirebaseUnknownError{})
		return models.Movement{}, wrapper
	case "2":
		repoWrapper.Wrap(error_const.FirestoreNotFoundError{DocID: id})
		return models.Movement{}, wrapper
	case "3":
		repoWrapper.Wrap(error_const.FirestoreParsingError{DocID: id, StructName: "movement"})
		return models.Movement{}, wrapper
	default:
		return models.Movement{}, error_const.TestInvalidTestCaseError{Param: id}
	}
}

func (r MovementServiceMock) DeleteMovement(id string, userID string) error {
	wrapper := &error_const.MovementRepositoryError{}
	repoWrapper := &error_const.MovementRepositoryError{}
	wrapper.Wrap(repoWrapper)
	switch id {
	case "0":
		return nil
	case "1":
		repoWrapper.Wrap(error_const.FirebaseUnknownError{})
		return wrapper
	case "2":
		repoWrapper.Wrap(error_const.FirestoreNotFoundError{DocID: id})
		return wrapper
	case "3":
		repoWrapper.Wrap(error_const.FirestoreParsingError{DocID: id, StructName: "movement"})
		return wrapper
	default:
		return error_const.TestInvalidTestCaseError{Param: id}
	}
}
