package repository

import (
	"context"
	"fmt"

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
		return models.User{}, error_utils.CheckFirebaseError(err, id, models.User{})
	}

	var user models.User

	if err = userDoc.DataTo(user); err != nil {
		return models.User{}, fmt.Errorf("Cant parse info for user %s", id)
	}

	return user, nil
}

func (r *UserRepositoryImplementation) CreateNewUser(user models.User) error {
	userCollection := r.client.Collection("users")

	_, err := userCollection.Doc(user.UID).Create(context.Background(), user)

	if err != nil {
		return error_utils.CheckFirebaseError(err, user.UID, user)
	}
	return nil
}
