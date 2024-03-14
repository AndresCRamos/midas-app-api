package repository

import (
	"testing"
	"time"

	"cloud.google.com/go/firestore"
	"github.com/AndresCRamos/midas-app-api/models"
	error_utils "github.com/AndresCRamos/midas-app-api/utils/errors"
	test_utils "github.com/AndresCRamos/midas-app-api/utils/test"
	firestore_utils "github.com/AndresCRamos/midas-app-api/utils/test/firestore"
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

func createTestSource(t *testing.T, firestoreClient *firestore.Client) string {
	userDuplicated := &SourceRepositoryImplementation{
		client: firestoreClient,
	}

	res, err := userDuplicated.CreateNewSource(models.Source{
		UID:       "0",
		Name:      "Test Source",
		OwnerId:   "0",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	})
	if err != nil {
		t.Fatalf("Cant connect to Firestore to create test source: %s", err.Error())
	}
	return res.UID
}

func deleteTestSource(firestoreClient *firestore.Client, id string) {
	args := map[string]interface{}{
		"Collection": "sources",
		"id":         id,
	}
	test_utils.ClearFireStoreTest(firestoreClient, "Create", args)
}

func deleteTestUser(firestoreClient *firestore.Client) {
	args := map[string]interface{}{
		"Collection": "users",
		"id":         "0",
	}
	test_utils.ClearFireStoreTest(firestoreClient, "Create", args)
}

func TestSourceRepositoryImplementation_GetMovementsBySourceAndDate(t *testing.T) {

	firestoreClient := test_utils.InitTestingFireStore(t)
	firestoreClientFail := test_utils.InitTestingFireStoreFail(t)

	createdOwner := firestore_utils.CreateTestUser(t, firestoreClient, "0")
	createdSourceUID := firestore_utils.CreateTestSource(t, firestoreClient, createdOwner.UID)
	createdMovements := createTestMovementList(t, firestoreClient, createdSourceUID.UID)

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
				"userID":           "0",
				"sourceID":         createdSourceUID.UID,
				"date_from":        time.Now().UTC().Add(-60 * 24 * time.Hour),
				"date_to":          time.Now().UTC(),
				"page":             1,
				"expectedPageSize": 50,
			},
			WantErr:     false,
			ExpectedErr: nil,
			PreTest:     nil,
		},
		{
			Name:   "Fail to connect",
			Fields: testFieldsFail,
			Args: test_utils.Args{
				"userID":           "0",
				"sourceID":         createdSourceUID.UID,
				"date_from":        time.Now().UTC().Add(-60 * 24 * time.Hour),
				"date_to":          time.Now().UTC(),
				"page":             1,
				"expectedPageSize": 50,
			},
			WantErr:     true,
			ExpectedErr: error_utils.FirebaseUnknownError{},
			PreTest:     nil,
		},
		{
			Name:   "Pagination test",
			Fields: testFields,
			Args: test_utils.Args{
				"userID":           "0",
				"sourceID":         createdSourceUID.UID,
				"date_from":        time.Now().UTC().Add(-60 * 24 * time.Hour),
				"date_to":          time.Now().UTC(),
				"page":             2,
				"expectedPageSize": 2,
			},
			WantErr:     false,
			ExpectedErr: nil,
			PreTest:     nil,
		},
		{
			Name:   "Not enough data",
			Fields: testFields,
			Args: test_utils.Args{
				"userID":    "0",
				"sourceID":  createdSourceUID.UID,
				"date_from": time.Now().UTC().Add(-60 * 24 * time.Hour),
				"date_to":   time.Now().UTC(),
				"page":      3,
			},
			WantErr:     true,
			ExpectedErr: error_utils.MovementNotEnoughData{},
			PreTest:     nil,
		},
		{
			Name:   "Source not found",
			Fields: testFields,
			Args: test_utils.Args{
				"userID":    "0",
				"sourceID":  "0",
				"date_from": time.Now().UTC().Add(-60 * 24 * time.Hour),
				"date_to":   time.Now().UTC(),
				"page":      1,
			},
			WantErr:     true,
			ExpectedErr: error_utils.SourceNotFound{},
			PreTest:     nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.Name, func(t *testing.T) {
			if tt.PreTest != nil {
				tt.PreTest(t)
			}
			r := &SourceRepositoryImplementation{
				client: test_utils.GetFieldByNameAndType[*firestore.Client](t, tt.Fields, "firestoreClient"),
			}
			userID := test_utils.GetArgByNameAndType[string](t, tt.Args, "userID")
			sourceID := test_utils.GetArgByNameAndType[string](t, tt.Args, "sourceID")
			page := test_utils.GetArgByNameAndType[int](t, tt.Args, "page")
			date_from := test_utils.GetArgByNameAndType[time.Time](t, tt.Args, "date_from")
			date_to := test_utils.GetArgByNameAndType[time.Time](t, tt.Args, "date_to")

			res, err := r.GetMovementsBySourceAndDate(sourceID, userID, page, date_from, date_to)
			if !tt.WantErr {
				expectedPageSize := test_utils.GetArgByNameAndType[int](t, tt.Args, "expectedPageSize")
				assert.NoError(t, err)
				assert.Equal(t, expectedPageSize, res.PageSize)
				assert.NotEmpty(t, res)
			} else {
				assert.Error(t, err)
				assert.ErrorAs(t, err, &tt.ExpectedErr, "Wanted: %v\nGot: %v", tt.ExpectedErr, err)
			}
		})
	}
	deleteTestMovementList(firestoreClient, createdMovements)
	firestore_utils.DeleteTestSource(t, firestoreClient, createdSourceUID.UID)
	firestore_utils.DeleteTestUser(t, firestoreClient, createdOwner.UID)
}

