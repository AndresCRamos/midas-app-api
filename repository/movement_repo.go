package repository

import (
	"context"
	"time"

	"cloud.google.com/go/firestore"
	"github.com/AndresCRamos/midas-app-api/models"
	util_models "github.com/AndresCRamos/midas-app-api/utils/api/models"
	error_utils "github.com/AndresCRamos/midas-app-api/utils/errors"
	firebase_utils "github.com/AndresCRamos/midas-app-api/utils/firebase"
	"google.golang.org/api/iterator"
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

	if err = movementDocSnap.DataTo(&movement); err != nil {
		wrapErr := error_utils.SourceRepositoryError{}
		logged_err := error_utils.FirestoreParsingError{DocID: id, StructName: "movement"}
		wrapErr.Wrap(logged_err)
		return models.Movement{}, wrapErr
	}

	if userID != movement.OwnerId {
		wrapErr := error_utils.SourceRepositoryError{}
		logged_err := error_utils.MovementDifferentOwner{MovementID: id, OwnerID: userID}
		wrapErr.Wrap(logged_err)
		return models.Movement{}, wrapErr
	}

	return movement, nil
}
func (r *movementRepositoryImplementation) GetMovementsByUserAndDate(userID string, page int, from_date time.Time, to_date time.Time) (util_models.PaginatedSearch[models.Movement], error) {
	movementCollection := r.client.Collection("movement")
	totalQuery := movementCollection.Where("owner", "==", userID).Where("movement_date", ">=", from_date).Where("movement_date", "<=", to_date).OrderBy("movement_date", firestore.Desc)
	iterSource := totalQuery.Offset((page - 1) * pageSize).Limit(pageSize).Documents(context.Background())

	totalSize, _ := getTotalSizeOfQuery(totalQuery)

	minPageData := ((page - 1) * pageSize) + 1

	if totalSize < minPageData {
		wrapErr := error_utils.MovementRepositoryError{
			Err: error_utils.MovementNotEnoughData{},
		}
		return util_models.PaginatedSearch[models.Movement]{}, wrapErr
	}

	var movements []models.Movement

	for {
		sourceDoc, err := iterSource.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			wrapErr := error_utils.SourceRepositoryError{}
			return util_models.PaginatedSearch[models.Movement]{}, error_utils.CheckFirebaseError(err, "", &wrapErr)
		}
		var movementModel models.Movement
		if err := sourceDoc.DataTo(&movementModel); err != nil {
			wrapErr := error_utils.SourceRepositoryError{}
			return util_models.PaginatedSearch[models.Movement]{}, error_utils.CheckFirebaseError(err, "", &wrapErr)
		}
		movements = append(movements, movementModel)
	}

	return util_models.PaginatedSearch[models.Movement]{
		CurrentPage: page,
		TotalData:   totalSize,
		PageSize:    len(movements),
		Data:        movements,
	}, nil
}

func getMovementDocSnapByID(id string, client *firestore.Client) (*firestore.DocumentSnapshot, error) {

	movementDocSnap, err := client.Collection("movements").Doc(id).Get(context.Background())

	if err != nil {
		return nil, err
	}
	return movementDocSnap, nil
}
