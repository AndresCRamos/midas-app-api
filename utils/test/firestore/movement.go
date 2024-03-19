package firestore

import (
	"context"
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
