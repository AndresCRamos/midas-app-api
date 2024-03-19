package firestore

import (
	"time"

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
