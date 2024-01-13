package handlers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"slices"
	"strings"
	"testing"

	"github.com/AndresCRamos/midas-app-api/models"
	"github.com/AndresCRamos/midas-app-api/services"
	error_utils "github.com/AndresCRamos/midas-app-api/utils/errors"
	test_utils "github.com/AndresCRamos/midas-app-api/utils/test"
	"github.com/AndresCRamos/midas-app-api/utils/test/mocks"
	"github.com/AndresCRamos/midas-app-api/utils/validations"
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
				"user":         &models.User{Name: "Success", UID: "0"},
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
				"user":         &models.User{Name: "CantConnect", UID: "0"},
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
				"user":         &models.User{Name: "Duplicated", UID: "0"},
				"expectedCode": http.StatusBadRequest,
			},
			WantErr:     true,
			ExpectedErr: error_utils.UserDuplicated{UserID: "0"},
			PreTest:     nil,
		},
		{
			Name:   "No UID",
			Fields: fields,
			Args: test_utils.Args{
				"user":              &models.User{Name: "Bad request"},
				"expectedCode":      http.StatusBadRequest,
				"expectedErrDetail": []string{"field uid is required"},
			},
			WantErr:     true,
			ExpectedErr: error_utils.APIInvalidRequestBody{},
			PreTest:     nil,
		},
		{
			Name:   "No Name nor alias",
			Fields: fields,
			Args: test_utils.Args{
				"user":              &models.User{UID: "0"},
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
				"user":              &models.User{UID: "0", LastName: "test_last", Alias: "alias_test"},
				"expectedCode":      http.StatusBadRequest,
				"expectedErrDetail": []string{"field lastname depends on name, which is not supplied"},
			},
			WantErr:     true,
			ExpectedErr: error_utils.APIInvalidRequestBody{},
			PreTest:     nil,
		},
	}

	for _, tt := range tests {
		gin.SetMode(gin.ReleaseMode)
		testRouter := gin.Default()
		err := validations.AddCustomValidations()
		if err != nil {
			t.Fatal(err)
		}
		w := httptest.NewRecorder()
		t.Run(tt.Name, func(t *testing.T) {
			if tt.PreTest != nil {
				tt.PreTest(t)
			}
			mockService := test_utils.GetFieldByNameAndType(t, tt.Fields, "mockService", new(services.UserService))
			h := &userHandler{
				s: mockService.(services.UserService),
			}

			body := getTestBody(t, tt)

			testRouter.POST("/", h.CreateNewUser)
			req, _ := http.NewRequest("POST", "/", bytes.NewBuffer(body))
			testRouter.ServeHTTP(w, req)
			expectedCode := test_utils.GetArgByNameAndType(t, tt.Args, "expectedCode", 0)
			assert.Equal(t, expectedCode, w.Code)
			if !tt.WantErr {
			} else {
				var errMessage map[string]interface{}
				err := json.Unmarshal(w.Body.Bytes(), &errMessage)
				assert.NoError(t, err)
				assert.Equal(t, tt.ExpectedErr.Error(), errMessage["error"])

				if slices.Contains(userValidationTests, tt.Name) {
					expectedDetail := test_utils.GetArgByNameAndType(t, tt.Args, "expectedErrDetail", []string{}).([]string)
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
		gin.SetMode(gin.ReleaseMode)
		testRouter := gin.Default()
		w := httptest.NewRecorder()
		t.Run(tt.Name, func(t *testing.T) {
			var body []byte
			if tt.PreTest != nil {
				tt.PreTest(t)
			}
			mockService := test_utils.GetFieldByNameAndType(t, tt.Fields, "mockService", new(services.UserService))
			h := &userHandler{
				s: mockService.(services.UserService),
			}

			userID := test_utils.GetArgByNameAndType(t, tt.Args, "userID", "").(string)
			url := fmt.Sprintf("/%s", userID)

			testRouter.GET("/:id", h.GetUserByID)
			req, _ := http.NewRequest("GET", url, bytes.NewBuffer(body))
			testRouter.ServeHTTP(w, req)
			expectedCode := test_utils.GetArgByNameAndType(t, tt.Args, "expectedCode", 0)
			assert.Equal(t, expectedCode, w.Code)
			if !tt.WantErr {
				var resUser models.User
				testUser := test_utils.GetArgByNameAndType(t, tt.Args, "expectedUser", models.User{})
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

func getTestBody(test *testing.T, testCase test_utils.TestCase) []byte {
	testName := strings.Split(test.Name(), "/")[1]

	switch testName {
	case "Bad_request":
		body, _ := json.Marshal(map[string]any{
			"InvalidUser": "Username",
		})
		return body
	default:
		bodyStruct := test_utils.GetArgByNameAndType(test, testCase.Args, "user", new(models.User)).(*models.User)
		body, _ := json.Marshal(bodyStruct)
		return body
	}
}
