package repository

import (
	"testing"
	"time"

	"cloud.google.com/go/firestore"
	"github.com/AndresCRamos/midas-app-api/models"
	error_utils "github.com/AndresCRamos/midas-app-api/utils/errors"
	test_utils "github.com/AndresCRamos/midas-app-api/utils/test"
	"github.com/stretchr/testify/assert"
)

func Test_movementRepositoryImplementation_CreateNewMovement(t *testing.T) {
	firestoreClient := test_utils.InitTestingFireStore(t)
	firestoreClientFail := test_utils.InitTestingFireStoreFail(t)

	createTestOwner(t, firestoreClient)
	sourceID := createTestSource(t, firestoreClient)

	tests := []test_utils.TestCase{
		{
			Name: "Success",
			Fields: test_utils.Fields{
				"firestoreClient": firestoreClient,
			},
			Args: test_utils.Args{
				"movement": models.Movement{
					UID:       "0",
					Name:      "TestMovement",
					OwnerId:   "0",
					SourceID:  sourceID,
					CreatedAt: time.Now(),
					UpdatedAt: time.Now(),
				},
			},
			WantErr:     false,
			ExpectedErr: nil,
			PreTest:     nil,
		},
		{
			Name: "Fail to connect",
			Fields: test_utils.Fields{
				"firestoreClient": firestoreClientFail,
			},
			Args: test_utils.Args{
				"movement": models.Movement{},
			},
			WantErr:     true,
			ExpectedErr: error_utils.FirebaseUnknownError{},
			PreTest:     nil,
		},
		{
			Name: "Cant find owner",
			Fields: test_utils.Fields{
				"firestoreClient": firestoreClient,
			},
			Args: test_utils.Args{
				"movement": models.Movement{UID: "0", Name: "TestMovement", OwnerId: "1", SourceID: sourceID},
			},
			WantErr:     true,
			ExpectedErr: error_utils.MovementOwnerNotFound{MovementID: "0", OwnerId: "1"},
			PreTest:     nil,
		},
		{
			Name: "Cant find source",
			Fields: test_utils.Fields{
				"firestoreClient": firestoreClient,
			},
			Args: test_utils.Args{
				"movement": models.Movement{UID: "0", Name: "TestMovement", OwnerId: "0", SourceID: "not_found"},
			},
			WantErr:     true,
			ExpectedErr: error_utils.MovementSourceNotFound{MovementID: "0", SourceID: "not_found"},
			PreTest:     nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.Name, func(t *testing.T) {

			if tt.PreTest != nil {
				tt.PreTest(t)
			}
			testFirestoreClient := test_utils.GetFieldByNameAndType[*firestore.Client](t, tt.Fields, "firestoreClient")
			r := &movementRepositoryImplementation{
				client: testFirestoreClient,
			}
			movementTest := test_utils.GetArgByNameAndType[models.Movement](t, tt.Args, "movement")
			res, err := r.CreateNewMovement(movementTest)
			if !tt.WantErr {
				if assert.NoError(t, err) {
					checkEqualMovement(t, movementTest, res)
				}

			} else {
				assert.ErrorAs(t, err, &tt.ExpectedErr, "Expected error as: %s", tt.ExpectedErr.Error())
			}
			defer func() {
				deleteTestMovement(firestoreClient, res.UID)
			}()
		})
	}
	deleteTestUser(firestoreClient)
	deleteTestSource(firestoreClient, sourceID)
}

func deleteTestMovement(client *firestore.Client, id string) {
	args := map[string]interface{}{
		"Collection": "movements",
		"id":         id,
	}
	test_utils.ClearFireStoreTest(client, "Create", args)
}

func checkEqualMovement(t *testing.T, expected models.Movement, got models.Movement) {
	assert.Equal(t, expected.OwnerId, got.OwnerId)
	assert.Equal(t, expected.SourceID, got.SourceID)
	assert.Equal(t, expected.Name, got.Name)
	assert.Equal(t, expected.Description, got.Description)
	assert.Equal(t, expected.Amount, got.Amount)
	assert.Equal(t, expected.MovementDate, got.MovementDate)
	assert.Equal(t, expected.Tags, got.Tags)
	assert.WithinDuration(t, expected.CreatedAt, got.CreatedAt, 10*time.Second)
	assert.WithinDuration(t, expected.UpdatedAt, got.UpdatedAt, 10*time.Second)
}
