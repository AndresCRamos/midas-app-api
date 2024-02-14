package repository

import (
	"testing"

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
				"movement": models.Movement{UID: "0", Name: "TestMovement", OwnerId: "0", SourceID: sourceID},
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
				assert.NoError(t, err)
				assert.Equal(t, movementTest.Name, res.Name)
				assert.Equal(t, movementTest.Description, res.Description)
			} else {
				assert.ErrorAs(t, err, &tt.ExpectedErr, "Expected error as: %s", tt.ExpectedErr.Error())
			}
			defer func() {
				args := map[string]interface{}{
					"Collection": "movements",
					"id":         res.UID,
				}
				test_utils.ClearFireStoreTest(firestoreClient, "Create", args)
			}()
		})
	}
	deleteTestUser(firestoreClient)
	deleteTestSource(firestoreClient, sourceID)
}
