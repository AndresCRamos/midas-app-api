package firestore

import (
	"context"
	"testing"

	"cloud.google.com/go/firestore"
	"github.com/AndresCRamos/midas-app-api/models"
)

func CreateTestUser(t *testing.T, client *firestore.Client, uid string) models.User {
	testUser := models.User{
		UID:      uid,
		Alias:    "TEST_USER",
		Name:     "John",
		LastName: "Doe",
	}

	_, err := client.Collection("user").Doc(uid).Set(context.Background(), testUser)

	if err != nil {
		t.Fatalf("Can't create test user: %s", err)
	}

	return testUser
}
