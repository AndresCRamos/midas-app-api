package handlers

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/AndresCRamos/midas-app-api/models"
	"github.com/AndresCRamos/midas-app-api/services"
	error_utils "github.com/AndresCRamos/midas-app-api/utils/errors"
	test_utils "github.com/AndresCRamos/midas-app-api/utils/test"
	"github.com/AndresCRamos/midas-app-api/utils/test/mocks"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
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
				"user":         &models.User{Name: "Success"},
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
				"user":         &models.User{Name: "CantConnect"},
				"expectedCode": http.StatusInternalServerError,
			},
			WantErr:     true,
			ExpectedErr: error_utils.FirebaseUnknownError{},
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
			Name:   "Bad request",
			Fields: fields,
			Args: test_utils.Args{
				"user":         &models.User{Name: "Bad request"},
				"expectedCode": http.StatusBadRequest,
			},
			WantErr:     true,
			ExpectedErr: error_utils.FirebaseUnknownError{},
			PreTest:     nil,
		},
	}

	for _, tt := range tests {
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

			bodyStruct := test_utils.GetArgByNameAndType(t, tt.Args, "user", new(models.User)).(*models.User)

			if bodyStruct.Name == "Bad request" {
				body, _ = json.Marshal(map[string]any{
					"InvalidUser": "Username",
				})
			} else {
				body, _ = json.Marshal(bodyStruct)
			}

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
			}
		})
	}
}
