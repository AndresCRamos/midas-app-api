package e2e

import (
	"bytes"
	"net/http"
	"testing"

	"github.com/AndresCRamos/midas-app-api/handlers"
	"github.com/AndresCRamos/midas-app-api/models"
	"github.com/AndresCRamos/midas-app-api/repository"
	"github.com/AndresCRamos/midas-app-api/services"
	"github.com/AndresCRamos/midas-app-api/utils/firebase"
	test_utils "github.com/AndresCRamos/midas-app-api/utils/test"
	"github.com/stretchr/testify/assert"
)

func initUserTest(t *testing.T) *handlers.UserHandler {
	client, err := firebase.GetFireStoreClient()
	if err != nil {
		t.Fatalf("Cant initialize firestore client: %s", err)
	}

	repo := repository.NewUserRepository(client)
	service := services.NewUserService(repo)
	handler := handlers.NewUserHandler(service)

	return handler
}

func Test_user_CreateNewUser(t *testing.T) {

	userHandler := initUserTest(t)

	field := test_utils.Fields{
		"userHandler": userHandler,
	}

	tests := []test_utils.TestCase{
		{
			Name: "",
			Args: test_utils.Args{
				"user":         models.UserCreate{},
				"expectedCode": http.StatusBadRequest,
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
				Method:   http.MethodPost,
				Handler:  handler.CreateNewUser,
				BasePath: "/",
				Body:     bytes.NewBuffer(body),
			}

			w := testRequest.ServeRequest(t)

			expectedCode := test_utils.GetArgByNameAndType[int](t, tt.Args, "expectedCode")

			assert.Equal(t, expectedCode, w.Code)

		})
	}
}

func Test_user_GetUserByID(t *testing.T) {

}
