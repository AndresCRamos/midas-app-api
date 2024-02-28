package mocks

import (
	"time"

	"github.com/AndresCRamos/midas-app-api/models"
	util_models "github.com/AndresCRamos/midas-app-api/utils/api/models"
	error_const "github.com/AndresCRamos/midas-app-api/utils/errors"
)

type SourceServiceMock struct{}

func (r SourceServiceMock) CreateNewSource(user models.Source) (models.Source, error) {
	ServiceWrapper := error_const.SourceRepositoryError{}
	RepoWrapper := error_const.SourceServiceError{}
	switch user.Name {
	case "Success":
		return user, nil
	case "CantConnect":
		RepoWrapper.Wrap(error_const.FirebaseUnknownError{})
		ServiceWrapper.Wrap(RepoWrapper)
		return models.Source{}, ServiceWrapper
	case "NoOwner":
		RepoWrapper.Wrap(error_const.SourceOwnerNotFound{SourceID: user.UID, OwnerId: user.OwnerId})
		ServiceWrapper.Wrap(RepoWrapper)
		return models.Source{}, ServiceWrapper
	case "Duplicated":
		RepoWrapper.Wrap(error_const.FirestoreAlreadyExistsError{DocID: user.UID})
		ServiceWrapper.Wrap(RepoWrapper)
		return models.Source{}, ServiceWrapper
	default:
		return models.Source{}, error_const.TestInvalidTestCaseError{Param: user.Name}
	}
}

func (r SourceServiceMock) GetSourcesByUser(userID string, page int) (util_models.PaginatedSearch[models.Source], error) {
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

func (r SourceServiceMock) GetMovementsBySourceAndDate(id string, userID string, page int, date_from time.Time, date_to time.Time) (util_models.PaginatedSearch[models.Movement], error) {
	RepoWrapper := &error_const.SourceRepositoryError{}
	wrapper := &error_const.SourceServiceError{}
	wrapper.Wrap(RepoWrapper)
	switch id {
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
		RepoWrapper.Wrap(error_const.FirebaseUnknownError{})
		return util_models.PaginatedSearch[models.Movement]{}, wrapper
	case "2":
		RepoWrapper.Wrap(error_const.SourceNotFound{SourceID: id})
		return util_models.PaginatedSearch[models.Movement]{}, wrapper
	case "3":
		RepoWrapper.Wrap(error_const.SourceDifferentOwner{SourceID: id, OwnerID: userID})
		return util_models.PaginatedSearch[models.Movement]{}, wrapper
	case "4":
		RepoWrapper.Wrap(error_const.MovementNotEnoughData{})
		return util_models.PaginatedSearch[models.Movement]{}, wrapper
	case "5":
		RepoWrapper.Wrap(error_const.SourceBadDates{})
		return util_models.PaginatedSearch[models.Movement]{}, wrapper
	default:
		return util_models.PaginatedSearch[models.Movement]{}, error_const.TestInvalidTestCaseError{Param: userID}
	}

}

func (r SourceServiceMock) GetSourceByID(id string, userId string) (models.Source, error) {
	ServiceWrapper := error_const.SourceServiceError{}
	RepoWrapper := error_const.SourceRepositoryError{}
	switch id {
	case "0":
		return TestSource, nil
	case "1":
		RepoWrapper.Wrap(error_const.FirebaseUnknownError{})
		ServiceWrapper.Wrap(RepoWrapper)
		return models.Source{}, ServiceWrapper
	case "2":
		RepoWrapper.Wrap(error_const.FirestoreNotFoundError{DocID: id})
		ServiceWrapper.Wrap(RepoWrapper)
		return models.Source{}, ServiceWrapper
	case "3":
		RepoWrapper.Wrap(error_const.FirestoreParsingError{DocID: id, StructName: "source"})
		ServiceWrapper.Wrap(RepoWrapper)
		return models.Source{}, ServiceWrapper
	case "4":
		RepoWrapper.Wrap(error_const.SourceDifferentOwner{SourceID: id, OwnerID: userId})
		ServiceWrapper.Wrap(RepoWrapper)
		return models.Source{}, ServiceWrapper
	default:
		return models.Source{}, error_const.TestInvalidTestCaseError{Param: id}
	}
}

func (r SourceServiceMock) UpdateSource(source models.Source) (models.Source, error) {
	id := source.UID
	ServiceWrapper := error_const.SourceRepositoryError{}
	RepoWrapper := error_const.SourceServiceError{}
	switch id {
	case "0":
		return source, nil
	case "1":
		RepoWrapper.Wrap(error_const.FirebaseUnknownError{})
		ServiceWrapper.Wrap(RepoWrapper)
		return models.Source{}, ServiceWrapper
	case "2":
		RepoWrapper.Wrap(error_const.FirestoreNotFoundError{DocID: id})
		ServiceWrapper.Wrap(RepoWrapper)
		return models.Source{}, ServiceWrapper
	case "3":
		RepoWrapper.Wrap(error_const.FirestoreParsingError{DocID: id, StructName: "source"})
		ServiceWrapper.Wrap(RepoWrapper)
		return models.Source{}, ServiceWrapper
	case "4":
		RepoWrapper.Wrap(error_const.SourceDifferentOwner{SourceID: id, OwnerID: source.OwnerId})
		ServiceWrapper.Wrap(RepoWrapper)
		return models.Source{}, ServiceWrapper
	default:
		return models.Source{}, error_const.TestInvalidTestCaseError{Param: id}
	}
}

func (r SourceServiceMock) DeleteSource(id string, userID string) error {
	ServiceWrapper := error_const.SourceRepositoryError{}
	RepoWrapper := error_const.SourceServiceError{}
	switch id {
	case "0":
		return nil
	case "1":
		RepoWrapper.Wrap(error_const.FirebaseUnknownError{})
		ServiceWrapper.Wrap(RepoWrapper)
		return ServiceWrapper
	case "2":
		RepoWrapper.Wrap(error_const.FirestoreNotFoundError{DocID: id})
		ServiceWrapper.Wrap(RepoWrapper)
		return ServiceWrapper
	case "3":
		RepoWrapper.Wrap(error_const.FirestoreParsingError{DocID: id, StructName: "source"})
		ServiceWrapper.Wrap(RepoWrapper)
		return ServiceWrapper
	case "4":
		RepoWrapper.Wrap(error_const.SourceDifferentOwner{SourceID: id, OwnerID: userID})
		ServiceWrapper.Wrap(RepoWrapper)
		return ServiceWrapper
	default:
		return error_const.TestInvalidTestCaseError{Param: id}
	}
}
