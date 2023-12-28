package mocks

import (
	"github.com/AndresCRamos/midas-app-api/models"
	error_const "github.com/AndresCRamos/midas-app-api/utils/errors"
)

type UserRepositoryMock struct{}

func (r UserRepositoryMock) CreateNewUser(user models.User) error {
	switch user.Name {
	case "Success":
		return nil
	case "CantConnect":
		return error_const.UNKNOWN
	case "Duplicated":
		return error_const.AlreadyExistsError{DocID: user.UID}
	default:
		return error_const.InvalidTestCaseError{Param: user.Name}
	}
}

var TestUser = models.User{UID: "0", Alias: "TEST_USER"}

func (r UserRepositoryMock) GetUserByID(id string) (models.User, error) {
	switch id {
	case "0":
		return TestUser, nil
	case "1":
		return models.User{}, error_const.UNKNOWN
	case "2":
		return models.User{}, error_const.FirestoreNotFoundError{DocID: id}
	case "3":
		return models.User{}, error_const.ParsingError{DocID: id, StructName: "user"}
	default:
		return models.User{}, error_const.InvalidTestCaseError{Param: id}
	}
}
