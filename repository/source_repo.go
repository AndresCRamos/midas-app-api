package repository

import (
	"context"
	"errors"

	"cloud.google.com/go/firestore"
	"cloud.google.com/go/firestore/apiv1/firestorepb"
	"github.com/AndresCRamos/midas-app-api/models"
	util_models "github.com/AndresCRamos/midas-app-api/utils/api/models"
	error_utils "github.com/AndresCRamos/midas-app-api/utils/errors"
	"google.golang.org/api/iterator"
)

type SourceRepository interface {
	GetSourcesByUser(userID string, page int) (util_models.PaginatedSearch[models.Source], error)
	GetSourceByID(id string, userID string) (models.Source, error)
	CreateNewSource(Source models.Source) (models.Source, error)
	UpdateSource(Source models.Source) (models.Source, error)
	DeleteSource(id string, userID string) error
}

const (
	pageSize = 50
)

type SourceRepositoryImplementation struct {
	client *firestore.Client
}

func NewSourceRepository(client *firestore.Client) *SourceRepositoryImplementation {
	return &SourceRepositoryImplementation{
		client: client,
	}
}

func getTotalSizeOfQuery(query firestore.Query) (int, error) {
	aggregationQuery := query.NewAggregationQuery().WithCount("all")
	results, err := aggregationQuery.Get(context.Background())
	if err != nil {
		return 0, err
	}

	count, ok := results["all"]
	if !ok {
		return 0, errors.New("firestore: couldn't get alias for COUNT from results")
	}

	countValue := count.(*firestorepb.Value)
	return int(countValue.GetIntegerValue()), nil
}

func (r *SourceRepositoryImplementation) GetSourcesByUser(userID string, page int) (util_models.PaginatedSearch[models.Source], error) {
	var sources []models.Source
	sourceCollection := r.client.Collection("sources")
	totalQuery := sourceCollection.Where("owner", "==", userID).OrderBy("created_at", firestore.Desc)
	iterSource := totalQuery.Offset((page - 1) * pageSize).Limit(pageSize).Documents(context.Background())

	totalSize, _ := getTotalSizeOfQuery(totalQuery)

	minPageData := ((page - 1) * pageSize) + 1

	if totalSize < minPageData {
		wrapErr := error_utils.SourceRepositoryError{
			Err: error_utils.SourceNotEnoughData{},
		}
		return util_models.PaginatedSearch[models.Source]{}, wrapErr
	}

	for {
		sourceDoc, err := iterSource.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			wrapErr := error_utils.SourceRepositoryError{}
			return util_models.PaginatedSearch[models.Source]{}, error_utils.CheckFirebaseError(err, "", &wrapErr)
		}
		var sourceModel models.Source
		if err := sourceDoc.DataTo(&sourceModel); err != nil {
			wrapErr := error_utils.SourceRepositoryError{}
			return util_models.PaginatedSearch[models.Source]{}, error_utils.CheckFirebaseError(err, "", &wrapErr)
		}
		sources = append(sources, sourceModel)
	}

	return util_models.PaginatedSearch[models.Source]{
		CurrentPage: page,
		TotalData:   totalSize,
		PageSize:    len(sources),
		Data:        sources,
	}, nil
}

func (r *SourceRepositoryImplementation) GetSourceByID(id string, userID string) (models.Source, error) {

	SourceDoc, err := getSourceDocSnapByID(id, r.client)

	if err != nil {
		wrapEr := error_utils.SourceRepositoryError{}
		return models.Source{}, error_utils.CheckFirebaseError(err, id, &wrapEr)
	}

	var Source models.Source

	if err = SourceDoc.DataTo(&Source); err != nil {
		wrapErr := error_utils.SourceRepositoryError{}
		logged_err := error_utils.FirestoreParsingError{DocID: Source.UID, StructName: "Source"}
		wrapErr.Wrap(logged_err)
		return models.Source{}, wrapErr
	}

	if userID != Source.OwnerId {
		return models.Source{}, error_utils.SourceRepositoryError{
			Err: error_utils.SourceDifferentOwner{
				SourceID: id,
				OwnerID:  userID,
			},
		}
	}

	return Source, nil
}

