package firebase

import (
	"context"
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
		fbErr := error_const.FirebaseError{}
		fbErr.Wrap(err)
		return nil, fbErr
	}

	// First, we try to initialize a Firebase Auth client to check for a possible error
	firebaseAuthClient, err := firebaseApp.Auth(ctx)

	if err != nil {
		authErr := error_const.AuthError{}
		authErr.Wrap(err)
		return nil, authErr
	}

	// Then a dummy user is requested to force the initialization and check if the process was successful
	// This presupposes the existence of a dummy user with email dummy@example.com, if not, this will always fail

	_, err = firebaseAuthClient.GetUserByEmail(ctx, "dummy@example.com")

	if err != nil {
		authErr := error_const.AuthError{}
		authErr.Wrap(err)
		return nil, authErr
	}

	return firebaseAuthClient, nil
}
