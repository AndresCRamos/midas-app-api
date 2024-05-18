package compare

import (
	"encoding/json"
	"testing"
	"time"

	"github.com/AndresCRamos/midas-app-api/models"
)

func CompareMovements(expected models.Movement, got models.Movement, checkID bool, maxDelta time.Duration) bool {
	if checkID {
		if expected.UID != got.UID {
			return false
		}
	}
	if expected.OwnerId != got.OwnerId {
		return false
	}
	if expected.SourceID != got.SourceID {
		return false
	}
	if expected.Name != got.Name {
		return false
	}
	if expected.Description != got.Description {
		return false
	}
	if expected.Amount != got.Amount {
		return false
	}
	if !expected.MovementDate.Equal(got.MovementDate) {
		return false
	}
	if !CompareSlices(expected.Tags, got.Tags) {
		return false
	}

	delta := expected.CreatedAt.Sub(got.CreatedAt)
	if delta > maxDelta {
		return false
	}

	delta = expected.UpdatedAt.Sub(got.UpdatedAt)
	return delta <= maxDelta
}

func ContainsMovement(t *testing.T, expectedList []models.Movement, got models.Movement, delta time.Duration) {
	for _, elem := range expectedList {
		if CompareMovements(elem, got, true, delta) {
			return
		}
	}
	bList, _ := json.MarshalIndent(expectedList, "", " ")
	dList, _ := json.MarshalIndent(got, "", " ")
	t.Fatalf("List\n%v\ndoes not contain\n%v", string(bList), string(dList))
}
