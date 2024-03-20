package repository

import (
	"encoding/json"
	"strconv"
	"testing"
	"time"

	"cloud.google.com/go/firestore"
	"github.com/AndresCRamos/midas-app-api/models"
	error_utils "github.com/AndresCRamos/midas-app-api/utils/errors"
	test_utils "github.com/AndresCRamos/midas-app-api/utils/test"
	firestore_utils "github.com/AndresCRamos/midas-app-api/utils/test/firestore"
	"github.com/stretchr/testify/assert"
)

func Test_movementRepositoryImplementation_CreateNewMovement(t *testing.T) {
	firestoreClient := test_utils.InitTestingFireStore(t)
	firestoreClientFail := test_utils.InitTestingFireStoreFail(t)

	createdOwner := firestore_utils.CreateTestUser(t, firestoreClient, "0")
	createdSource := firestore_utils.CreateTestSource(t, firestoreClient, createdOwner.UID)

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
					SourceID:  createdSource.UID,
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
				"movement": models.Movement{UID: "0", Name: "TestMovement", OwnerId: "1", SourceID: createdSource.UID},
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
				firestore_utils.DeleteTestMovement(t, firestoreClient, res.UID)
			}()
		})
	}
	firestore_utils.DeleteTestSource(t, firestoreClient, createdSource.UID)
	firestore_utils.DeleteTestUser(t, firestoreClient, createdOwner.UID)
}

func Test_movementRepositoryImplementation_GetMovementByID(t *testing.T) {

	firestoreClient := test_utils.InitTestingFireStore(t)
	firestoreClientFail := test_utils.InitTestingFireStoreFail(t)

	createdOwner := firestore_utils.CreateTestUser(t, firestoreClient, "0")
	createdSource := firestore_utils.CreateTestSource(t, firestoreClient, createdOwner.UID)
	createdMovement := firestore_utils.CreateTestMovement(t, firestoreClient, createdOwner.UID, createdSource.UID)

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
	firestore_utils.DeleteTestMovement(t, firestoreClient, createdMovement.UID)
	firestore_utils.DeleteTestSource(t, firestoreClient, createdSource.UID)
	firestore_utils.DeleteTestUser(t, firestoreClient, createdOwner.UID)
}

func Test_movementRepositoryImplementation_GetMovementsByUserAndDate(t *testing.T) {

	firestoreClient := test_utils.InitTestingFireStore(t)
	firestoreClientFail := test_utils.InitTestingFireStoreFail(t)

	createdOwner := firestore_utils.CreateTestUser(t, firestoreClient, "0")
	createdSource := firestore_utils.CreateTestSource(t, firestoreClient, createdOwner.UID)
	createdMovements := firestore_utils.CreateTestMovementList(t, firestoreClient, createdOwner.UID, createdSource.UID)

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
				"date_from":        time.Now().Add(-100 * 24 * time.Hour),
				"date_to":          time.Now(),
			},
			WantErr:     false,
			ExpectedErr: nil,
			PreTest:     nil,
		},
		{
			Name:   "Fail to connect",
			Fields: testFieldsFail,
			Args: test_utils.Args{
				"userID":    "0",
				"page":      1,
				"date_from": time.Now().Add(-100 * 24 * time.Hour),
				"date_to":   time.Now(),
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
				"expectedPageSize": 2,
				"date_from":        time.Now().Add(-100 * 24 * time.Hour),
				"date_to":          time.Now(),
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
				"page":      3,
				"date_from": time.Now().Add(-100 * 24 * time.Hour),
				"date_to":   time.Now(),
			},
			WantErr:     true,
			ExpectedErr: nil,
			PreTest:     nil,
		},
		{
			Name:   "Date Limit",
			Fields: testFields,
			Args: test_utils.Args{
				"userID":           "0",
				"page":             1,
				"expectedPageSize": 30,
				"date_from":        time.Now().Add(-30 * 24 * time.Hour),
				"date_to":          time.Now(),
			},
			WantErr:     false,
			ExpectedErr: nil,
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
			userID := test_utils.GetArgByNameAndType[string](t, tt.Args, "userID")
			page := test_utils.GetArgByNameAndType[int](t, tt.Args, "page")
			dateFrom := test_utils.GetArgByNameAndType[time.Time](t, tt.Args, "date_from")
			dateTo := test_utils.GetArgByNameAndType[time.Time](t, tt.Args, "date_to")

			res, err := r.GetMovementsByUserAndDate(userID, page, dateFrom, dateTo)
			if !tt.WantErr {
				assert.NoError(t, err)
				assert.NotEmpty(t, res)
				expectedPageSize := test_utils.GetArgByNameAndType[int](t, tt.Args, "expectedPageSize")
				assert.Equal(t, expectedPageSize, res.PageSize)
				for _, data := range res.Data {
					containsMovement(t, createdMovements, data)
				}
			} else {
				assert.Error(t, err)
				assert.ErrorAs(t, err, &tt.ExpectedErr, "Wanted: %v\nGot: %v", tt.ExpectedErr, err)
			}
		})
	}
	firestore_utils.DeleteTestMovementList(t, firestoreClient, createdMovements)
	firestore_utils.DeleteTestSource(t, firestoreClient, createdSource.UID)
	firestore_utils.DeleteTestUser(t, firestoreClient, createdOwner.UID)
}

