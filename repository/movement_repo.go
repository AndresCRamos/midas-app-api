package repository

import (
	"cloud.google.com/go/firestore"
	"github.com/AndresCRamos/midas-app-api/models"
)

type MovementRepository interface {
	GetMovementByID(id string) (models.Movement, error)
}

type movementRepositoryImplementation struct {
	client *firestore.Client
}

func NewMovementRepository(client *firestore.Client) *movementRepositoryImplementation {
	return &movementRepositoryImplementation{
		client: client,
	}
}

func (r *movementRepositoryImplementation) GetMovementByID(id string) (models.Movement, error) {
	return models.Movement{}, nil
}
