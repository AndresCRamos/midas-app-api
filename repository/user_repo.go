package repository

import (
	"context"

	"cloud.google.com/go/firestore"
	"github.com/AndresCRamos/midas-app-api/models"
	error_utils "github.com/AndresCRamos/midas-app-api/utils/errors"
)

type UserRepository interface {
	GetUserByID(id string) (models.User, error)
	CreateNewUser(user models.User) error
}

type UserRepositoryImplementation struct {
	client *firestore.Client
}

func NewUserRepository(client *firestore.Client) *UserRepositoryImplementation {
	return &UserRepositoryImplementation{
		client: client,
	}
}

func (r *UserRepositoryImplementation) GetUserByID(id string) (models.User, error) {
	userCollection := r.client.Collection("users")

	userDoc, err := userCollection.Doc(id).Get(context.Background())

	if err != nil {
		wrapEr := error_utils.UserRepositoryError{}
		return models.User{}, error_utils.CheckFirebaseError(err, id, models.User{}, &wrapEr)
	}

	var user models.User

	if err = userDoc.DataTo(&user); err != nil {
		wrapErr := error_utils.UserRepositoryError{}
		logged_err := error_utils.ParsingError{DocID: user.UID, StructName: "user"}
		wrapErr.Wrap(logged_err)
		return models.User{}, wrapErr
	}

	return user, nil
}

func (r *UserRepositoryImplementation) CreateNewUser(user models.User) error {
	userCollection := r.client.Collection("users")

	_, err := userCollection.Doc(user.UID).Create(context.Background(), user)

	if err != nil {
		wrapErr := error_utils.UserRepositoryError{}
		return error_utils.CheckFirebaseError(err, user.UID, user, &wrapErr)
	}
	return nil
}
