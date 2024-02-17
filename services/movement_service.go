package services

import (
	"time"

	"github.com/AndresCRamos/midas-app-api/models"
	"github.com/AndresCRamos/midas-app-api/repository"
	util_models "github.com/AndresCRamos/midas-app-api/utils/api/models"
	error_utils "github.com/AndresCRamos/midas-app-api/utils/errors"
)

type MovementService interface {
	CreateNewMovement(Movement models.Movement) (models.Movement, error)
	GetMovementByID(id string, userID string) (models.Movement, error)
	GetMovementsByUserAndDate(userID string, page int, date_from time.Time, date_to time.Time) (util_models.PaginatedSearch[models.Movement], error)
	UpdateMovement(Movement models.Movement) (models.Movement, error)
	DeleteMovement(id string, userID string) error
}

type movementServiceImplementation struct {
	r repository.MovementRepository
}

func NewMovementService(r repository.MovementRepository) *movementServiceImplementation {
	return &movementServiceImplementation{
		r: r,
	}
}

func (s *movementServiceImplementation) CreateNewMovement(movement models.Movement) (models.Movement, error) {
	movement, err := s.r.CreateNewMovement(movement)
	if err != nil {
		movementServiceErr := error_utils.MovementServiceError{Err: err, Method: "Create"}
		return models.Movement{}, movementServiceErr
	}
	return movement, nil
}

func (s *movementServiceImplementation) GetMovementsByUserAndDate(userID string, page int, date_from time.Time, date_to time.Time) (util_models.PaginatedSearch[models.Movement], error) {
	movement, err := s.r.GetMovementsByUserAndDate(userID, page, date_from, date_to)
	if err != nil {
		movementServiceErr := error_utils.MovementServiceError{Err: err, Method: "List"}
		return util_models.PaginatedSearch[models.Movement]{}, movementServiceErr
	}
	return movement, nil
}

func (s *movementServiceImplementation) GetMovementByID(id string, userID string) (models.Movement, error) {
	res, err := s.r.GetMovementByID(id, userID)
	if err != nil {
		movementServiceErr := error_utils.MovementServiceError{Err: err, Method: "Retrieve"}
		return models.Movement{}, movementServiceErr
	}
	return res, nil
}

func (s *movementServiceImplementation) UpdateMovement(movement models.Movement) (models.Movement, error) {
	movement, err := s.r.UpdateMovement(movement)
	if err != nil {
		movementServiceErr := error_utils.MovementServiceError{Err: err, Method: "Update"}
		return models.Movement{}, movementServiceErr
	}
	return movement, nil
}
func (s *movementServiceImplementation) DeleteMovement(id string, userID string) error {
	err := s.r.DeleteMovement(id, userID)
	if err != nil {
		movementServiceErr := error_utils.MovementServiceError{Err: err, Method: "Delete"}
		return movementServiceErr
	}
	return nil
}
