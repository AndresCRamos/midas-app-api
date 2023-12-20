package repository

import (
	"testing"

	"cloud.google.com/go/firestore"
	"github.com/AndresCRamos/midas-app-api/models"
	"github.com/AndresCRamos/midas-app-api/utils/errors"
	test_utils "github.com/AndresCRamos/midas-app-api/utils/test"
	"github.com/stretchr/testify/assert"
)

type testCase struct {
	name        string
	fields      fields
	args        args
	wantErr     bool
	expectedErr error
}

type fields struct {
	firestoreClient *firestore.Client
}

type args map[string]interface{}

func TestUserRepositoryImplementation_CreateNewUser(t *testing.T) {

	firestoreClient := test_utils.InitTestingFireStore(t)
	firestoreClientFail := test_utils.InitTestingFireStoreFail(t)

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
			errors.UNKNOWN,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
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
