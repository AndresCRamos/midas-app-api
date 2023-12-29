package mocks

import (
	"github.com/AndresCRamos/midas-app-api/models"
	error_const "github.com/AndresCRamos/midas-app-api/utils/errors"
)

type UserServiceMock struct{}

func (r UserServiceMock) CreateNewUser(user models.User) error {
	switch user.Name {
	case "Success":
		return nil
	case "CantConnect":
		return error_const.FirebaseUnknownError{}
	case "Duplicated":
		return error_const.FirestoreAlreadyExistsError{DocID: user.UID}
	default:
		return error_const.TestInvalidTestCaseError{Param: user.Name}
	}
}

func (r UserServiceMock) GetUserByID(id string) (models.User, error) {
	switch id {
	case "0":
		return TestUser, nil
	case "1":
		return models.User{}, error_const.FirebaseUnknownError{}
	case "2":
		return models.User{}, error_const.FirestoreNotFoundError{DocID: id}
	case "3":
		return models.User{}, error_const.FirestoreParsingError{DocID: id, StructName: "user"}
	default:
		return models.User{}, error_const.TestInvalidTestCaseError{Param: id}
	}
}
