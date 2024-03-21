package e2e

import (
	"bytes"
	"encoding/json"
	"net/http"
	"testing"

	firestore "cloud.google.com/go/firestore"
	"github.com/AndresCRamos/midas-app-api/handlers"
	"github.com/AndresCRamos/midas-app-api/models"
	"github.com/AndresCRamos/midas-app-api/repository"
	"github.com/AndresCRamos/midas-app-api/services"
	error_utils "github.com/AndresCRamos/midas-app-api/utils/errors"
	"github.com/AndresCRamos/midas-app-api/utils/firebase"
	test_utils "github.com/AndresCRamos/midas-app-api/utils/test"
	firestore_utils "github.com/AndresCRamos/midas-app-api/utils/test/firestore"
	test_middleware "github.com/AndresCRamos/midas-app-api/utils/test/middleware"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func initUserTest(t *testing.T) (*firestore.Client, *firestore.Client, *handlers.UserHandler) {
	client, err := firebase.GetFireStoreClient()
	if err != nil {
		t.Fatalf("Cant initialize firestore client: %s", err)
	}

	repo := repository.NewUserRepository(client)
	service := services.NewUserService(repo)
	handler := handlers.NewUserHandler(service)

	failedClient := test_utils.InitTestingFireStoreFail(t)

	return client, failedClient, handler
}

func Test_user_CreateNewUser(t *testing.T) {

	firestoreClient, userHandler := initUserTest(t)

	createDupUser := func(t *testing.T) {
		firestore_utils.CreateTestUser(t, firestoreClient, "0")
	}

	field := test_utils.Fields{
		"userHandler": userHandler,
	}

	tests := []test_utils.TestCase{
		{
			Name: "Success",
			Args: test_utils.Args{
				"user":         models.UserCreate{Alias: "TEST_USER", Name: "John", LastName: "Doe"},
				"expectedUser": models.User{Alias: "TEST_USER", Name: "John", LastName: "Doe", UID: "0"},
				"expectedCode": http.StatusCreated,
			},
			Fields:  field,
			WantErr: false,
			PreTest: nil,
		},
		{
			Name:   "Duplicated user",
			Fields: field,
			Args: test_utils.Args{
				"user":         models.UserCreate{Alias: "TEST_USER", Name: "John", LastName: "Doe"},
				"expectedCode": http.StatusBadRequest,
			},
			WantErr:     true,
			ExpectedErr: error_utils.UserDuplicated{UserID: "0"},
			PreTest:     createDupUser,
		},
	}

	for _, tt := range tests {
		t.Run(tt.Name, func(t *testing.T) {
			if tt.PreTest != nil {
				tt.PreTest(t)
			}
			handler := test_utils.GetFieldByNameAndType[*handlers.UserHandler](t, tt.Fields, "userHandler")
			body := test_utils.GetTestBody[models.UserCreate](t, tt.Args, "user")

			testRequest := test_utils.TestRequest{
				Method:      http.MethodPost,
				Handler:     handler.CreateNewUser,
				BasePath:    "/",
				Body:        bytes.NewBuffer(body),
				Middlewares: []gin.HandlerFunc{test_middleware.TestMiddleware("0")},
			}

			w := testRequest.ServeRequest(t)

			expectedCode := test_utils.GetArgByNameAndType[int](t, tt.Args, "expectedCode")
			assert.Equal(t, expectedCode, w.Code)

			if !tt.WantErr {
				assert.Empty(t, w.Body.Bytes())
			} else {
				var errMessage map[string]interface{}
				err := json.Unmarshal(w.Body.Bytes(), &errMessage)
				assert.NoError(t, err)
				assert.Equal(t, tt.ExpectedErr.Error(), errMessage["error"])
			}
			defer firestore_utils.DeleteTestUser(t, firestoreClient, "0")

		})
	}
}

func Test_user_GetUserByID(t *testing.T) {

}
