package handlers

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"slices"
	"strings"
	"testing"
	"time"

	"github.com/AndresCRamos/midas-app-api/models"
	"github.com/AndresCRamos/midas-app-api/services"
	util_models "github.com/AndresCRamos/midas-app-api/utils/api/models"
	error_utils "github.com/AndresCRamos/midas-app-api/utils/errors"
	"github.com/AndresCRamos/midas-app-api/utils/test"
	test_utils "github.com/AndresCRamos/midas-app-api/utils/test"
	test_middleware "github.com/AndresCRamos/midas-app-api/utils/test/middleware"
	"github.com/AndresCRamos/midas-app-api/utils/test/mocks"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

var (
	movementValidationTests = []string{
		"No Name",
		"Empty Body",
		"Bad request",
	}
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

func Test_movementHandler_GetMovementsByUserAndDate(t *testing.T) {
	mockService := mocks.MovementServiceMock{}

	fields := test_utils.Fields{
		"mockService": mockService,
	}

	tests := []test_utils.TestCase{
		{
			Name:   "Success",
			Fields: fields,
			Args: test_utils.Args{
				"movementID":   "0",
				"expectedCode": http.StatusOK,
				"expectedMovement": util_models.PaginatedSearch[models.MovementRetrieve]{
					CurrentPage: 1,
					TotalData:   1,
					PageSize:    1,
					Data:        []models.MovementRetrieve{mocks.TestMovementRetrieve},
				},
				"userID": "0",
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
				"userID":       "1",
			},
			WantErr:     true,
			ExpectedErr: error_utils.APIUnknown{},
			PreTest:     nil,
		},
		{
			Name:   "Not enough data",
			Fields: fields,
			Args: test_utils.Args{
				"movementID":   "1",
				"expectedCode": http.StatusNotFound,
				"userID":       "2",
			},
			WantErr:     true,
			ExpectedErr: error_utils.MovementNotEnoughData{},
			PreTest:     nil,
		},
		{
			Name:   "Page bad type",
			Fields: fields,
			Args: test_utils.Args{
				"movementID":   "1",
				"userID":       "2",
				"expectedCode": http.StatusBadRequest,
				"page":         "page1",
			},
			WantErr:     true,
			ExpectedErr: util_models.PaginatedTypeError{},
			PreTest:     nil,
		},
		{
			Name:   "Bad date from type",
			Fields: fields,
			Args: test_utils.Args{
				"movementID":   "1",
				"userID":       "2",
				"expectedCode": http.StatusBadRequest,
				"date_from":    "bad_date_txt",
			},
			WantErr:     true,
			ExpectedErr: error_utils.APIBadDateFormat{DateString: "bad_date_txt", DateField: "date_from"},
			PreTest:     nil,
		},
		{
			Name:   "Bad date to type",
			Fields: fields,
			Args: test_utils.Args{
				"movementID":   "1",
				"userID":       "2",
				"expectedCode": http.StatusBadRequest,
				"date_to":      "bad_date_txt",
			},
			WantErr:     true,
			ExpectedErr: error_utils.APIBadDateFormat{DateString: "bad_date_txt", DateField: "date_to"},
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

			userID := test_utils.GetArgByNameAndType[string](t, tt.Args, "userID")
			page, _ := test_utils.ShouldGetArgByNameAndType[string](tt.Args, "page")
			dateTo, _ := test_utils.ShouldGetArgByNameAndType[string](tt.Args, "date_to")
			dateFrom, _ := test_utils.ShouldGetArgByNameAndType[string](tt.Args, "date_from")

			testRequest := test.TestRequest{
				Method:      http.MethodGet,
				BasePath:    "/",
				Handler:     h.GetMovementsByUserAndDate,
				Middlewares: []gin.HandlerFunc{test_middleware.TestMiddleware(userID)},
				QueryParams: map[string]string{
					"page":      page,
					"date_to":   dateTo,
					"date_from": dateFrom,
				},
			}

			w := testRequest.ServeRequest(t)

			expectedCode := test_utils.GetArgByNameAndType[int](t, tt.Args, "expectedCode")
			assert.Equal(t, expectedCode, w.Code)
			if !tt.WantErr {
				var resMovement util_models.PaginatedSearch[models.MovementRetrieve]
				testMovement := test_utils.GetArgByNameAndType[util_models.PaginatedSearch[models.MovementRetrieve]](t, tt.Args, "expectedMovement")
				err := json.Unmarshal(w.Body.Bytes(), &resMovement)
				assert.NoError(t, err)
				assert.NotEmpty(t, resMovement, "Error parsing response, got %v", w.Body.String())
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

func Test_movementHandler_CreateNewMovement(t *testing.T) {
	mockService := mocks.MovementServiceMock{}

	fields := test_utils.Fields{
		"mockService": mockService,
	}

	tests := []test_utils.TestCase{
		{
			Name:   "Success",
			Fields: fields,
			Args: test_utils.Args{
				"movement":     models.MovementCreate{Name: "Success", Amount: 100, SourceID: "123", MovementDate: time.Date(2024, time.January, 25, 0, 0, 0, 0, time.UTC)},
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
				"movement":     models.MovementCreate{Name: "CantConnect", Amount: 100, SourceID: "123", MovementDate: time.Date(2024, time.January, 25, 0, 0, 0, 0, time.UTC)},
				"expectedCode": http.StatusInternalServerError,
			},
			WantErr:     true,
			ExpectedErr: error_utils.APIUnknown{},
			PreTest:     nil,
		},
		{
			Name:   "Cant find source",
			Fields: fields,
			Args: test_utils.Args{
				"movement":          models.MovementCreate{Name: "NoSource", Amount: 100, SourceID: "123", MovementDate: time.Date(2024, time.January, 25, 0, 0, 0, 0, time.UTC)},
				"expectedCode":      http.StatusNotFound,
				"expectedErrDetail": []string{"The movement  cant be created, because source 123 doesn't exists"},
			},
			WantErr:     true,
			ExpectedErr: errors.New("The movement  cant be created, because source 123 doesn't exists"),
			PreTest:     nil,
		},
		{
			Name:   "No Name",
			Fields: fields,
			Args: test_utils.Args{
				"movement":          models.MovementCreate{Amount: 100, SourceID: "123", MovementDate: time.Date(2024, time.January, 25, 0, 0, 0, 0, time.UTC)},
				"expectedCode":      http.StatusBadRequest,
				"expectedErrDetail": []string{"field name is required"},
			},
			WantErr:     true,
			ExpectedErr: error_utils.APIInvalidRequestBody{},
			PreTest:     nil,
		},
		{
			Name:   "Empty Body",
			Fields: fields,
			Args: test_utils.Args{
				"movement":          models.MovementCreate{},
				"expectedCode":      http.StatusBadRequest,
				"expectedErrDetail": []string{"field name is required", "field sourceid is required", "field amount is required", "field movementdate is required"},
			},
			WantErr:     true,
			ExpectedErr: error_utils.APIInvalidRequestBody{},
			PreTest:     nil,
		},
		{
			Name:   "Bad request",
			Fields: fields,
			Args: test_utils.Args{
				"movement":          models.MovementCreate{},
				"expectedCode":      http.StatusBadRequest,
				"expectedErrDetail": []string{"invalid character 'I' looking for beginning of value"},
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
			mockService := test_utils.GetFieldByNameAndType[services.MovementService](t, tt.Fields, "mockService")
			h := &movementHandler{
				s: mockService,
			}

			body := getMovementTestBody[models.MovementCreate](t, tt)

			testRequest := test.TestRequest{
				Method:      http.MethodPost,
				BasePath:    "/",
				Handler:     h.CreateNewMovement,
				Middlewares: []gin.HandlerFunc{test_middleware.TestMiddleware("123")},
				Body:        bytes.NewBuffer(body),
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

				if slices.Contains(movementValidationTests, tt.Name) {
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

func Test_movementHandler_UpdateMovement(t *testing.T) {
	mockServiceMain := mocks.MovementServiceMock{}

	fields := test_utils.Fields{
		"mockService": mockServiceMain,
	}

	tests := []test_utils.TestCase{
		{
			Name:   "Success",
			Fields: fields,
			Args: test_utils.Args{
				"movement":     models.MovementUpdate{Name: "Success"},
				"id":           "0",
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
				"movement":     models.MovementUpdate{Name: "CantConnect"},
				"id":           "1",
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
				"movement":     models.MovementUpdate{Name: "NotFound"},
				"id":           "2",
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
				"movement":     models.MovementUpdate{Name: "NoOwner"},
				"id":           "4",
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

			id := test_utils.GetArgByNameAndType[string](t, tt.Args, "id")

			testRequest := test.TestRequest{
				Method:      http.MethodPut,
				BasePath:    "/:id",
				RequestPath: "/" + id,
				Middlewares: []gin.HandlerFunc{test_middleware.TestMiddleware("123")},
				Handler:     h.UpdateMovement,
				Body:        bytes.NewBuffer(getMovementTestBody[models.MovementUpdate](t, tt)),
			}
			w := testRequest.ServeRequest(t)

			expectedCode := test_utils.GetArgByNameAndType[int](t, tt.Args, "expectedCode")
			assert.Equal(t, expectedCode, w.Code)
			assert.NotEmpty(t, w.Body.String())
			if !tt.WantErr {
				var created models.MovementRetrieve
				err := json.Unmarshal(w.Body.Bytes(), &created)
				assert.NoError(t, err)

			} else {
				var errMessage map[string]interface{}
				err := json.Unmarshal(w.Body.Bytes(), &errMessage)
				assert.NoErrorf(t, err, "Response was: %v", w.Body.String())
				assert.Equal(t, tt.ExpectedErr.Error(), errMessage["error"])

				if slices.Contains(movementValidationTests, tt.Name) {
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

func getMovementTestBody[T any](test *testing.T, testCase test_utils.TestCase) []byte {
	testName := strings.Split(test.Name(), "/")[1]

	switch testName {
	case "Bad_request":
		return []byte("Invalid json")
	default:
		bodyStruct := test_utils.GetArgByNameAndType[T](test, testCase.Args, "movement")
		body, _ := json.Marshal(bodyStruct)
		return body
	}
}
