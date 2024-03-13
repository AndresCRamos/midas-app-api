package handlers

import (
	"bytes"
	"encoding/json"
	"net/http"
	"slices"
	"testing"

	"github.com/AndresCRamos/midas-app-api/models"
	"github.com/AndresCRamos/midas-app-api/services"
	error_utils "github.com/AndresCRamos/midas-app-api/utils/errors"
	test_utils "github.com/AndresCRamos/midas-app-api/utils/test"
	test_middleware "github.com/AndresCRamos/midas-app-api/utils/test/middleware"
	"github.com/AndresCRamos/midas-app-api/utils/test/mocks"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

var (
	userValidationTests = []string{"No UID", "No Name nor alias, No Name nor alias", "Lastname but no name"}
)

func Test_userHandler_CreateNewUser(t *testing.T) {
	mockService := mocks.UserServiceMock{}

	fields := test_utils.Fields{
		"mockService": mockService,
	}

	tests := []test_utils.TestCase{
		{
			Name:   "Success",
			Fields: fields,
			Args: test_utils.Args{
				"user":         models.User{Name: "Success", UID: "0"},
				"expectedCode": http.StatusCreated,
			},
			WantErr:     false,
			ExpectedErr: nil,
			PreTest:     nil,
		},
		{
			Name:   "Fail to connect",
			Fields: fields,
			Args: test_utils.Args{
				"user":         models.User{Name: "CantConnect", UID: "0"},
				"expectedCode": http.StatusInternalServerError,
			},
			WantErr:     true,
			ExpectedErr: error_utils.APIUnknown{},
			PreTest:     nil,
		},
		{
			Name:   "Duplicated User",
			Fields: fields,
			Args: test_utils.Args{
				"user":         models.User{Name: "Duplicated", UID: "0"},
				"expectedCode": http.StatusBadRequest,
			},
			WantErr:     true,
			ExpectedErr: error_utils.UserDuplicated{UserID: "0"},
			PreTest:     nil,
		},
		{
			Name:   "No Name nor alias",
			Fields: fields,
			Args: test_utils.Args{
				"user":              models.User{UID: "0"},
				"expectedCode":      http.StatusBadRequest,
				"expectedErrDetail": []string{"field alias is required if name is not supplied", "field name is required if alias is not supplied"},
			},
			WantErr:     true,
			ExpectedErr: error_utils.APIInvalidRequestBody{},
			PreTest:     nil,
		},
		{
			Name:   "Lastname but no name",
			Fields: fields,
			Args: test_utils.Args{
				"user":              models.User{UID: "0", LastName: "test_last", Alias: "alias_test"},
				"expectedCode":      http.StatusBadRequest,
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
			mockService := test_utils.GetFieldByNameAndType[services.UserService](t, tt.Fields, "mockService")
			h := &userHandler{
				s: mockService,
			}

			body := test_utils.GetTestBody[models.User](t, tt.Args, "user")

			testRequest := test_utils.TestRequest{
				Method:      http.MethodPost,
				BasePath:    "/",
				Handler:     h.CreateNewUser,
				Body:        bytes.NewBuffer(body),
				Middlewares: []gin.HandlerFunc{test_middleware.TestMiddleware("0")},
			}

			w := testRequest.ServeRequest(t)

			expectedCode := test_utils.GetArgByNameAndType[int](t, tt.Args, "expectedCode")
			assert.Equal(t, expectedCode, w.Code)
			if !tt.WantErr {
			} else {
				var errMessage map[string]interface{}
				err := json.Unmarshal(w.Body.Bytes(), &errMessage)
				assert.NoError(t, err)
				assert.Equal(t, tt.ExpectedErr.Error(), errMessage["error"])

				if slices.Contains(userValidationTests, tt.Name) {
					expectedDetail := test_utils.GetArgByNameAndType[[]string](t, tt.Args, "expectedErrDetail")
					val, ok := errMessage["detail"]
					if ok {
						assert.Equal(t, expectedDetail[0], val.(string))
						return
					}
					val, ok = errMessage["details"]
					if ok {
						assert.ElementsMatch(t, expectedDetail, val)
						return
					}

					errMsjByte, _ := json.MarshalIndent(errMessage, "", " ")

					t.Fatalf("Cant get error details from error message:\n %s", string(errMsjByte))
				}

			}
		})
	}
}

func Test_userHandler_GetUserByID(t *testing.T) {
	mockService := mocks.UserServiceMock{}

	fields := test_utils.Fields{
		"mockService": mockService,
	}

	tests := []test_utils.TestCase{
		{
			Name:   "Success",
			Fields: fields,
			Args: test_utils.Args{
				"userID":       "0",
				"expectedCode": http.StatusOK,
				"expectedUser": mocks.TestUser,
			},
			WantErr:     false,
			ExpectedErr: nil,
			PreTest:     nil,
		},
		{
			Name:   "Fail to connect",
			Fields: fields,
			Args: test_utils.Args{
				"userID":       "1",
				"expectedCode": http.StatusInternalServerError,
			},
			WantErr:     true,
			ExpectedErr: error_utils.APIUnknown{},
			PreTest:     nil,
		},
		{
			Name:   "Not Found",
			Fields: fields,
			Args: test_utils.Args{
				"userID":       "2",
				"expectedCode": http.StatusNotFound,
			},
			WantErr:     true,
			ExpectedErr: error_utils.UserNotFound{UserID: "2"},
			PreTest:     nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.Name, func(t *testing.T) {
			if tt.PreTest != nil {
				tt.PreTest(t)
			}
			mockService := test_utils.GetFieldByNameAndType[services.UserService](t, tt.Fields, "mockService")
			h := &userHandler{
				s: mockService,
			}

			userID := test_utils.GetArgByNameAndType[string](t, tt.Args, "userID")

			testRequest := test_utils.TestRequest{
				Method:      http.MethodGet,
				BasePath:    "/:id",
				RequestPath: "/" + userID,
				Handler:     h.GetUserByID,
				Middlewares: []gin.HandlerFunc{test_middleware.TestMiddleware(userID)},
			}

			w := testRequest.ServeRequest(t)

			expectedCode := test_utils.GetArgByNameAndType[int](t, tt.Args, "expectedCode")
			assert.Equal(t, expectedCode, w.Code)
			if !tt.WantErr {
				var resUser models.User
				testUser := test_utils.GetArgByNameAndType[models.User](t, tt.Args, "expectedUser")
				err := json.Unmarshal(w.Body.Bytes(), &resUser)
				assert.NoError(t, err)
				assert.Equal(t, testUser, resUser)
			} else {
				var errMessage map[string]interface{}
				err := json.Unmarshal(w.Body.Bytes(), &errMessage)
				assert.NoError(t, err)
				assert.Equal(t, tt.ExpectedErr.Error(), errMessage["error"])
			}
		})
	}
}
