package firestore

import (
	"context"
	"testing"
	"time"

	"cloud.google.com/go/firestore"
	"github.com/AndresCRamos/midas-app-api/models"
)

var (
	TestSourceCreate = models.SourceCreate{
		Name:        "TEST_SOURCE",
		Description: "TEST DESCRIPTION",
	}

	TestSourceRetrieve = models.SourceRetrieve{
		Name:        "TEST_SOURCE",
		Description: "TEST DESCRIPTION",
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	TestSourceRetrieveUpdated = models.SourceRetrieve{
		Name:        "TEST_SOURCE",
		Description: "TEST DESCRIPTION UPDATED",
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	TestSourceUpdate = models.SourceUpdate{
		Name:        "TEST_SOURCE",
		Description: "TEST DESCRIPTION UPDATED",
	}

	TestSource = models.Source{
		Name:        "TEST_SOURCE",
		Description: "TEST DESCRIPTION",
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}
)

func CreateTestSource(t *testing.T, client *firestore.Client, ownerID string) models.Source {
	sourceDocRef := client.Collection("sources").NewDoc()
	tUser := TestSource
	tUser.UID = sourceDocRef.ID
	tUser.OwnerId = ownerID
	_, err := sourceDocRef.Set(context.Background(), tUser)
	if err != nil {
		t.Fatalf("Can't create test source: %s", err)
	}
	return tUser
}