func Test_movementRepositoryImplementation_UpdateMovement(t *testing.T) {
	firestoreClient := test_utils.InitTestingFireStore(t)
	firestoreClientFail := test_utils.InitTestingFireStoreFail(t)

	createdOwner := firestore_utils.CreateTestUser(t, firestoreClient, "0")
	createdSource := firestore_utils.CreateTestSource(t, firestoreClient, createdOwner.UID)
	createdMovement := firestore_utils.CreateTestMovement(t, firestoreClient, createdOwner.UID, createdSource.UID)
	updatedMovement := models.Movement{
		UID:          createdMovement.UID,
		Name:         "Update Movement",
		OwnerId:      "0",
		MovementDate: time.Now().UTC(),
		CreatedAt:    time.Now().UTC(),
		UpdatedAt:    time.Now().UTC(),
	}

	tests := []test_utils.TestCase{
		{
			Name: "Success",
			Fields: test_utils.Fields{
				"firestoreClient": firestoreClient,
			},
			Args: test_utils.Args{
				"movement": updatedMovement,
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
			Name: "Cant find",
			Fields: test_utils.Fields{
				"firestoreClient": firestoreClient,
			},
			Args: test_utils.Args{
				"movement": models.Movement{UID: "100", Name: "Not found Movement", OwnerId: "0"},
			},
			WantErr:     true,
			ExpectedErr: error_utils.FirestoreNotFoundError{DocID: "100"},
			PreTest:     nil,
		},
		{
			Name: "Different owner",
			Fields: test_utils.Fields{
				"firestoreClient": firestoreClient,
			},
			Args: test_utils.Args{
				"movement": models.Movement{UID: createdMovement.UID, Name: "Not found Movement", OwnerId: "1"},
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
			testFirestoreClient := test_utils.GetFieldByNameAndType[*firestore.Client](t, tt.Fields, "firestoreClient")
			r := &movementRepositoryImplementation{
				client: testFirestoreClient,
			}
			movementTest := test_utils.GetArgByNameAndType[models.Movement](t, tt.Args, "movement")
			updatedRes, err := r.UpdateMovement(movementTest)
			if !tt.WantErr {
				assert.NoError(t, err)
				checkEqualMovement(t, movementTest, updatedRes)
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
				"Collection":   "movements",
				"id":           createdMovement.UID,
				"originalData": createdMovement,
			}
			test_utils.ClearFireStoreTest(firestoreClient, "Update", args)

		})
	}
	defer func() {
		firestore_utils.DeleteTestMovement(t, firestoreClient, createdMovement.UID)
		firestore_utils.DeleteTestSource(t, firestoreClient, createdSource.UID)
		firestore_utils.DeleteTestUser(t, firestoreClient, createdOwner.UID)
	}()
}

func Test_movementRepositoryImplementation_DeleteMovement(t *testing.T) {

	firestoreClient := test_utils.InitTestingFireStore(t)
	firestoreClientFail := test_utils.InitTestingFireStoreFail(t)

	createdOwner := firestore_utils.CreateTestUser(t, firestoreClient, "0")
	createdSource := firestore_utils.CreateTestSource(t, firestoreClient, createdOwner.UID)
	createdMovement := firestore_utils.CreateTestMovement(t, firestoreClient, createdOwner.UID, createdSource.UID)

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
				"userID": createdMovement.OwnerId,
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
			Name:   "Different owner",
			Fields: testFields,
			Args: test_utils.Args{
				"id":     createdMovement.UID,
				"userID": "1000",
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
			r := &movementRepositoryImplementation{
				client: test_utils.GetFieldByNameAndType[*firestore.Client](t, tt.Fields, "firestoreClient"),
			}
			sourceTestId := test_utils.GetArgByNameAndType[string](t, tt.Args, "id")
			userTestID := test_utils.GetArgByNameAndType[string](t, tt.Args, "userID")

			err := r.DeleteMovement(sourceTestId, userTestID)
			if !tt.WantErr {
				assert.NoError(t, err)
			} else {
				assert.ErrorAs(t, err, &tt.ExpectedErr)
			}
		})
	}
	firestore_utils.DeleteTestSource(t, firestoreClient, createdSource.UID)
	firestore_utils.DeleteTestUser(t, firestoreClient, createdOwner.UID)
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
			MovementDate: time.Now().AddDate(0, 0, -i).UTC().Truncate(time.Hour * 24),
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

func containsMovement(t *testing.T, expectedList []models.Movement, got models.Movement) {
	for _, elem := range expectedList {
		if compareMovements(elem, got) {
			return
		}
	}
	bList, _ := json.MarshalIndent(expectedList, "", " ")
	dList, _ := json.MarshalIndent(got, "", " ")
	t.Fatalf("List\n%v\ndoes not contain\n%v", string(bList), string(dList))
}

func compareMovements(expected models.Movement, got models.Movement) bool {
	if expected.UID != got.UID {
		return false
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
	if !compareStringSlices(expected.Tags, got.Tags) {
		return false
	}

	delta := expected.CreatedAt.Sub(got.CreatedAt)

	if delta > time.Second*10 {
		return false
	}

	delta = expected.UpdatedAt.Sub(got.UpdatedAt)
	return delta <= time.Second*10
}

func compareStringSlices(slice1, slice2 []string) bool {
	if len(slice1) != len(slice2) {
		return false
	}
	for i := range slice1 {
		if slice1[i] != slice2[i] {
			return false
		}
	}
	return true
}
