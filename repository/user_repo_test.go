package repository

import (
	"testing"

	"cloud.google.com/go/firestore"
	"github.com/AndresCRamos/midas-app-api/models"
	error_utils "github.com/AndresCRamos/midas-app-api/utils/errors"
	test_utils "github.com/AndresCRamos/midas-app-api/utils/test"
	"github.com/stretchr/testify/assert"
)

func TestUserRepositoryImplementation_CreateNewUser(t *testing.T) {

	firestoreClient := test_utils.InitTestingFireStore(t)
	firestoreClientFail := test_utils.InitTestingFireStoreFail(t)

	dupUser := models.User{
		UID:   "0",
		Alias: "DupUser",
	}

	createDupUser := func(t *testing.T) {
		rDuplicated := &UserRepositoryImplementation{
			client: firestoreClient,
		}

		err := rDuplicated.CreateNewUser(dupUser)
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
				"user": &models.User{UID: "0", Alias: "TestUser"},
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
				"user": &models.User{},
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
				"user": &dupUser,
			},
			WantErr:     true,
			ExpectedErr: error_utils.FirestoreAlreadyExistsError{DocID: dupUser.UID},
			PreTest:     createDupUser,
		},
	}

	for _, tt := range tests {
		t.Run(tt.Name, func(t *testing.T) {
			if tt.PreTest != nil {
				tt.PreTest(t)
			}
			testFirestoreClient := test_utils.GetFieldByNameAndType(t, tt.Fields, "firestoreClient", new(firestore.Client))
			r := &UserRepositoryImplementation{
				client: testFirestoreClient.(*firestore.Client),
			}
			// userTest := tt.Args["user"].(models.User)
			userTest := test_utils.GetArgByNameAndType(t, tt.Args, "user", new(models.User)).(*models.User)
			err := r.CreateNewUser(*userTest)
			if !tt.WantErr {
				assert.NoError(t, err)
			} else {
				assert.ErrorAs(t, err, &tt.ExpectedErr, "Expected error as: %s", tt.ExpectedErr.Error())
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

	searchUser := models.User{
		UID:   "1",
		Alias: "UserToSearch",
	}

	rSearch := UserRepositoryImplementation{
		client: firestoreClient,
	}

	testFields := test_utils.Fields{
		"firestoreClient": firestoreClient,
	}

	testFieldsFail := test_utils.Fields{
		"firestoreClient": firestoreClientFail,
	}

	err := rSearch.CreateNewUser(searchUser)
	if err != nil {
		t.Fatal("Cant create user to search")
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
				client: tt.Fields["firestoreClient"].(*firestore.Client),
			}
			userTestId := tt.Args["id"].(string)
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
