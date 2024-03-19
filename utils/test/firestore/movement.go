package firestore

import (
	"context"
	"fmt"
	"testing"
	"time"

	"cloud.google.com/go/firestore"
	"github.com/AndresCRamos/midas-app-api/models"
)

var (
	TestMovementCreate = models.MovementCreate{
		Name:        "TEST_SOURCE",
		Description: "TEST DESCRIPTION",
	}

	TestMovementRetrieve = models.MovementRetrieve{
		Name:        "TEST_SOURCE",
		Description: "TEST DESCRIPTION",
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	TestMovementUpdated = models.Movement{
		Name:        "TEST_SOURCE",
		Description: "TEST DESCRIPTION UPDATED",
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	TestMovementRetrieveUpdated = models.MovementRetrieve{
		Name:        "TEST_SOURCE",
		Description: "TEST DESCRIPTION UPDATED",
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	TestMovementUpdate = models.MovementUpdate{
		Name:        "TEST_SOURCE",
		Description: "TEST DESCRIPTION UPDATED",
	}

	TestMovement = models.Movement{
		Name:        "TEST_SOURCE",
		Description: "TEST DESCRIPTION",
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}
)

func SetTestMovementData(movement models.Movement, uid string, ownerID string) models.Movement {
	movement.UID = uid
	movement.OwnerId = ownerID
	return movement
}

func CreateTestMovement(t *testing.T, client *firestore.Client, ownerID string) models.Movement {
	movementDocRef := client.Collection("movements").NewDoc()
	tUser := TestMovement
	tUser.UID = movementDocRef.ID
	tUser.OwnerId = ownerID
	_, err := movementDocRef.Set(context.Background(), tUser)
	if err != nil {
		t.Fatalf("Can't create test movement: %s", err)
	}
	return tUser
}

func createTestMovementListItem(t *testing.T, client *firestore.Client, ownerID string, n int) models.Movement {
	movementDocRef := client.Collection("movements").NewDoc()
	tUser := TestMovement
	tUser.UID = movementDocRef.ID
	tUser.OwnerId = ownerID
	tUser.Name += "_N" + fmt.Sprint(n)
	_, err := movementDocRef.Set(context.Background(), tUser)
	if err != nil {
		t.Fatalf("Can't create test movement: %s", err)
	}
	return tUser
}

func CreateTestMovementList(t *testing.T, client *firestore.Client, ownerID string) []models.Movement {
	createdList := []models.Movement{}
	for i := 0; i < 51; i++ {
		createdList = append(createdList, createTestMovementListItem(t, client, ownerID, i))
	}

	return createdList
}

func DeleteTestMovementList(t *testing.T, client *firestore.Client, deleteMovements []models.Movement) {
	for _, movement := range deleteMovements {
		DeleteTestMovement(t, client, movement.UID)
	}
}

func DeleteTestMovement(t *testing.T, client *firestore.Client, uid string) {
	_, err := client.Collection("movements").Doc(uid).Delete(context.Background())

	if err != nil {
		t.Logf("Cant delete test user: %s", err.Error())
	}
}
