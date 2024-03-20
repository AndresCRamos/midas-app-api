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

func SetTestMovementData(movement models.Movement, uid string, ownerID string, sourceID string, movementDate time.Time) models.Movement {
	movement.UID = uid
	movement.OwnerId = ownerID
	movement.SourceID = sourceID
	movement.MovementDate = movementDate.UTC().Truncate(24 * time.Hour)
	return movement
}

func CreateTestMovement(t *testing.T, client *firestore.Client, ownerID string, sourceID string) models.Movement {
	movementDocRef := client.Collection("movements").NewDoc()
	tMovement := SetTestMovementData(TestMovement, movementDocRef.ID, ownerID, sourceID, time.Now())
	_, err := movementDocRef.Set(context.Background(), tMovement)
	if err != nil {
		t.Fatalf("Can't create test movement: %s", err)
	}
	return tMovement
}

func createTestMovementListItem(t *testing.T, client *firestore.Client, ownerID string, sourceID string, n int) models.Movement {
	movementDocRef := client.Collection("movements").NewDoc()
	tMovement := SetTestMovementData(TestMovement, movementDocRef.ID, ownerID, sourceID, time.Now())
	tMovement.Name += "_N" + fmt.Sprint(n)

	_, err := movementDocRef.Set(context.Background(), tMovement)
	if err != nil {
		t.Fatalf("Can't create test movement: %s", err)
	}
	return tMovement
}

func CreateTestMovementList(t *testing.T, client *firestore.Client, ownerID string, sourceID string) []models.Movement {
	createdList := []models.Movement{}
	for i := 0; i < 51; i++ {
		createdList = append(createdList, createTestMovementListItem(t, client, ownerID, sourceID, i))
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
