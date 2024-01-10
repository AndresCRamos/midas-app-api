package repository

import (
	"testing"

	"cloud.google.com/go/firestore"
	"github.com/AndresCRamos/midas-app-api/models"
	error_utils "github.com/AndresCRamos/midas-app-api/utils/errors"
	test_utils "github.com/AndresCRamos/midas-app-api/utils/test"
	"github.com/stretchr/testify/assert"
)

func createTestOwner(t *testing.T, firestoreClient *firestore.Client) {
	userDuplicated := &UserRepositoryImplementation{
		client: firestoreClient,
	}

	err := userDuplicated.CreateNewUser(models.User{UID: "0", Alias: "TEST USER"})
	if err != nil {
		t.Fatalf("Cant connect to Firestore to create test user: %s", err.Error())
	}

}

func DeleteTestUser(firestoreClient *firestore.Client) {
	args := map[string]interface{}{
		"Collection": "users",
		"id":         "0",
	}
	test_utils.ClearFireStoreTest(firestoreClient, "Create", args)
}

func TestSourceRepositoryImplementation_CreateNewSource(t *testing.T) {

	firestoreClient := test_utils.InitTestingFireStore(t)
	firestoreClientFail := test_utils.InitTestingFireStoreFail(t)

	dupSource := models.Source{
		UID:     "0",
		Name:    "DupSource",
		OwnerId: "0",
	}

	createDupSource := func(t *testing.T) {
		rDuplicated := &SourceRepositoryImplementation{
			client: firestoreClient,
		}

		err := rDuplicated.CreateNewSource(dupSource)
		if err != nil {
			t.Fatalf("Cant connect to Firestore to check for duplication test: %s", err.Error())
		}
	}

	tests := []test_utils.TestCase{
		{
			Name: "Success",
			Fields: test_utils.Fields{
				"firestoreClient": firestoreClient,
			},
			Args: test_utils.Args{
				"source": &models.Source{UID: "0", Name: "TestSource", OwnerId: "0"},
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
				"source": &models.Source{},
			},
			WantErr:     true,
			ExpectedErr: error_utils.FirebaseUnknownError{},
			PreTest:     nil,
		},
		{
			Name: "Duplicated source",
			Fields: test_utils.Fields{
				"firestoreClient": firestoreClient,
			},
			Args: test_utils.Args{
				"source": &dupSource,
			},
			WantErr:     true,
			ExpectedErr: error_utils.FirestoreAlreadyExistsError{DocID: dupSource.UID},
			PreTest:     createDupSource,
		},
		{
			Name: "Cant find owner",
			Fields: test_utils.Fields{
				"firestoreClient": firestoreClient,
			},
			Args: test_utils.Args{
				"source": &models.Source{UID: "0", Name: "TestSource", OwnerId: "1"},
			},
			WantErr:     true,
			ExpectedErr: error_utils.SourceOwnerNotFound{SourceID: "0", OwnerId: "1"},
			PreTest:     createDupSource,
		},
	}

	createTestOwner(t, firestoreClient)

	for _, tt := range tests {
		t.Run(tt.Name, func(t *testing.T) {

			if tt.PreTest != nil {
				tt.PreTest(t)
			}
			testFirestoreClient := test_utils.GetFieldByNameAndType(t, tt.Fields, "firestoreClient", new(firestore.Client))
			r := &SourceRepositoryImplementation{
				client: testFirestoreClient.(*firestore.Client),
			}
			sourceTest := test_utils.GetArgByNameAndType(t, tt.Args, "source", new(models.Source)).(*models.Source)
			err := r.CreateNewSource(*sourceTest)
			if !tt.WantErr {
				assert.NoError(t, err)
			} else {
				assert.ErrorAs(t, err, &tt.ExpectedErr, "Expected error as: %s", tt.ExpectedErr.Error())
			}
			args := map[string]interface{}{
				"Collection": "sources",
				"id":         sourceTest.UID,
			}
			test_utils.ClearFireStoreTest(firestoreClient, "Create", args)

		})
	}

	DeleteTestUser(firestoreClient)
}

func TestSourceRepositoryImplementation_GetSourceByID(t *testing.T) {

	firestoreClient := test_utils.InitTestingFireStore(t)
	firestoreClientFail := test_utils.InitTestingFireStoreFail(t)

	createTestOwner(t, firestoreClient)

	searchSource := models.Source{
		UID:     "0",
		Name:    "SourceToSearch",
		OwnerId: "0",
	}

	rSearch := SourceRepositoryImplementation{
		client: firestoreClient,
	}

	testFields := test_utils.Fields{
		"firestoreClient": firestoreClient,
	}

	testFieldsFail := test_utils.Fields{
		"firestoreClient": firestoreClientFail,
	}

	err := rSearch.CreateNewSource(searchSource)
	if err != nil {
		t.Fatalf("Cant create source to search: %s", err.Error())
	}

	tests := []test_utils.TestCase{
		{
			Name:   "Success",
			Fields: testFields,
			Args: test_utils.Args{
				"id": searchSource.UID,
			},
			WantErr:     false,
			ExpectedErr: nil,
			PreTest:     nil,
		},
		{
			Name:   "Fail to connect",
			Fields: testFieldsFail,
			Args: test_utils.Args{
				"id": searchSource.UID,
			},
			WantErr:     true,
			ExpectedErr: error_utils.FirebaseUnknownError{},
			PreTest:     nil,
		},
		{
			Name:   "Cant find",
			Fields: testFields,
			Args: test_utils.Args{
				"id": "100",
			},
			WantErr:     true,
			ExpectedErr: error_utils.FirestoreNotFoundError{DocID: "100"},
			PreTest:     nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.Name, func(t *testing.T) {
			if tt.PreTest != nil {
				tt.PreTest(t)
			}
			r := &SourceRepositoryImplementation{
				client: test_utils.GetFieldByNameAndType(t, tt.Fields, "firestoreClient", new(firestore.Client)).(*firestore.Client),
			}
			sourceTestId := test_utils.GetArgByNameAndType(t, tt.Args, "id", "").(string)

			res, err := r.GetSourceByID(sourceTestId)
			if !tt.WantErr {
				assert.NoError(t, err)
				assert.Equal(t, searchSource, res)
			} else {
				assert.ErrorAs(t, err, &tt.ExpectedErr)
			}
		})
	}
	DeleteTestUser(firestoreClient)
}
