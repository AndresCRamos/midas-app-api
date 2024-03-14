package firestore

import (
	"time"

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

	TestSourceUpdate = models.SourceUpdate{
		Name:        "TEST_SOURCE",
		Description: "TEST DESCRIPTION",
	}

	TestSource = models.Source{
		Name:        "TEST_SOURCE",
		Description: "TEST DESCRIPTION",
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}
)
