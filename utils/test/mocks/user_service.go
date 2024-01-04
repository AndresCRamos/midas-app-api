package mocks

import (
	"github.com/AndresCRamos/midas-app-api/models"
	error_const "github.com/AndresCRamos/midas-app-api/utils/errors"
)

type UserServiceMock struct{}

func (r UserServiceMock) CreateNewUser(user models.User) error {
	ServiceWrapper := error_const.UserServiceError{}
	RepoWrapper := error_const.UserServiceError{}
	switch user.Name {
	case "Success":
		return nil
	case "CantConnect":
		RepoWrapper.Wrap(error_const.FirebaseUnknownError{})
		ServiceWrapper.Wrap(RepoWrapper)
		return ServiceWrapper
	case "Duplicated":
		RepoWrapper.Wrap(error_const.FirestoreAlreadyExistsError{DocID: user.UID})
		ServiceWrapper.Wrap(RepoWrapper)
		return ServiceWrapper
	default:
		return error_const.TestInvalidTestCaseError{Param: user.Name}
	}
}

func (r UserServiceMock) GetUserByID(id string) (models.User, error) {
	ServiceWrapper := error_const.UserServiceError{}
	RepoWrapper := error_const.UserServiceError{}
	switch id {
	case "0":
		return TestUser, nil
	case "1":
		RepoWrapper.Wrap(error_const.FirebaseUnknownError{})
		ServiceWrapper.Wrap(RepoWrapper)
		return models.User{}, ServiceWrapper
	case "2":
		RepoWrapper.Wrap(error_const.FirestoreNotFoundError{DocID: id})
		ServiceWrapper.Wrap(RepoWrapper)
		return models.User{}, ServiceWrapper
	case "3":
		RepoWrapper.Wrap(error_const.FirestoreParsingError{DocID: id, StructName: "user"})
		ServiceWrapper.Wrap(RepoWrapper)
		return models.User{}, ServiceWrapper
	default:
		return models.User{}, error_const.TestInvalidTestCaseError{Param: id}
	}
}
