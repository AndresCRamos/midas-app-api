package e2e

import (
	"bytes"
	"net/http"
	"testing"

	firestore "cloud.google.com/go/firestore"
	"github.com/AndresCRamos/midas-app-api/handlers"
	"github.com/AndresCRamos/midas-app-api/models"
	"github.com/AndresCRamos/midas-app-api/repository"
	"github.com/AndresCRamos/midas-app-api/services"
	"github.com/AndresCRamos/midas-app-api/utils/firebase"
	test_utils "github.com/AndresCRamos/midas-app-api/utils/test"
	firestore_utils "github.com/AndresCRamos/midas-app-api/utils/test/firestore"
	test_middleware "github.com/AndresCRamos/midas-app-api/utils/test/middleware"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func initUserTest(t *testing.T) (*firestore.Client, *handlers.UserHandler) {
	client, err := firebase.GetFireStoreClient()
	if err != nil {
		t.Fatalf("Cant initialize firestore client: %s", err)
	}

	repo := repository.NewUserRepository(client)
	service := services.NewUserService(repo)
	handler := handlers.NewUserHandler(service)

	return client, handler
}

func Test_user_CreateNewUser(t *testing.T) {

	firestoreClient, userHandler := initUserTest(t)

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
	}

	for _, tt := range tests {
		t.Run(tt.Name, func(t *testing.T) {
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
			}
			defer firestore_utils.DeleteTestUser(t, firestoreClient, "0")

		})
	}
}

func Test_user_GetUserByID(t *testing.T) {

}