func (r *SourceRepositoryImplementation) CreateNewSource(Source models.Source) (models.Source, error) {

	SourceCollection := r.client.Collection("sources")

	userDocRef, _ := r.client.Collection("users").Doc(Source.OwnerId).Get(context.Background())

	Source.NewCreationAtDate()
	Source.NewUpdatedAtDate()

	if userDocRef != nil && !userDocRef.Exists() {
		wrapErr := error_utils.SourceRepositoryError{}
		wrapErr.Wrap(error_utils.SourceOwnerNotFound{SourceID: Source.UID, OwnerId: Source.OwnerId})
		return models.Source{}, wrapErr
	}

	docRef := SourceCollection.NewDoc()
	Source.UID = docRef.ID

	_, err := docRef.Set(context.Background(), Source)
	if err != nil {
		wrapErr := error_utils.SourceRepositoryError{}
		return models.Source{}, error_utils.CheckFirebaseError(err, Source.UID, &wrapErr)
	}
	return Source, nil
}

func (r *SourceRepositoryImplementation) UpdateSource(source models.Source) (models.Source, error) {
	sourceDoc, err := getSourceDocSnapByID(source.UID, r.client)

	source.NewUpdatedAtDate()

	if err != nil {
		wrapErr := error_utils.SourceRepositoryError{}
		return models.Source{}, error_utils.CheckFirebaseError(err, source.UID, &wrapErr)
	}

	var prevData models.Source

	if err := sourceDoc.DataTo(&prevData); err != nil {
		wrapErr := error_utils.SourceRepositoryError{}
		return models.Source{}, error_utils.CheckFirebaseError(err, source.UID, &wrapErr)
	}

	if prevData.OwnerId != source.OwnerId {
		wrapErr := error_utils.SourceRepositoryError{}
		wrapErr.Wrap(error_utils.SourceDifferentOwner{SourceID: source.UID, OwnerID: source.OwnerId})
		return models.Source{}, wrapErr
	}

	sourceMap, isChanged := sourceStructToMap(source)

	source.CreatedAt = prevData.CreatedAt

	if !isChanged {
		source.UpdatedAt = prevData.UpdatedAt
	}

	_, err = sourceDoc.Ref.Set(context.Background(), sourceMap, firestore.MergeAll)
	if err != nil {
		wrapErr := error_utils.SourceRepositoryError{}
		return models.Source{}, error_utils.CheckFirebaseError(err, source.UID, &wrapErr)
	}
	return source, nil
}

func (r *SourceRepositoryImplementation) DeleteSource(id string, userID string) error {
	sourceDoc, err := getSourceDocSnapByID(id, r.client)

	if err != nil {
		wrapErr := error_utils.SourceRepositoryError{}
		return error_utils.CheckFirebaseError(err, id, &wrapErr)
	}

	var sourceData models.Source

	if err := sourceDoc.DataTo(&sourceData); err != nil {
		return error_utils.SourceRepositoryError{
			Err: error_utils.FirestoreParsingError{DocID: id, StructName: "Source"},
		}
	}

	if sourceData.OwnerId != userID {
		return error_utils.SourceRepositoryError{
			Err: error_utils.SourceDifferentOwner{SourceID: id, OwnerID: userID},
		}
	}

	_, err = sourceDoc.Ref.Delete(context.Background())
	if err != nil {
		wrapErr := error_utils.SourceRepositoryError{}
		return error_utils.CheckFirebaseError(err, id, &wrapErr)
	}
	return nil
}

func getSourceDocSnapByID(id string, client *firestore.Client) (*firestore.DocumentSnapshot, error) {
	SourceCollection := client.Collection("sources")

	SourceDoc, err := SourceCollection.Doc(id).Get(context.Background())

	if err != nil {
		return nil, err
	}
	return SourceDoc, nil
}

func sourceStructToMap(source models.Source) (map[string]any, bool) {
	change := false
	fields := make(map[string]any)
	if source.Name != "" {
		fields["name"] = source.Name
		change = true
	}
	if source.Description != "" {
		fields["description"] = source.Description
		change = true
	}
	if change {
		fields["updated_at"] = source.UpdatedAt
	}

	return fields, change
}
