package test

import (
	"context"
	"log"
	"os"
	"testing"

	"cloud.google.com/go/firestore"
	firebase "firebase.google.com/go/v4"
)

const (
	fireStoreEmulatorHostEnvVar  = "FIRESTORE_EMULATOR_HOST"
	firestoreEmulatorHostDefault = "127.0.0.1:8080"
)

func InitTestingFireStore(t *testing.T) *firestore.Client {

	firebaseProject := os.Getenv("FIREBASE_PROJECT_ID")

	fireStoreEmulatorHost := os.Getenv(fireStoreEmulatorHostEnvVar)

	if fireStoreEmulatorHost == "" {
		t.Logf("Cant use empty Firestore Emulator Host, setting default value: %s", firestoreEmulatorHostDefault)
		os.Setenv(fireStoreEmulatorHostEnvVar, firestoreEmulatorHostDefault)
	} else {
		t.Logf("Firebase Emulator Host set to: %s", fireStoreEmulatorHost)
	}

	if firebaseProject == "" {
		t.Fatal("Firebase Project can not be empty")
	}
	t.Logf("Firebase Project set to: %s", firebaseProject)

	ctx := context.Background()

	conf := firebase.Config{
		ProjectID: firebaseProject,
	}

	firebaseApp, err := firebase.NewApp(ctx, &conf)

	if err != nil {
		t.Fatalf("Cant create firebase client: %s", err.Error())
	}

	// First, we try to initialize a Firestore client to check for a possible error
	firestoreClient, err := firebaseApp.Firestore(ctx)

	if err != nil {
		t.Fatalf("Cant create firestore client: %s", err.Error())
	}

	// Then a dummy document is requested to force the initialization and check if the process was successful
	// This presupposes the existence of the dummy collection and a document inside it with id 0, if thats not the case, it will always fail
	_, err = firestoreClient.Collection("dummy").Doc("0").Get(ctx)

	if err != nil {
		log.Fatalf(`
		Cant check connection to Firestore emulator client, make sure a document exists in "dummy/0":
		%s
		`, err.Error())
	}
	return firestoreClient
}

func InitTestingFireStoreFail(t *testing.T) *firestore.Client {

	firebaseProject := "non_existent"

	ctx := context.Background()

	conf := firebase.Config{
		ProjectID: firebaseProject,
	}

	firebaseApp, err := firebase.NewApp(ctx, &conf)

	if err != nil {
		t.Fatalf("Cant create firebase client: %s", err.Error())
	}

	// Since the project is null, this will always fail to connect
	firestoreClient, _ := firebaseApp.Firestore(ctx)

	_, err = firestoreClient.Collection("dummy").Doc("0").Get(ctx)

	if err == nil {
		log.Fatal(`Cant create Firestore Fail client"`)
	}

	return firestoreClient
}

func ClearFireStoreTest(client *firestore.Client, operation string, args map[string]interface{}) {
	ctx := context.Background()
	if operation == "Create" {
		collectionDelete := args["Collection"].(string)
		deleteID := args["id"].(string)
		_, _ = client.Collection(collectionDelete).Doc(deleteID).Delete(ctx)
	}
	if operation == "Update" {
		collectionDelete := args["Collection"].(string)
		updateID := args["id"].(string)
		originalData := args["originalData"]
		_, _ = client.Collection(collectionDelete).Doc(updateID).Set(ctx, originalData)
	}

}
