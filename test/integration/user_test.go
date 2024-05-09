package integration

import (
	"bytes"
	"context"
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

	firestoreClient, firestoreClientFail, userHandler := initUserTest(t)

	createDupUser := func(t *testing.T) {
		firestore_utils.CreateTestUser(t, firestoreClient, "0")
	}

	field := test_utils.Fields{
		"userHandler": userHandler,
		"client":      firestoreClient,
	}

	failedField := test_utils.Fields{
		"userHandler": userHandler,
		"client":      firestoreClientFail,
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
		{
			Name:   "Fail to connect",
			Fields: failedField,
			Args: test_utils.Args{
				"user":         models.UserCreate{Alias: "TEST_USER", Name: "John", LastName: "Doe"},
				"expectedCode": http.StatusBadRequest,
			},
			WantErr:     true,
			ExpectedErr: error_utils.UserDuplicated{UserID: "0"},
			PreTest:     createDupUser,
		},
		{
			Name:   "No Name nor alias",
			Fields: field,
			Args: test_utils.Args{
				"user":              models.UserCreate{},
				"expectedCode":      http.StatusBadRequest,
				"isValidation":      true,
				"expectedErrDetail": []string{"field alias is required if name is not supplied", "field name is required if alias is not supplied"},
			},
			WantErr:     true,
			ExpectedErr: error_utils.APIInvalidRequestBody{},
			PreTest:     nil,
		},
		{
			Name:   "Lastname but no name",
			Fields: field,
			Args: test_utils.Args{
				"user":              models.UserCreate{LastName: "test_last", Alias: "alias_test"},
				"expectedCode":      http.StatusBadRequest,
				"isValidation":      true,
				"expectedErrDetail": []string{"field lastname depends on name, which is not supplied"},
			},
			WantErr:     true,
			ExpectedErr: error_utils.APIInvalidRequestBody{},
			PreTest:     nil,
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
				assert.Empty(t, w.Body.Bytes(), w.Body.String())
			} else {
				var errMessage map[string]interface{}
				err := json.Unmarshal(w.Body.Bytes(), &errMessage)
				assert.NoError(t, err)
				assert.Equal(t, tt.ExpectedErr.Error(), errMessage["error"])
				isValidation, err := test_utils.ShouldGetArgByNameAndType[bool](tt.Args, "isValidation")

				if err == nil && isValidation {
					expectErrDetail := test_utils.GetArgByNameAndType[[]string](t, tt.Args, "expectedErrDetail")

					if len(expectErrDetail) == 1 {
						assert.Equal(t, expectErrDetail[0], errMessage["detail"])
					} else {
						assert.ElementsMatch(t, expectErrDetail, errMessage["details"])
					}
				}
			}
			defer firestore_utils.DeleteTestUser(t, firestoreClient, "0")

		})
	}
}

func Test_user_GetUserByID(t *testing.T) {
	testUserID := "0"
	firestoreClient, _, userHandler := initUserTest(t)
	firestore_utils.CreateTestUser(t, firestoreClient, testUserID)

	createRandomData := func(t *testing.T) {
		_, err := firestoreClient.Collection("users").Doc("1").Set(context.Background(), map[string]any{
			"alias": 123,
		})
		if err != nil {
			t.Errorf("Cant create random data for parsing error: %v", err)
		}
	}

	deleteRandomData := func(t *testing.T) {
		_, err := firestoreClient.Collection("users").Doc("1").Delete(context.Background())
		if err != nil {
			t.Logf("Cant delete firestore random data: %v", err)
		}
	}

	field := test_utils.Fields{
		"userHandler": userHandler,
		"client":      firestoreClient,
	}

	tests := []test_utils.TestCase{
		{
			Name: "Success",
			Args: test_utils.Args{
				"expectedUser": firestore_utils.SetTestUserID(testUserID),
				"expectedCode": http.StatusOK,
				"userID":       testUserID,
			},
			Fields:  field,
			WantErr: false,
			PreTest: nil,
		},
		{
			Name: "Not Found",
			Args: test_utils.Args{
				"expectedCode": http.StatusNotFound,
				"userID":       "1000",
			},
			Fields:      field,
			WantErr:     true,
			ExpectedErr: error_utils.UserNotFound{UserID: "1000"},
			PreTest:     nil,
		},
		{
			Name: "Parsing Error",
			Args: test_utils.Args{
				"expectedCode": http.StatusInternalServerError,
				"userID":       "1",
			},
			Fields:      field,
			WantErr:     true,
			ExpectedErr: error_utils.APIUnknown{},
			PreTest:     createRandomData,
		},
	}

	for _, tt := range tests {
		t.Run(tt.Name, func(t *testing.T) {
			if tt.PreTest != nil {
				tt.PreTest(t)
			}
			handler := test_utils.GetFieldByNameAndType[*handlers.UserHandler](t, tt.Fields, "userHandler")
			userID := test_utils.GetArgByNameAndType[string](t, tt.Args, "userID")

			testRequest := test_utils.TestRequest{
				Method:      http.MethodGet,
				Handler:     handler.GetUserByID,
				BasePath:    "/:id",
				RequestPath: "/" + userID,
				Middlewares: []gin.HandlerFunc{test_middleware.TestMiddleware(userID)},
			}

			w := testRequest.ServeRequest(t)

			expectedCode := test_utils.GetArgByNameAndType[int](t, tt.Args, "expectedCode")
			assert.Equal(t, expectedCode, w.Code)

			if !tt.WantErr {
				expectedUser := test_utils.GetArgByNameAndType[models.User](t, tt.Args, "expectedUser")
				var resUser models.User
				err := json.Unmarshal(w.Body.Bytes(), &resUser)
				if assert.NoError(t, err) {
					assert.Equal(t, expectedUser, resUser)
				}
			} else {
				var errMessage map[string]interface{}
				err := json.Unmarshal(w.Body.Bytes(), &errMessage)
				assert.NoError(t, err)
				assert.Equal(t, tt.ExpectedErr.Error(), errMessage["error"])
			}
		})
	}
	defer func() {
		firestore_utils.DeleteTestUser(t, firestoreClient, testUserID)
		deleteRandomData(t)
	}()
}
