package handlers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"reflect"
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
				"user": models.User{Name: "Success"},
			},
			WantErr:     false,
			ExpectedErr: nil,
			PreTest:     nil,
		},
		{
			Name:   "Fail to connect",
			Fields: fields,
			Args: test_utils.Args{
				"user":         models.User{Name: "CantConnect"},
				"expectedCode": http.StatusInternalServerError,
			},
			WantErr:     true,
			ExpectedErr: error_utils.UNKNOWN,
			PreTest:     nil,
		},
		{
			Name:   "Duplicated User",
			Fields: fields,
			Args: test_utils.Args{
				"user":               models.User{Name: "Duplicated", UID: "0"},
				"expectedCode":       http.StatusBadRequest,
				"expectedErrMessage": "Document 0 already exists",
			},
			WantErr:     true,
			ExpectedErr: fmt.Errorf(error_utils.ALREADY_EXISTS, "1"),
			PreTest:     nil,
		},
		{
			Name:   "Bad request",
			Fields: fields,
			Args: test_utils.Args{
				"user":               models.User{Name: "Bad request"},
				"expectedCode":       http.StatusBadRequest,
				"expectedErrMessage": "EOF",
			},
			WantErr:     true,
			ExpectedErr: nil,
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
			h := &userHandler{
				s: tt.Fields["mockService"].(services.UserService),
			}

			bodyStruct := tt.Args["user"].(models.User)

			if bodyStruct.Name != "Bad request" {
				body, _ = json.Marshal(&bodyStruct)
			}

			testRouter.POST("/", h.CreateNewUser)
			req, _ := http.NewRequest("POST", "/", bytes.NewBuffer(body))
			testRouter.ServeHTTP(w, req)
			if !tt.WantErr {
				assert.Equal(t, http.StatusCreated, w.Code)
			} else {
				var errMessage error_utils.APIError
				assert.Equal(t, http.StatusBadRequest, w.Code)
				err := json.Unmarshal(w.Body.Bytes(), &errMessage)
				assert.NoError(t, err)
				expectedMsg := test_utils.GetValueByNameAndType(t, tt.Args, "expectedErrMessage", reflect.TypeOf(""))
				assert.Equal(t, expectedMsg, errMessage.Error)
			}
		})

	}
}
