package repository

import (
	"strconv"
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

func Test_movementRepositoryImplementation_GetMovementByID(t *testing.T) {

	firestoreClient := test_utils.InitTestingFireStore(t)
	firestoreClientFail := test_utils.InitTestingFireStoreFail(t)

	createTestOwner(t, firestoreClient)
	createdSourceUID := createTestSource(t, firestoreClient)
	createdMovement := createTestMovement(t, firestoreClient, createdSourceUID)

	testFields := test_utils.Fields{
		"firestoreClient": firestoreClient,
	}

	testFieldsFail := test_utils.Fields{
		"firestoreClient": firestoreClientFail,
	}

	tests := []test_utils.TestCase{
		{
			Name:   "Success",
			Fields: testFields,
			Args: test_utils.Args{
				"id":     createdMovement.UID,
				"userID": "0",
			},
			WantErr:     false,
			ExpectedErr: nil,
			PreTest:     nil,
		},
		{
			Name:   "Fail to connect",
			Fields: testFieldsFail,
			Args: test_utils.Args{
				"id":     "0",
				"userID": "0",
			},
			WantErr:     true,
			ExpectedErr: error_utils.FirebaseUnknownError{},
			PreTest:     nil,
		},
		{
			Name:   "Cant find",
			Fields: testFields,
			Args: test_utils.Args{
				"id":     "100",
				"userID": "0",
			},
			WantErr:     true,
			ExpectedErr: error_utils.FirestoreNotFoundError{DocID: "100"},
			PreTest:     nil,
		},
		{
			Name:   "Different user",
			Fields: testFields,
			Args: test_utils.Args{
				"id":     createdMovement.UID,
				"userID": "1",
			},
			WantErr:     true,
			ExpectedErr: error_utils.MovementDifferentOwner{MovementID: createdMovement.UID, OwnerID: "1"},
			PreTest:     nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.Name, func(t *testing.T) {
			if tt.PreTest != nil {
				tt.PreTest(t)
			}
			r := &movementRepositoryImplementation{
				client: test_utils.GetFieldByNameAndType[*firestore.Client](t, tt.Fields, "firestoreClient"),
			}
			sourceTestId := test_utils.GetArgByNameAndType[string](t, tt.Args, "id")
			userID := test_utils.GetArgByNameAndType[string](t, tt.Args, "userID")

			res, err := r.GetMovementByID(sourceTestId, userID)
			if !tt.WantErr {
				if assert.NoError(t, err) {
					checkEqualMovement(t, createdMovement, res)
				}
			} else {
				assert.ErrorAs(t, err, &tt.ExpectedErr)
			}
		})
	}
	deleteTestMovement(firestoreClient, createdMovement.UID)
	deleteTestSource(firestoreClient, createdSourceUID)
	deleteTestUser(firestoreClient)
}

func createTestMovement(t *testing.T, client *firestore.Client, sourceID string) models.Movement {
	testMovement := movementRepositoryImplementation{
		client: client,
	}

	res, err := testMovement.CreateNewMovement(models.Movement{
		Name:     "test Owner",
		SourceID: sourceID,
		OwnerId:  "0",
	})

	if err != nil {
		t.Fatalf("Cant connect to Firestore to create test source: %s", err.Error())
	}

	return res
}

func createTestMovementList(t *testing.T, client *firestore.Client, sourceID string) []models.Movement {
	movementRepo := movementRepositoryImplementation{
		client: client,
	}
	var movementList []models.Movement

	for i := 0; i < 52; i++ {
		createdMovement, err := movementRepo.CreateNewMovement(models.Movement{
			Name:         "Test movement N" + strconv.Itoa(i),
			OwnerId:      "0",
			SourceID:     sourceID,
			MovementDate: time.Now().AddDate(0, 0, i),
		})
		if err != nil {
			t.Fatalf("Cant connect to Firestore to create test movement: %s", err.Error())
		}
		movementList = append(movementList, createdMovement)
	}
	return movementList
}

func deleteTestMovementList(client *firestore.Client, movements []models.Movement) {
	for _, movement := range movements {
		deleteTestMovement(client, movement.UID)
	}
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