func TestSourceRepositoryImplementation_CreateNewSource(t *testing.T) {

	firestoreClient := test_utils.InitTestingFireStore(t)
	firestoreClientFail := test_utils.InitTestingFireStoreFail(t)

	tests := []test_utils.TestCase{
		{
			Name: "Success",
			Fields: test_utils.Fields{
				"firestoreClient": firestoreClient,
			},
			Args: test_utils.Args{
				"source": models.Source{UID: "0", Name: "TestSource", OwnerId: "0"},
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
				"source": models.Source{},
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
				"source": models.Source{UID: "0", Name: "TestSource", OwnerId: "1"},
			},
			WantErr:     true,
			ExpectedErr: error_utils.SourceOwnerNotFound{SourceID: "0", OwnerId: "1"},
			PreTest:     nil,
		},
	}

	createdOwner := firestore_utils.CreateTestUser(t, firestoreClient, "0")

	for _, tt := range tests {
		t.Run(tt.Name, func(t *testing.T) {

			if tt.PreTest != nil {
				tt.PreTest(t)
			}
			testFirestoreClient := test_utils.GetFieldByNameAndType[*firestore.Client](t, tt.Fields, "firestoreClient")
			r := &SourceRepositoryImplementation{
				client: testFirestoreClient,
			}
			sourceTest := test_utils.GetArgByNameAndType[models.Source](t, tt.Args, "source")
			res, err := r.CreateNewSource(sourceTest)
			if !tt.WantErr {
				assert.NoError(t, err)
				assert.Equal(t, sourceTest.Name, res.Name)
				assert.Equal(t, sourceTest.Description, res.Description)
			} else {
				assert.ErrorAs(t, err, &tt.ExpectedErr, "Expected error as: %s", tt.ExpectedErr.Error())
			}
			defer func() {
				args := map[string]interface{}{
					"Collection": "sources",
					"id":         res.UID,
				}
				test_utils.ClearFireStoreTest(firestoreClient, "Create", args)
			}()
		})
	}
	firestore_utils.DeleteTestUser(t, firestoreClient, createdOwner.UID)
}

