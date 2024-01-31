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

func createTestSourceList(t *testing.T, firestoreClient *firestore.Client) []string {
	userDuplicated := &SourceRepositoryImplementation{
		client: firestoreClient,
	}

	createdIDs := []string{}

	for i := 0; i < 51; i++ {
		createdSource, err := userDuplicated.CreateNewSource(models.Source{
			UID:       "0",
			Name:      "Test Source N" + strconv.Itoa(i),
			OwnerId:   "0",
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		})
		if err != nil {
			t.Fatalf("Cant connect to Firestore to create test source: %s", err.Error())
		}

		createdIDs = append(createdIDs, createdSource.UID)

	}
	return createdIDs
}

func deleteTestSourceList(firestoreClient *firestore.Client, idList []string) {
	for _, id := range idList {
		args := map[string]interface{}{
			"Collection": "sources",
			"id":         id,
		}
		test_utils.ClearFireStoreTest(firestoreClient, "Create", args)
	}
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
			Name: "Cant find owner",
			Fields: test_utils.Fields{
				"firestoreClient": firestoreClient,
			},
			Args: test_utils.Args{
				"source": &models.Source{UID: "0", Name: "TestSource", OwnerId: "1"},
			},
			WantErr:     true,
			ExpectedErr: error_utils.SourceOwnerNotFound{SourceID: "0", OwnerId: "1"},
			PreTest:     nil,
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
			res, err := r.CreateNewSource(*sourceTest)
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
	deleteTestUser(firestoreClient)
}

func TestSourceRepositoryImplementation_GetSourcesByUser(t *testing.T) {

	firestoreClient := test_utils.InitTestingFireStore(t)
	firestoreClientFail := test_utils.InitTestingFireStoreFail(t)

	createTestOwner(t, firestoreClient)
	createdSources := createTestSourceList(t, firestoreClient)

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
				client: test_utils.GetFieldByNameAndType(t, tt.Fields, "firestoreClient", new(firestore.Client)).(*firestore.Client),
			}
			userID := test_utils.GetArgByNameAndType(t, tt.Args, "userID", "").(string)
			page := test_utils.GetArgByNameAndType(t, tt.Args, "page", 0).(int)

			res, err := r.GetSourcesByUser(userID, page)
			if !tt.WantErr {
				assert.NoError(t, err)
				assert.NotEmpty(t, res)
				expectedPageSize := test_utils.GetArgByNameAndType(t, tt.Args, "expectedPageSize", 0).(int)
				assert.Equal(t, expectedPageSize, res.PageSize)
			} else {
				assert.Error(t, err)
				assert.ErrorAs(t, err, &tt.ExpectedErr, "Wanted: %v\nGot: %v", tt.ExpectedErr, err)
			}
		})
	}
	deleteTestSourceList(firestoreClient, createdSources)
	deleteTestUser(firestoreClient)
}

func TestSourceRepositoryImplementation_GetSourceByID(t *testing.T) {

	firestoreClient := test_utils.InitTestingFireStore(t)
	firestoreClientFail := test_utils.InitTestingFireStoreFail(t)

	createTestOwner(t, firestoreClient)
	createdSourceUID := createTestSource(t, firestoreClient)

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
				"id":     createdSourceUID,
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
				"id":     createdSourceUID,
				"userID": "1",
			},
			WantErr:     true,
			ExpectedErr: error_utils.SourceDifferentOwner{SourceID: createdSourceUID, OwnerID: "1"},
			PreTest:     nil,
		},
	}

	searchSource := models.Source{
		UID:       "0",
		Name:      "Test Source",
		OwnerId:   "0",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
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
			userID := test_utils.GetArgByNameAndType(t, tt.Args, "userID", "").(string)

			res, err := r.GetSourceByID(sourceTestId, userID)
			if !tt.WantErr {
				assert.NoError(t, err)
				checkEqualSource(t, searchSource, res)
			} else {
				assert.ErrorAs(t, err, &tt.ExpectedErr)
			}
		})
	}
	deleteTestSource(firestoreClient, createdSourceUID)
	deleteTestUser(firestoreClient)
}

func TestSourceRepositoryImplementation_UpdateSource(t *testing.T) {

	firestoreClient := test_utils.InitTestingFireStore(t)
	firestoreClientFail := test_utils.InitTestingFireStoreFail(t)

	originalSource := models.Source{
		UID:       "0",
		Name:      "Original Source",
		OwnerId:   "0",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	createTestOwner(t, firestoreClient)
	createdSourceUID := createTestSource(t, firestoreClient)

	tests := []test_utils.TestCase{
		{
			Name: "Success",
			Fields: test_utils.Fields{
				"firestoreClient": firestoreClient,
			},
			Args: test_utils.Args{
				"source": &models.Source{UID: createdSourceUID, Name: "Update Source", OwnerId: "0"},
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
			Name: "Cant find",
			Fields: test_utils.Fields{
				"firestoreClient": firestoreClient,
			},
			Args: test_utils.Args{
				"source": &models.Source{UID: "100", Name: "Not found Source"},
			},
			WantErr:     true,
			ExpectedErr: error_utils.FirestoreNotFoundError{DocID: originalSource.UID},
			PreTest:     nil,
		},
	}

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
			_, err := r.UpdateSource(*sourceTest)
			if !tt.WantErr {
				assert.NoError(t, err)
			} else {
				assert.Error(t, err, "Expected: %s", tt.ExpectedErr.Error())
				if err != nil {
					assert.ErrorAs(
						t, err, &tt.ExpectedErr,
						`Expected error as: 
						%s
						Got:
						%s`,
						tt.ExpectedErr.Error(),
						err.Error(),
					)
				}

			}
			args := map[string]interface{}{
				"Collection":   "sources",
				"id":           originalSource.UID,
				"originalData": originalSource,
			}
			test_utils.ClearFireStoreTest(firestoreClient, "Update", args)

		})
	}
	defer func() {
		deleteTestSource(firestoreClient, createdSourceUID)
		deleteTestUser(firestoreClient)
	}()
}

func TestSourceRepositoryImplementation_DeleteSource(t *testing.T) {

	firestoreClient := test_utils.InitTestingFireStore(t)
	firestoreClientFail := test_utils.InitTestingFireStoreFail(t)

	createTestOwner(t, firestoreClient)
	createTestSource(t, firestoreClient)

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
				"id": "0",
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
				client: test_utils.GetFieldByNameAndType(t, tt.Fields, "firestoreClient", new(firestore.Client)).(*firestore.Client),
			}
			sourceTestId := test_utils.GetArgByNameAndType(t, tt.Args, "id", "").(string)

			err := r.DeleteSource(sourceTestId, "0")
			if !tt.WantErr {
				assert.NoError(t, err)
			} else {
				assert.ErrorAs(t, err, &tt.ExpectedErr)
			}
		})
	}
	deleteTestUser(firestoreClient)
}

func checkEqualSource(t *testing.T, expected models.Source, got models.Source) {
	assert.Equal(t, expected.Name, got.Name)
	assert.Equal(t, expected.Description, got.Description)
	assert.WithinDuration(t, expected.CreatedAt, got.CreatedAt, 10*time.Second)
	assert.WithinDuration(t, expected.UpdatedAt, got.UpdatedAt, 10*time.Second)
}
