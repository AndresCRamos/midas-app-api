package repository

import (
	"context"

	"cloud.google.com/go/firestore"
	"github.com/AndresCRamos/midas-app-api/models"
	error_utils "github.com/AndresCRamos/midas-app-api/utils/errors"
	firebase_utils "github.com/AndresCRamos/midas-app-api/utils/firebase"
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

func (r *movementRepositoryImplementation) CreateNewMovement(movement models.Movement) (models.Movement, error) {
	userExists, err := firebase_utils.CheckDocumentExists(r.client.Collection("users"), movement.OwnerId)
	if err != nil {
		return models.Movement{}, error_utils.CheckFirebaseError(err, movement.SourceID, &error_utils.MovementRepositoryError{})
	}
	if !userExists {
		repoErr := error_utils.MovementRepositoryError{}
		repoErr.Wrap(error_utils.MovementOwnerNotFound{MovementID: movement.UID, OwnerId: movement.OwnerId})
		return models.Movement{}, repoErr
	}

	sourceExists, err := firebase_utils.CheckDocumentExists(r.client.Collection("sources"), movement.SourceID)
	if err != nil {
		return models.Movement{}, error_utils.CheckFirebaseError(err, movement.SourceID, &error_utils.MovementRepositoryError{})
	}
	if !sourceExists {
		repoErr := error_utils.MovementRepositoryError{}
		repoErr.Wrap(error_utils.MovementSourceNotFound{MovementID: movement.UID, SourceID: movement.SourceID})
		return models.Movement{}, repoErr
	}

	movementCollection := r.client.Collection("movements")

	movement.NewCreationAtDate()
	movement.NewUpdatedAtDate()

	docRef := movementCollection.NewDoc()
	movement.UID = docRef.ID

	_, err = docRef.Set(context.Background(), movement)
	if err != nil {
		wrapErr := error_utils.MovementRepositoryError{}
		return models.Movement{}, error_utils.CheckFirebaseError(err, movement.UID, &wrapErr)
	}
	return movement, nil
}

func (r *movementRepositoryImplementation) GetMovementByID(id string) (models.Movement, error) {
	return models.Movement{}, nil
}
