package repository

import (
	"testing"

	"cloud.google.com/go/firestore"
	"github.com/AndresCRamos/midas-app-api/models"
	error_utils "github.com/AndresCRamos/midas-app-api/utils/errors"
	test_utils "github.com/AndresCRamos/midas-app-api/utils/test"
	firestore_utils "github.com/AndresCRamos/midas-app-api/utils/test/firestore"
	"github.com/stretchr/testify/assert"
)

func TestUserRepositoryImplementation_CreateNewUser(t *testing.T) {

	firestoreClient := test_utils.InitTestingFireStore(t)
	firestoreClientFail := test_utils.InitTestingFireStoreFail(t)

	createDupUser := func(t *testing.T) {
		firestore_utils.CreateTestUser(t, firestoreClient, "0")
	}

	tests := []test_utils.TestCase{
		{
			Name: "Success",
			Fields: test_utils.Fields{
				"firestoreClient": firestoreClient,
			},
			Args: test_utils.Args{
				"user": models.User{UID: "0", Alias: "TestUser"},
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
				"user": models.User{},
			},
			WantErr:     true,
			ExpectedErr: error_utils.FirebaseUnknownError{},
			PreTest:     nil,
		},
		{
			Name: "Duplicated user",
			Fields: test_utils.Fields{
				"firestoreClient": firestoreClient,
			},
			Args: test_utils.Args{
				"user": firestore_utils.SetTestUserID("0"),
			},
			WantErr:     true,
			ExpectedErr: error_utils.FirestoreAlreadyExistsError{DocID: "0"},
			PreTest:     createDupUser,
		},
	}

	for _, tt := range tests {
		t.Run(tt.Name, func(t *testing.T) {
			if tt.PreTest != nil {
				tt.PreTest(t)
			}
			testFirestoreClient := test_utils.GetFieldByNameAndType[*firestore.Client](t, tt.Fields, "firestoreClient")
			r := &UserRepositoryImplementation{
				client: testFirestoreClient,
			}
			userTest := test_utils.GetArgByNameAndType[models.User](t, tt.Args, "user")
			err := r.CreateNewUser(userTest)
			if !tt.WantErr {
				assert.NoError(t, err)
			} else {
				if assert.Error(t, err) {
					assert.ErrorAs(t, err, &tt.ExpectedErr, "Expected: %s\nGot: %s", tt.ExpectedErr.Error(), err.Error())
				}
			}
			args := map[string]interface{}{
				"Collection": "users",
				"id":         userTest.UID,
			}
			test_utils.ClearFireStoreTest(firestoreClient, "Create", args)
		})

	}
}

func TestUserRepositoryImplementation_GetUserByID(t *testing.T) {

	firestoreClient := test_utils.InitTestingFireStore(t)
	firestoreClientFail := test_utils.InitTestingFireStoreFail(t)

	searchUser := firestore_utils.CreateTestUser(t, firestoreClient, "1")

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
				"id": searchUser.UID,
			},
			WantErr:     false,
			ExpectedErr: nil,
			PreTest:     nil,
		},
		{
			Name:   "Fail to connect",
			Fields: testFieldsFail,
			Args: test_utils.Args{
				"id": searchUser.UID,
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
			r := &UserRepositoryImplementation{
				client: test_utils.GetFieldByNameAndType[*firestore.Client](t, tt.Fields, "firestoreClient"),
			}
			userTestId := test_utils.GetArgByNameAndType[string](t, tt.Args, "id")

			res, err := r.GetUserByID(userTestId)
			if !tt.WantErr {
				assert.NoError(t, err)
				assert.Equal(t, searchUser, res)
			} else {
				assert.ErrorAs(t, err, &tt.ExpectedErr)
			}
		})
	}
	args := test_utils.Args{
		"Collection": "users",
		"id":         searchUser.UID,
	}
	defer test_utils.ClearFireStoreTest(firestoreClient, "Create", args)
}
