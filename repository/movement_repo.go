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
	GetMovementsByUserAndDate(userID string, page int, from_date time.Time, to_date time.Time) (util_models.PaginatedSearch[models.Movement], error)
	UpdateMovement(movement models.Movement) (models.Movement, error)
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
	movementCollection := r.client.Collection("movements")
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

func (r *movementRepositoryImplementation) UpdateMovement(movement models.Movement) (models.Movement, error) {
	movementDocSnap, err := getMovementDocSnapByID(movement.UID, r.client)
	if err != nil {
		wrapErr := error_utils.MovementRepositoryError{}
		return models.Movement{}, error_utils.CheckFirebaseError(err, movement.UID, &wrapErr)
	}

	var prevData models.Movement

	if err = movementDocSnap.DataTo(&prevData); err != nil {
		wrapErr := error_utils.SourceRepositoryError{}
		logged_err := error_utils.FirestoreParsingError{DocID: movement.UID, StructName: "movement"}
		wrapErr.Wrap(logged_err)
		return models.Movement{}, wrapErr
	}

	if prevData.OwnerId != movement.OwnerId {
		wrapErr := error_utils.SourceRepositoryError{}
		logged_err := error_utils.MovementDifferentOwner{MovementID: movement.UID, OwnerID: movement.OwnerId}
		wrapErr.Wrap(logged_err)
		return models.Movement{}, wrapErr
	}

	movement.NewUpdatedAtDate()

	sourceMap, isChanged := movementToStruct(movement)

	if isChanged {
		_, err = movementDocSnap.Ref.Set(context.Background(), sourceMap, firestore.MergeAll)
		if err != nil {
			wrapErr := error_utils.SourceRepositoryError{}
			return models.Movement{}, error_utils.CheckFirebaseError(err, movement.UID, &wrapErr)
		}
	} else {
		movement.UpdatedAt = prevData.UpdatedAt
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

func movementToStruct(m models.Movement) (map[string]interface{}, bool) {
	var changed bool
	fields := make(map[string]interface{})

	if m.UID != "" {
		changed = true
		fields["uid"] = m.UID
	}

	if m.OwnerId != "" {
		changed = true
		fields["owner"] = m.OwnerId
	}

	if m.SourceID != "" {
		changed = true
		fields["source"] = m.SourceID
	}

	if m.Name != "" {
		changed = true
		fields["name"] = m.Name
	}

	if m.Description != "" {
		changed = true
		fields["description"] = m.Description
	}

	if m.Amount != 0 {
		changed = true
		fields["amount"] = m.Amount
	}

	if !m.MovementDate.IsZero() {
		changed = true
		fields["movement_date"] = m.MovementDate
	}

	if m.Tags != nil && len(m.Tags) != 0 {
		changed = true
		fields["tags"] = m.Tags
	}

	if changed {
		fields["updated_at"] = m.UpdatedAt
	}

	return fields, changed
}
