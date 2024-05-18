package compare

import (
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
