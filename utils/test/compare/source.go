package compare

import (
	"encoding/json"
	"testing"
	"time"

	"github.com/AndresCRamos/midas-app-api/models"
)

func CompareSources(expected models.Source, got models.Source, checkID bool, maxDelta time.Duration) bool {
	if checkID {
		if expected.UID != got.UID {
			return false
		}
	}
	if expected.UID != got.UID {
		return false
	}
	if expected.OwnerId != got.OwnerId {
		return false
	}
	if expected.Description != got.Description {
		return false
	}
	if expected.Name != got.Name {
		return false
	}
	if expected.Description != got.Description {
		return false
	}

	if expected.CreatedAt.Sub(got.CreatedAt).Abs() > maxDelta.Abs() {
		return false
	}

	return expected.UpdatedAt.Sub(got.UpdatedAt).Abs() <= maxDelta.Abs()
}

func ContainsSource(t *testing.T, expectedList []models.Source, got models.Source, maxDelta time.Duration) {
	for _, elem := range expectedList {
		if CompareSources(elem, got, true, maxDelta) {
			return
		}
	}
	bList, _ := json.MarshalIndent(expectedList, "", " ")
	dList, _ := json.MarshalIndent(got, "", " ")
	t.Fatalf("List\n%v\ndoes not contain\n%v", string(bList), string(dList))
}