func TestSourceRepositoryImplementation_GetSourcesByUser(t *testing.T) {

	firestoreClient := test_utils.InitTestingFireStore(t)
	firestoreClientFail := test_utils.InitTestingFireStoreFail(t)

	createdOwner := firestore_utils.CreateTestUser(t, firestoreClient, "0")
	createdSources := firestore_utils.CreateTestSourceList(t, firestoreClient, createdOwner.UID)

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
				"userID":           "0",
				"page":             1,
				"expectedPageSize": 50,
			},
			WantErr:     false,
			ExpectedErr: nil,
			PreTest:     nil,
		},
		{
			Name:   "Fail to connect",
			Fields: testFieldsFail,
			Args: test_utils.Args{
				"userID": "0",
				"page":   1,
			},
			WantErr:     true,
			ExpectedErr: error_utils.FirebaseUnknownError{},
			PreTest:     nil,
		},
		{
			Name:   "Pagination test",
			Fields: testFields,
			Args: test_utils.Args{
				"userID":           "0",
				"page":             2,
				"expectedPageSize": 1,
			},
			WantErr:     false,
			ExpectedErr: nil,
			PreTest:     nil,
		},
		{
			Name:   "Not enough data",
			Fields: testFields,
			Args: test_utils.Args{
				"userID": "0",
				"page":   3,
			},
			WantErr:     true,
			ExpectedErr: nil,
			PreTest:     nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.Name, func(t *testing.T) {
			if tt.PreTest != nil {
				tt.PreTest(t)
			}
			r := &SourceRepositoryImplementation{
				client: test_utils.GetFieldByNameAndType[*firestore.Client](t, tt.Fields, "firestoreClient"),
			}
			userID := test_utils.GetArgByNameAndType[string](t, tt.Args, "userID")
			page := test_utils.GetArgByNameAndType[int](t, tt.Args, "page")

			res, err := r.GetSourcesByUser(userID, page)
			if !tt.WantErr {
				assert.NoError(t, err)
				assert.NotEmpty(t, res)
				expectedPageSize := test_utils.GetArgByNameAndType[int](t, tt.Args, "expectedPageSize")
				assert.Equal(t, expectedPageSize, res.PageSize)
			} else {
				assert.Error(t, err)
				assert.ErrorAs(t, err, &tt.ExpectedErr, "Wanted: %v\nGot: %v", tt.ExpectedErr, err)
			}
		})
	}
	firestore_utils.DeleteTestSourceList(t, firestoreClient, createdSources)
	firestore_utils.DeleteTestUser(t, firestoreClient, createdOwner.UID)
}

func TestSourceRepositoryImplementation_GetSourceByID(t *testing.T) {

	firestoreClient := test_utils.InitTestingFireStore(t)
	firestoreClientFail := test_utils.InitTestingFireStoreFail(t)

	createdOwner := firestore_utils.CreateTestUser(t, firestoreClient, "0")
	createdSourceUID := firestore_utils.CreateTestSource(t, firestoreClient, createdOwner.UID)

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
				"id":     createdSourceUID.UID,
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
				"id":     createdSourceUID.UID,
				"userID": "1",
			},
			WantErr:     true,
			ExpectedErr: error_utils.SourceDifferentOwner{SourceID: createdSourceUID.UID, OwnerID: "1"},
			PreTest:     nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.Name, func(t *testing.T) {
			if tt.PreTest != nil {
				tt.PreTest(t)
			}
			r := &SourceRepositoryImplementation{
				client: test_utils.GetFieldByNameAndType[*firestore.Client](t, tt.Fields, "firestoreClient"),
			}
			sourceTestId := test_utils.GetArgByNameAndType[string](t, tt.Args, "id")
			userID := test_utils.GetArgByNameAndType[string](t, tt.Args, "userID")

			res, err := r.GetSourceByID(sourceTestId, userID)
			if !tt.WantErr {
				assert.NoError(t, err)
				checkEqualSource(t, firestore_utils.TestSource, res)
			} else {
				assert.ErrorAs(t, err, &tt.ExpectedErr)
			}
		})
	}
	firestore_utils.DeleteTestSource(t, firestoreClient, createdSourceUID.UID)
	firestore_utils.DeleteTestUser(t, firestoreClient, createdOwner.UID)
}

