package services

import (
	"fmt"

	"github.com/AndresCRamos/midas-app-api/models"
	"github.com/AndresCRamos/midas-app-api/repository"
)

type UserService interface {
	CreateNewUser(user models.User) error
	GetUserByID(id string) (models.User, error)
}

type userServiceImplementation struct {
	r repository.UserRepository
}

func NewService(r repository.UserRepository) *userServiceImplementation {
	return &userServiceImplementation{
		r: r,
	}
}

func (s *userServiceImplementation) CreateNewUser(user models.User) error {
	err := s.r.CreateNewUser(user)
	if err != nil {
		return fmt.Errorf("UserService: Cant create User: %w", err)
	}
	return nil
}

func (s *userServiceImplementation) GetUserByID(id string) (models.User, error) {
	res, err := s.r.GetUserByID(id)
	if err != nil {
		return models.User{}, fmt.Errorf("UserService: Cant get User: %w", err)
	}
	return res, nil
}
