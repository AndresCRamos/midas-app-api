package firestore

import (
	"context"
	"testing"

	"cloud.google.com/go/firestore"
	"github.com/AndresCRamos/midas-app-api/models"
)

var (
	TestUser = models.User{
		Alias:    "TEST_USER",
		Name:     "John",
		LastName: "Doe",
	}
)

func SetTestUserID(uid string) models.User {
	tUser := TestUser
	tUser.UID = uid
	return tUser
}

func CreateTestUser(t *testing.T, client *firestore.Client, uid string) models.User {
	tUser := SetTestUserID(uid)

	_, err := client.Collection("users").Doc(uid).Set(context.Background(), tUser)

	if err != nil {
		t.Fatalf("Can't create test user: %s", err)
	}

	return tUser
}