func TestSourceRepositoryImplementation_UpdateSource(t *testing.T) {

	firestoreClient := test_utils.InitTestingFireStore(t)
	firestoreClientFail := test_utils.InitTestingFireStoreFail(t)

	createdOwner := firestore_utils.CreateTestUser(t, firestoreClient, "0")
	createdSourceUID := firestore_utils.CreateTestSource(t, firestoreClient, createdOwner.UID)

	tests := []test_utils.TestCase{
		{
			Name: "Success",
			Fields: test_utils.Fields{
				"firestoreClient": firestoreClient,
			},
			Args: test_utils.Args{
				"source": firestore_utils.SetTestSourceData(firestore_utils.TestSourceUpdated, createdSourceUID.UID, createdOwner.UID),
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
				"source": models.Source{},
			},
			WantErr:     true,
			ExpectedErr: error_utils.FirebaseUnknownError{},
			PreTest:     nil,
		},
		{
			Name: "Cant find",
			Fields: test_utils.Fields{
				"firestoreClient": firestoreClient,
			},
			Args: test_utils.Args{
				"source": models.Source{UID: "100", Name: "Not found Source"},
			},
			WantErr:     true,
			ExpectedErr: error_utils.FirestoreNotFoundError{DocID: createdSourceUID.UID},
			PreTest:     nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.Name, func(t *testing.T) {

			if tt.PreTest != nil {
				tt.PreTest(t)
			}
			testFirestoreClient := test_utils.GetFieldByNameAndType[*firestore.Client](t, tt.Fields, "firestoreClient")
			r := &SourceRepositoryImplementation{
				client: testFirestoreClient,
			}
			sourceTest := test_utils.GetArgByNameAndType[models.Source](t, tt.Args, "source")
			_, err := r.UpdateSource(sourceTest)
			if !tt.WantErr {
				assert.NoError(t, err)
			} else {
				assert.Error(t, err, "Expected: %s", tt.ExpectedErr.Error())
				if err != nil {
					assert.ErrorAs(
						t, err, &tt.ExpectedErr, "Expected error as:\n%s\nGot:\n%s",
						tt.ExpectedErr.Error(),
						err.Error(),
					)
				}

			}
			args := map[string]interface{}{
				"Collection":   "sources",
				"id":           createdSourceUID.UID,
				"originalData": firestore_utils.TestSource,
			}
			test_utils.ClearFireStoreTest(firestoreClient, "Update", args)

		})
	}
	defer func() {
		firestore_utils.DeleteTestSource(t, firestoreClient, createdSourceUID.UID)
		firestore_utils.DeleteTestUser(t, firestoreClient, createdOwner.UID)
	}()
}

func TestSourceRepositoryImplementation_DeleteSource(t *testing.T) {

	firestoreClient := test_utils.InitTestingFireStore(t)
	firestoreClientFail := test_utils.InitTestingFireStoreFail(t)

	createdOwner := firestore_utils.CreateTestUser(t, firestoreClient, "0")
	createdSourceUID := firestore_utils.CreateTestSource(t, firestoreClient, createdOwner.UID)

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
				"id": createdSourceUID.UID,
			},
			WantErr:     false,
			ExpectedErr: nil,
			PreTest:     nil,
		},
		{
			Name:   "Fail to connect",
			Fields: testFieldsFail,
			Args: test_utils.Args{
				"id": "0",
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
				client: test_utils.GetFieldByNameAndType[*firestore.Client](t, tt.Fields, "firestoreClient"),
			}
			sourceTestId := test_utils.GetArgByNameAndType[string](t, tt.Args, "id")

			err := r.DeleteSource(sourceTestId, "0")
			if !tt.WantErr {
				assert.NoError(t, err)
			} else {
				assert.ErrorAs(t, err, &tt.ExpectedErr)
			}
		})
	}
	firestore_utils.DeleteTestUser(t, firestoreClient, createdOwner.UID)
}

func checkEqualSource(t *testing.T, expected models.Source, got models.Source) {
	assert.Equal(t, expected.Name, got.Name)
	assert.Equal(t, expected.Description, got.Description)
	assert.WithinDuration(t, expected.CreatedAt, got.CreatedAt, 10*time.Second)
	assert.WithinDuration(t, expected.UpdatedAt, got.UpdatedAt, 10*time.Second)
}
