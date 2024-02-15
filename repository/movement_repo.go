package repository

import (
	"context"

	"cloud.google.com/go/firestore"
	"github.com/AndresCRamos/midas-app-api/models"
	error_utils "github.com/AndresCRamos/midas-app-api/utils/errors"
	firebase_utils "github.com/AndresCRamos/midas-app-api/utils/firebase"
)

type MovementRepository interface {
	CreateNewMovement(movement models.Movement) (models.Movement, error)
	GetMovementByID(id string, userID string) (models.Movement, error)
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

func (r *movementRepositoryImplementation) GetMovementByID(id string, userID string) (models.Movement, error) {
	movementDocSnap, err := getMovementDocSnapByID(id, r.client)
	if err != nil {
		wrapErr := error_utils.MovementRepositoryError{}
		return models.Movement{}, error_utils.CheckFirebaseError(err, id, &wrapErr)
	}

	var movement models.Movement
	if err = movementDocSnap.DataTo(movement); err != nil {
		wrapErr := error_utils.SourceRepositoryError{}
		logged_err := error_utils.FirestoreParsingError{DocID: movement.UID, StructName: "movement"}
		wrapErr.Wrap(logged_err)
		return models.Movement{}, wrapErr
	}

	if userID != movement.OwnerId {
		return models.Movement{}, error_utils.SourceRepositoryError{
			Err: error_utils.MovementDifferentOwner{
				MovementID: id,
				OwnerID:    userID,
			},
		}
	}

	return movement, nil
}

func getMovementDocSnapByID(id string, client *firestore.Client) (*firestore.DocumentSnapshot, error) {

	movementDocSnap, err := client.Collection("movements").Doc(id).Get(context.Background())

	if err != nil {
		return nil, err
	}
	return movementDocSnap, nil
}
