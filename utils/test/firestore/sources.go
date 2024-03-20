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

	TestSourceUpdated = models.Source{
		Name:        "TEST_SOURCE",
		Description: "TEST DESCRIPTION UPDATED",
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

func SetTestSourceData(source models.Source, uid string, ownerID string) models.Source {
	source.UID = uid
	source.OwnerId = ownerID
	return source
}

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

func createTestSourceListItem(t *testing.T, client *firestore.Client, ownerID string, n int) models.Source {
	sourceDocRef := client.Collection("sources").NewDoc()
	tUser := TestSource
	tUser.UID = sourceDocRef.ID
	tUser.OwnerId = ownerID
	tUser.Name += "_N" + fmt.Sprint(n)
	_, err := sourceDocRef.Set(context.Background(), tUser)
	if err != nil {
		t.Fatalf("Can't create test source: %s", err)
	}
	return tUser
}

func CreateTestSourceList(t *testing.T, client *firestore.Client, ownerID string) []models.Source {
	createdList := []models.Source{}
	for i := 0; i < 51; i++ {
		createdList = append(createdList, createTestSourceListItem(t, client, ownerID, i))
	}

	return createdList
}

func DeleteTestSourceList(t *testing.T, client *firestore.Client, deleteSources []models.Source) {
	for _, source := range deleteSources {
		DeleteTestSource(t, client, source.UID)
	}
}

func DeleteTestSource(t *testing.T, client *firestore.Client, uid string) {
	_, err := client.Collection("sources").Doc(uid).Delete(context.Background())

	if err != nil {
		t.Logf("Cant delete test user: %s", err.Error())
	}
}
