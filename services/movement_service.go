package services

import (
	"github.com/AndresCRamos/midas-app-api/models"
	"github.com/AndresCRamos/midas-app-api/repository"
)

type MovementService interface {
	GetMovementByID(id string) (models.Movement, error)
}

type movementServiceImplementation struct {
	r repository.MovementRepository
}

func NewMovementService(r repository.MovementRepository) *movementServiceImplementation {
	return &movementServiceImplementation{
		r: r,
	}
}

func (s *movementServiceImplementation) GetMovementByID(id string) (models.Movement, error) {
	return s.r.GetMovementByID(id)
}
