package mocks

import (
	"fmt"

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
		return fmt.Errorf(error_const.ALREADY_EXISTS, user.UID)
	default:
		return fmt.Errorf(error_const.INVALID_TEST_CASE, user)
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
		return models.User{}, fmt.Errorf(error_const.FIRESTORE_NOT_FOUND, id)
	case "3":
		return models.User{}, fmt.Errorf(error_const.PARSING_ERROR, id, "user")
	default:
		return models.User{}, fmt.Errorf(error_const.INVALID_TEST_CASE, id)
	}
}
