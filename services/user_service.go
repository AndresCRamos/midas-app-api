package services

import (
	"github.com/AndresCRamos/midas-app-api/models"
	"github.com/AndresCRamos/midas-app-api/repository"
	error_utils "github.com/AndresCRamos/midas-app-api/utils/errors"
)

type UserService interface {
	CreateNewUser(user models.User) error
	GetUserByID(id string) (models.User, error)
}

type userServiceImplementation struct {
	r repository.UserRepository
}

func NewUserService(r repository.UserRepository) *userServiceImplementation {
	return &userServiceImplementation{
		r: r,
	}
}

func (s *userServiceImplementation) CreateNewUser(user models.User) error {
	err := s.r.CreateNewUser(user)
	if err != nil {
		userServiceErr := error_utils.UserServiceError{Err: err, Method: "Create"}
		return userServiceErr
	}
	return nil
}

func (s *userServiceImplementation) GetUserByID(id string) (models.User, error) {
	res, err := s.r.GetUserByID(id)
	if err != nil {
		userServiceErr := error_utils.UserServiceError{Err: err, Method: "Retrieve"}
		return models.User{}, userServiceErr
	}
	return res, nil
}
