package firebase

import (
	"context"
	"os"

	"cloud.google.com/go/firestore"
	firebase "firebase.google.com/go/v4"
	error_const "github.com/AndresCRamos/midas-app-api/utils/errors"
)

func GetFireStoreClient() (*firestore.Client, error) {

	firebaseProject := os.Getenv("FIREBASE_PROJECT_ID")

	if firebaseProject == "" {
		return nil, error_const.EmptyProject{}
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

	// First, we try to initialize a Firestore client to check for a possible error
	firestoreClient, err := firebaseApp.Firestore(ctx)

	if err != nil {
		fsError := error_const.FirestoreError{}
		fsError.Wrap(err)
		return nil, fsError
	}

	// Then a dummy document is requested to force the initialization and check if the process was successful
	// This presupposes the existence of the dummy collection and a document inside it with id 0, if thats not the case, it will always fail
	_, err = firestoreClient.Collection("dummy").Doc("0").Get(ctx)

	if err != nil {
		fsError := error_const.FirestoreError{}
		fsError.Wrap(err)
		return nil, fsError
	}
	return firestoreClient, nil
}
