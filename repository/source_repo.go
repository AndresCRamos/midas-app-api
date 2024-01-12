package repository

import (
	"context"

	"cloud.google.com/go/firestore"
	"github.com/AndresCRamos/midas-app-api/models"
	error_utils "github.com/AndresCRamos/midas-app-api/utils/errors"
)

type SourceRepository interface {
	GetSourceByID(id string) (models.Source, error)
	CreateNewSource(Source models.Source) error
	UpdateSource(Source models.Source) error
	DeleteSource(id string) error
}

type SourceRepositoryImplementation struct {
	client *firestore.Client
}

func NewSourceRepository(client *firestore.Client) *SourceRepositoryImplementation {
	return &SourceRepositoryImplementation{
		client: client,
	}
}

func (r *SourceRepositoryImplementation) GetSourceByID(id string) (models.Source, error) {

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

	return Source, nil
}

func (r *SourceRepositoryImplementation) CreateNewSource(Source models.Source) error {

	SourceCollection := r.client.Collection("sources")

	userDocRef, _ := r.client.Collection("users").Doc(Source.OwnerId).Get(context.Background())

	Source.NewCreationAtDate()
	Source.NewUpdatedAtDate()

	if userDocRef == nil {
		wrapErr := error_utils.SourceRepositoryError{}
		wrapErr.Wrap(error_utils.SourceOwnerNotFound{SourceID: Source.UID, OwnerId: Source.OwnerId})
		return wrapErr
	}

	_, err := SourceCollection.Doc(Source.UID).Create(context.Background(), Source)
	if err != nil {
		wrapErr := error_utils.SourceRepositoryError{}
		return error_utils.CheckFirebaseError(err, Source.UID, &wrapErr)
	}
	return nil
}

func (r *SourceRepositoryImplementation) UpdateSource(source models.Source) error {
	source.NewUpdatedAtDate()

	sourceDoc, err := getSourceDocSnapByID(source.UID, r.client)

	if err != nil {
		wrapErr := error_utils.SourceRepositoryError{}
		return error_utils.CheckFirebaseError(err, source.UID, &wrapErr)
	}

	_, err = sourceDoc.Ref.Set(context.Background(), source)
	if err != nil {
		wrapErr := error_utils.SourceRepositoryError{}
		return error_utils.CheckFirebaseError(err, source.UID, &wrapErr)
	}
	return nil
}

func (r *SourceRepositoryImplementation) DeleteSource(id string) error {
	sourceDoc, err := getSourceDocSnapByID(id, r.client)

	if err != nil {
		wrapErr := error_utils.SourceRepositoryError{}
		return error_utils.CheckFirebaseError(err, id, &wrapErr)
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
