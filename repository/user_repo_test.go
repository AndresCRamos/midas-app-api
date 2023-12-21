package repository

import (
	"fmt"
	"testing"

	"cloud.google.com/go/firestore"
	"github.com/AndresCRamos/midas-app-api/models"
	error_utils "github.com/AndresCRamos/midas-app-api/utils/errors"
	test_utils "github.com/AndresCRamos/midas-app-api/utils/test"
	"github.com/stretchr/testify/assert"
)

type testCase struct {
	name        string
	fields      fields
	args        args
	wantErr     bool
	expectedErr error
	preTest     preTestFunc
}

type preTestFunc func(t *testing.T)

type fields struct {
	firestoreClient *firestore.Client
}

type args map[string]interface{}

func TestUserRepositoryImplementation_CreateNewUser(t *testing.T) {

	firestoreClient := test_utils.InitTestingFireStore(t)
	firestoreClientFail := test_utils.InitTestingFireStoreFail(t)

	dupUser := models.User{
		UID:   "0",
		Alias: "DupUser",
	}

	tests := []testCase{
		{
			"Success",
			fields{
				firestoreClient: firestoreClient,
			},
			args{
				"user": models.User{Name: "John", LastName: "Doe", UID: "0", Alias: "TestUser"},
			},
			false,
			nil,
			nil,
		},
		{
			"Fail to connect",
			fields{
				firestoreClient: firestoreClientFail,
			},
			args{
				"user": models.User{},
			},
			true,
			error_utils.UNKNOWN,
			func(t *testing.T) {},
		},
		{
			"Duplicated user",
			fields{
				firestoreClient: firestoreClient,
			},
			args{
				"user": dupUser,
			},
			true,
			fmt.Errorf(error_utils.ALREADY_EXISTS, dupUser.UID),
			func(t *testing.T) {
				rDuplicated := &UserRepositoryImplementation{
					client: firestoreClient,
				}

				err := rDuplicated.CreateNewUser(dupUser)
				if err != nil {
					t.Fatalf("Cant connect to Firestore to check for duplication test: %s", err.Error())
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.preTest != nil {
				tt.preTest(t)
			}
			r := &UserRepositoryImplementation{
				client: tt.fields.firestoreClient,
			}
			userTest := tt.args["user"].(models.User)
			err := r.CreateNewUser(userTest)
			if !tt.wantErr {
				assert.Nil(t, err)
			} else {
				assert.NotNil(t, err)
				assert.EqualError(t, err, tt.expectedErr.Error())
			}
			args := map[string]interface{}{
				"Collection": "users",
				"id":         userTest.UID,
			}
			test_utils.ClearFireStoreTest(firestoreClient, "Create", args)
		})

	}
}
