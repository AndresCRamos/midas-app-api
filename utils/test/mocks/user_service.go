package mocks

import (
	"github.com/AndresCRamos/midas-app-api/models"
	error_const "github.com/AndresCRamos/midas-app-api/utils/errors"
)

type UserServiceMock struct{}

func (r UserServiceMock) CreateNewUser(user models.User) error {
	wrapper := error_const.UserServiceError{}
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
		wrapper.Wrap(error_const.TestInvalidTestCaseError{Param: user.Name})
		return wrapper
	}
}

func (r UserServiceMock) GetUserByID(id string) (models.User, error) {
	wrapper := error_const.UserServiceError{}
	switch id {
	case "0":
		return TestUser, nil
	case "1":
		wrapper.Wrap(error_const.FirebaseUnknownError{})
		return models.User{}, wrapper
	case "2":
		wrapper.Wrap(error_const.FirestoreNotFoundError{DocID: id})
		return models.User{}, wrapper
	case "3":
		wrapper.Wrap(error_const.FirestoreParsingError{DocID: id, StructName: "user"})
		return models.User{}, wrapper
	default:
		return models.User{}, error_const.TestInvalidTestCaseError{Param: id}
	}
}
