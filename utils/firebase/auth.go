package firebase

import (
	"context"
	"fmt"
	"os"

	firebase "firebase.google.com/go/v4"
	"firebase.google.com/go/v4/auth"
	error_const "github.com/AndresCRamos/midas-app-api/utils/errors"
)

func GetFirebaseAuthClient() (*auth.Client, error) {

	firebaseProject := os.Getenv("FIREBASE_PROJECT_ID")

	if firebaseProject == "" {
		return nil, error_const.EMPTY_PROJECT
	}

	ctx := context.Background()

	conf := firebase.Config{
		ProjectID: firebaseProject,
	}

	firebaseApp, err := firebase.NewApp(ctx, &conf)

	if err != nil {
		return nil, fmt.Errorf(error_const.FIREBASE_ERROR, err)
	}

	firebaseAuthClient, err := firebaseApp.Auth(ctx)

	if err != nil {
		return nil, fmt.Errorf(error_const.AUTH_ERROR, err)
	}

	dummyUser := auth.UserToCreate{}

	_, err = firebaseAuthClient.CreateUser(ctx, &dummyUser)

	if err != nil {
		return nil, fmt.Errorf(error_const.AUTH_ERROR, error_const.FIREBASE_AUTH_CANT_CONNECT)
	}

	return firebaseAuthClient, nil
}
