package handlers

import (
	"encoding/json"
	"net/http"
	"testing"

	"github.com/AndresCRamos/midas-app-api/models"
	"github.com/AndresCRamos/midas-app-api/services"
	error_utils "github.com/AndresCRamos/midas-app-api/utils/errors"
	"github.com/AndresCRamos/midas-app-api/utils/test"
	test_utils "github.com/AndresCRamos/midas-app-api/utils/test"
	test_middleware "github.com/AndresCRamos/midas-app-api/utils/test/middleware"
	"github.com/AndresCRamos/midas-app-api/utils/test/mocks"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func Test_movementHandler_GetMovementByID(t *testing.T) {
	mockService := mocks.MovementServiceMock{}

	fields := test_utils.Fields{
		"mockService": mockService,
	}

	tests := []test_utils.TestCase{
		{
			Name:   "Success",
			Fields: fields,
			Args: test_utils.Args{
				"movementID":       "0",
				"expectedCode":     http.StatusOK,
				"expectedMovement": mocks.TestMovement,
			},
			WantErr:     false,
			ExpectedErr: nil,
			PreTest:     nil,
		},
		{
			Name:   "Fail to connect",
			Fields: fields,
			Args: test_utils.Args{
				"movementID":   "1",
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
				"movementID":   "2",
				"expectedCode": http.StatusNotFound,
			},
			WantErr:     true,
			ExpectedErr: error_utils.MovementNotFound{MovementID: "2"},
			PreTest:     nil,
		},
		{
			Name:   "Different user",
			Fields: fields,
			Args: test_utils.Args{
				"movementID":   "4",
				"expectedCode": http.StatusNotFound,
			},
			WantErr:     true,
			ExpectedErr: error_utils.MovementDifferentOwner{MovementID: "4", OwnerID: "123"},
			PreTest:     nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.Name, func(t *testing.T) {
			if tt.PreTest != nil {
				tt.PreTest(t)
			}
			mockService := test_utils.GetFieldByNameAndType[services.MovementService](t, tt.Fields, "mockService")
			h := &movementHandler{
				s: mockService,
			}

			movementID := test_utils.GetArgByNameAndType[string](t, tt.Args, "movementID")

			testRequest := test.TestRequest{
				Method:      http.MethodGet,
				BasePath:    "/:id",
				RequestPath: "/" + movementID,
				Handler:     h.GetMovementByID,
				Middlewares: []gin.HandlerFunc{test_middleware.TestMiddleware("123")},
			}

			w := testRequest.ServeRequest(t)
			expectedCode := test_utils.GetArgByNameAndType[int](t, tt.Args, "expectedCode")
			assert.Equal(t, expectedCode, w.Code)
			if !tt.WantErr {
				var resMovement models.Movement
				testMovement := test_utils.GetArgByNameAndType[models.Movement](t, tt.Args, "expectedMovement")
				err := json.Unmarshal(w.Body.Bytes(), &resMovement)
				assert.NoError(t, err)
				assert.Equal(t, testMovement, resMovement)
			} else {
				var errMessage map[string]interface{}
				err := json.Unmarshal(w.Body.Bytes(), &errMessage)
				assert.NoError(t, err)
				assert.Equal(t, tt.ExpectedErr.Error(), errMessage["error"])
			}
		})
	}
}
