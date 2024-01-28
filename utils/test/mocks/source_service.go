package mocks

import (
	"github.com/AndresCRamos/midas-app-api/models"
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

func (r SourceServiceMock) GetSourceByID(id string) (models.Source, error) {
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
	default:
		return models.Source{}, error_const.TestInvalidTestCaseError{Param: id}
	}
}

func (r SourceServiceMock) DeleteSource(id string) error {
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
	default:
		return error_const.TestInvalidTestCaseError{Param: id}
	}
}
