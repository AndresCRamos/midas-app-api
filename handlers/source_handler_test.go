package handlers

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"slices"
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
	sourceValidationTests = []string{
		"No Name",
	}
	mapNameID = map[string]string{
		"Success":         "0",
		"Fail to connect": "1",
		"Not Found":       "2",
		"Different user":  "4",
	}
)

func Test_sourceHandler_CreateNewSource(t *testing.T) {
	mockService := mocks.SourceServiceMock{}

	fields := test_utils.Fields{
		"mockService": mockService,
	}

	tests := []test_utils.TestCase{
		{
			Name:   "Success",
			Fields: fields,
			Args: test_utils.Args{
				"source":       models.SourceCreate{Name: "Success"},
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
				"source":       models.SourceCreate{Name: "CantConnect"},
				"expectedCode": http.StatusInternalServerError,
			},
			WantErr:     true,
			ExpectedErr: error_utils.APIUnknown{},
			PreTest:     nil,
		},
		{
			Name:   "No Name",
			Fields: fields,
			Args: test_utils.Args{
				"source":            models.SourceCreate{},
				"expectedCode":      http.StatusBadRequest,
				"expectedErrDetail": []string{"field name is required"},
			},
			WantErr:     true,
			ExpectedErr: error_utils.APIInvalidRequestBody{},
			PreTest:     nil,
		},
		{
			Name:   "Cant find owner",
			Fields: fields,
			Args: test_utils.Args{
				"source":            models.SourceCreate{Name: "NoOwner"},
				"expectedCode":      http.StatusNotFound,
				"expectedErrDetail": []string{"The source  cant be created, because owner 123 doesn't exists"},
			},
			WantErr:     true,
			ExpectedErr: errors.New("The source  cant be created, because owner 123 doesn't exists"),
			PreTest:     nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.Name, func(t *testing.T) {
			if tt.PreTest != nil {
				tt.PreTest(t)
			}
			mockService := test_utils.GetFieldByNameAndType[services.SourceService](t, tt.Fields, "mockService")
			h := &SourceHandler{
				s: mockService,
			}

			body := test_utils.GetTestBody[models.SourceCreate](t, tt.Args, "source")

			testRequest := test.TestRequest{
				Method:      http.MethodPost,
				BasePath:    "/",
				Handler:     h.CreateNewSource,
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

				if slices.Contains(sourceValidationTests, tt.Name) {
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

func Test_sourceHandler_GetSourcesByUser(t *testing.T) {
	mockService := mocks.SourceServiceMock{}

	fields := test_utils.Fields{
		"mockService": mockService,
	}

	tests := []test_utils.TestCase{
		{
			Name:   "Success",
			Fields: fields,
			Args: test_utils.Args{
				"sourceID":     "0",
				"expectedCode": http.StatusOK,
				"expectedSource": util_models.PaginatedSearch[models.SourceRetrieve]{
					CurrentPage: 1,
					TotalData:   1,
					PageSize:    1,
					Data:        []models.SourceRetrieve{mocks.TestSourceRetrieve},
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
				"sourceID":     "1",
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
				"sourceID":     "1",
				"expectedCode": http.StatusNotFound,
				"userID":       "2",
			},
			WantErr:     true,
			ExpectedErr: error_utils.SourceNotEnoughData{},
			PreTest:     nil,
		},
		{
			Name:   "Page bad type",
			Fields: fields,
			Args: test_utils.Args{
				"sourceID":     "1",
				"userID":       "2",
				"expectedCode": http.StatusBadRequest,
				"page":         "page1",
			},
			WantErr:     true,
			ExpectedErr: util_models.PaginatedTypeError{},
			PreTest:     nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.Name, func(t *testing.T) {
			if tt.PreTest != nil {
				tt.PreTest(t)
			}
			mockService := test_utils.GetFieldByNameAndType[services.SourceService](t, tt.Fields, "mockService")
			h := &SourceHandler{
				s: mockService,
			}

			userID := test_utils.GetArgByNameAndType[string](t, tt.Args, "userID")

			testRequest := test.TestRequest{
				Method:      http.MethodGet,
				BasePath:    "/",
				Handler:     h.GetSourcesByUser,
				Middlewares: []gin.HandlerFunc{test_middleware.TestMiddleware(userID)},
				QueryParams: map[string]string{},
			}
			page, err := test_utils.ShouldGetArgByNameAndType[string](tt.Args, "page")
			if err == nil {
				testRequest.QueryParams["page"] = page
			}

			w := testRequest.ServeRequest(t)

			expectedCode := test_utils.GetArgByNameAndType[int](t, tt.Args, "expectedCode")
			assert.Equal(t, expectedCode, w.Code)
			if !tt.WantErr {
				var resSource util_models.PaginatedSearch[models.SourceRetrieve]
				testSource := test_utils.GetArgByNameAndType[util_models.PaginatedSearch[models.SourceRetrieve]](t, tt.Args, "expectedSource")
				err := json.Unmarshal(w.Body.Bytes(), &resSource)
				assert.NoError(t, err)
				assert.NotEmpty(t, resSource, "Error parsing response, got %v", w.Body.String())
				assert.Equal(t, testSource, resSource)
			} else {
				var errMessage map[string]interface{}
				err := json.Unmarshal(w.Body.Bytes(), &errMessage)
				assert.NoError(t, err)
				assert.Equal(t, tt.ExpectedErr.Error(), errMessage["error"])
			}
		})
	}
}

func Test_sourceHandler_GetSourceByID(t *testing.T) {
	mockService := mocks.SourceServiceMock{}

	fields := test_utils.Fields{
		"mockService": mockService,
	}

	tests := []test_utils.TestCase{
		{
			Name:   "Success",
			Fields: fields,
			Args: test_utils.Args{
				"sourceID":       "0",
				"expectedCode":   http.StatusOK,
				"expectedSource": mocks.TestSource,
			},
			WantErr:     false,
			ExpectedErr: nil,
			PreTest:     nil,
		},
		{
			Name:   "Fail to connect",
			Fields: fields,
			Args: test_utils.Args{
				"sourceID":     "1",
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
				"sourceID":     "2",
				"expectedCode": http.StatusNotFound,
			},
			WantErr:     true,
			ExpectedErr: error_utils.SourceNotFound{SourceID: "2"},
			PreTest:     nil,
		},
		{
			Name:   "Different user",
			Fields: fields,
			Args: test_utils.Args{
				"sourceID":     "2",
				"expectedCode": http.StatusNotFound,
			},
			WantErr:     true,
			ExpectedErr: error_utils.SourceNotFound{SourceID: "2"},
			PreTest:     nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.Name, func(t *testing.T) {
			if tt.PreTest != nil {
				tt.PreTest(t)
			}
			mockService := test_utils.GetFieldByNameAndType[services.SourceService](t, tt.Fields, "mockService")
			h := &SourceHandler{
				s: mockService,
			}

			sourceID := test_utils.GetArgByNameAndType[string](t, tt.Args, "sourceID")

			testRequest := test.TestRequest{
				Method:      http.MethodGet,
				BasePath:    "/:id",
				RequestPath: "/" + sourceID,
				Handler:     h.GetSourceByID,
				Middlewares: []gin.HandlerFunc{test_middleware.TestMiddleware("123")},
			}

			w := testRequest.ServeRequest(t)
			expectedCode := test_utils.GetArgByNameAndType[int](t, tt.Args, "expectedCode")
			assert.Equal(t, expectedCode, w.Code)
			if !tt.WantErr {
				var resSource models.Source
				testSource := test_utils.GetArgByNameAndType[models.Source](t, tt.Args, "expectedSource")
				err := json.Unmarshal(w.Body.Bytes(), &resSource)
				assert.NoError(t, err)
				assert.Equal(t, testSource, resSource)
			} else {
				var errMessage map[string]interface{}
				err := json.Unmarshal(w.Body.Bytes(), &errMessage)
				assert.NoError(t, err)
				assert.Equal(t, tt.ExpectedErr.Error(), errMessage["error"])
			}
		})
	}
}

func Test_sourceHandler_GetMovementsBySourceAndDate(t *testing.T) {
	mockService := mocks.SourceServiceMock{}

	fields := test_utils.Fields{
		"mockService": mockService,
	}

	tests := []test_utils.TestCase{
		{
			Name:   "Success",
			Fields: fields,
			Args: test_utils.Args{
				"sourceID":       "0",
				"expectedCode":   http.StatusOK,
				"expectedSource": mocks.TestSource,
			},
			WantErr:     false,
			ExpectedErr: nil,
			PreTest:     nil,
		},
		{
			Name:   "Fail to connect",
			Fields: fields,
			Args: test_utils.Args{
				"sourceID":     "1",
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
				"sourceID":     "2",
				"expectedCode": http.StatusNotFound,
			},
			WantErr:     true,
			ExpectedErr: error_utils.SourceNotFound{SourceID: "2"},
			PreTest:     nil,
		},
		{
			Name:   "Different user",
			Fields: fields,
			Args: test_utils.Args{
				"sourceID":     "2",
				"expectedCode": http.StatusNotFound,
			},
			WantErr:     true,
			ExpectedErr: error_utils.SourceNotFound{SourceID: "2"},
			PreTest:     nil,
		},
		{
			Name:   "Bad dates",
			Fields: fields,
			Args: test_utils.Args{
				"sourceID":     "5",
				"expectedCode": http.StatusBadRequest,
			},
			WantErr:     true,
			ExpectedErr: error_utils.MovementBadDates{},
			PreTest:     nil,
		},
		{
			Name:   "Syntax date to error",
			Fields: fields,
			Args: test_utils.Args{
				"sourceID":     "5",
				"expectedCode": http.StatusBadRequest,
				"date_to":      "bad_date",
			},
			WantErr:     true,
			ExpectedErr: error_utils.APIBadDateFormat{DateString: "bad_date", DateField: "date_to", Format: time.DateOnly},
			PreTest:     nil,
		},
		{
			Name:   "Syntax date from error",
			Fields: fields,
			Args: test_utils.Args{
				"sourceID":     "5",
				"expectedCode": http.StatusBadRequest,
				"date_from":    "bad_date",
			},
			WantErr:     true,
			ExpectedErr: error_utils.APIBadDateFormat{DateString: "bad_date", DateField: "date_from", Format: time.DateOnly},
			PreTest:     nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.Name, func(t *testing.T) {
			if tt.PreTest != nil {
				tt.PreTest(t)
			}
			// mockService := test_utils.GetFieldByNameAndType[services.SourceService](t, tt.Fields, "mockService")
			h := &SourceHandler{
				s: mockService,
			}

			sourceID := test_utils.GetArgByNameAndType[string](t, tt.Args, "sourceID")
			page, _ := test_utils.ShouldGetArgByNameAndType[string](tt.Args, "page")
			date_from, _ := test_utils.ShouldGetArgByNameAndType[string](tt.Args, "date_from")
			date_to, _ := test_utils.ShouldGetArgByNameAndType[string](tt.Args, "date_to")

			testRequest := test.TestRequest{
				Method:      http.MethodGet,
				BasePath:    "/:id",
				RequestPath: "/" + sourceID,
				Handler:     h.GetMovementsBySourceAndDate,
				Middlewares: []gin.HandlerFunc{test_middleware.TestMiddleware("123")},
				QueryParams: map[string]string{
					"page":      page,
					"date_from": date_from,
					"date_to":   date_to,
				},
			}

			w := testRequest.ServeRequest(t)
			expectedCode := test_utils.GetArgByNameAndType[int](t, tt.Args, "expectedCode")
			assert.Equal(t, expectedCode, w.Code)
			if !tt.WantErr {
				var resSource util_models.PaginatedSearch[models.Source]
				testSource := test_utils.GetArgByNameAndType[models.Source](t, tt.Args, "expectedSource")
				err := json.Unmarshal(w.Body.Bytes(), &resSource)
				assert.NoError(t, err)
				containsSource(t, resSource.Data, testSource)

			} else {
				var errMessage map[string]interface{}
				err := json.Unmarshal(w.Body.Bytes(), &errMessage)
				assert.NoError(t, err)
				assert.Equal(t, tt.ExpectedErr.Error(), errMessage["error"])
			}
		})
	}
}

func Test_sourceHandler_UpdateSource(t *testing.T) {
	mockServiceMain := mocks.SourceServiceMock{}

	fields := test_utils.Fields{
		"mockService": mockServiceMain,
	}

	tests := []test_utils.TestCase{
		{
			Name:   "Success",
			Fields: fields,
			Args: test_utils.Args{
				"source":       models.SourceUpdate{Name: "Success"},
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
				"source":       models.SourceUpdate{Name: "CantConnect"},
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
				"source":       models.SourceUpdate{Name: "NotFound"},
				"expectedCode": http.StatusNotFound,
			},
			WantErr:     true,
			ExpectedErr: error_utils.SourceNotFound{SourceID: "2"},
			PreTest:     nil,
		},
		{
			Name:   "Different user",
			Fields: fields,
			Args: test_utils.Args{
				"source":       models.SourceUpdate{Name: "NoOwner"},
				"expectedCode": http.StatusNotFound,
			},
			WantErr:     true,
			ExpectedErr: error_utils.SourceDifferentOwner{SourceID: "4", OwnerID: "123"},
			PreTest:     nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.Name, func(t *testing.T) {
			if tt.PreTest != nil {
				tt.PreTest(t)
			}
			mockService := test_utils.GetFieldByNameAndType[services.SourceService](t, tt.Fields, "mockService")
			h := &SourceHandler{
				s: mockService,
			}

			body := test_utils.GetTestBody[models.SourceUpdate](t, tt.Args, "source")

			testRequest := test.TestRequest{
				Method:      http.MethodPut,
				BasePath:    "/:id",
				RequestPath: "/" + mapNameID[tt.Name],
				Middlewares: []gin.HandlerFunc{test_middleware.TestMiddleware("123")},
				Handler:     h.UpdateSource,
				Body:        bytes.NewBuffer(body),
			}
			w := testRequest.ServeRequest(t)

			expectedCode := test_utils.GetArgByNameAndType[int](t, tt.Args, "expectedCode")
			assert.Equal(t, expectedCode, w.Code)
			assert.NotEmpty(t, w.Body.String())
			if !tt.WantErr {
				var created models.SourceRetrieve
				err := json.Unmarshal(w.Body.Bytes(), &created)
				assert.NoError(t, err)

			} else {
				var errMessage map[string]interface{}
				err := json.Unmarshal(w.Body.Bytes(), &errMessage)
				assert.NoErrorf(t, err, "Response was: %v", w.Body.String())
				assert.Equal(t, tt.ExpectedErr.Error(), errMessage["error"])

				if slices.Contains(sourceValidationTests, tt.Name) {
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

func Test_sourceHandler_DeleteSource(t *testing.T) {
	mockService := mocks.SourceServiceMock{}

	fields := test_utils.Fields{
		"mockService": mockService,
	}

	tests := []test_utils.TestCase{
		{
			Name:   "Success",
			Fields: fields,
			Args: test_utils.Args{
				"sourceID":       "0",
				"expectedCode":   http.StatusOK,
				"expectedSource": mocks.TestSource,
			},
			WantErr:     false,
			ExpectedErr: nil,
			PreTest:     nil,
		},
		{
			Name:   "Fail to connect",
			Fields: fields,
			Args: test_utils.Args{
				"sourceID":     "1",
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
				"sourceID":     "2",
				"expectedCode": http.StatusNotFound,
			},
			WantErr:     true,
			ExpectedErr: error_utils.SourceNotFound{SourceID: "2"},
			PreTest:     nil,
		},
		{
			Name:   "Different user",
			Fields: fields,
			Args: test_utils.Args{
				"sourceID":     "4",
				"expectedCode": http.StatusNotFound,
			},
			WantErr:     true,
			ExpectedErr: error_utils.SourceDifferentOwner{SourceID: "4", OwnerID: "123"},
			PreTest:     nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.Name, func(t *testing.T) {
			if tt.PreTest != nil {
				tt.PreTest(t)
			}
			mockService := test_utils.GetFieldByNameAndType[services.SourceService](t, tt.Fields, "mockService")
			h := &SourceHandler{
				s: mockService,
			}

			sourceID := test_utils.GetArgByNameAndType[string](t, tt.Args, "sourceID")

			testRequest := test.TestRequest{
				Method:      http.MethodDelete,
				BasePath:    "/:id",
				RequestPath: "/" + sourceID,
				Handler:     h.DeleteSource,
				Middlewares: []gin.HandlerFunc{test_middleware.TestMiddleware("123")},
			}

			w := testRequest.ServeRequest(t)

			expectedCode := test_utils.GetArgByNameAndType[int](t, tt.Args, "expectedCode")
			assert.Equal(t, expectedCode, w.Code)
			if !tt.WantErr {
				assert.Empty(t, w.Body.String())

			} else {
				var errMessage map[string]interface{}
				err := json.Unmarshal(w.Body.Bytes(), &errMessage)
				assert.NoError(t, err)
				assert.Equal(t, tt.ExpectedErr.Error(), errMessage["error"])
			}
		})
	}
}

// func getSourceTestBody[T any](test *testing.T, testCase test_utils.TestCase) []byte {
// 	testName := strings.Split(test.Name(), "/")[1]

// 	switch testName {
// 	case "Bad_request":
// 		body, _ := json.Marshal(map[string]any{
// 			"InvalidBody": "Invalid",
// 		})
// 		return body
// 	default:
// 		bodyStruct := test_utils.GetArgByNameAndType[T](test, testCase.Args, "source")
// 		body, _ := json.Marshal(bodyStruct)
// 		return body
// 	}
// }

func containsSource(t *testing.T, expectedList []models.Source, got models.Source) {
	for _, elem := range expectedList {
		if compareSources(elem, got) {
			return
		}
	}
	bList, _ := json.MarshalIndent(expectedList, "", " ")
	dList, _ := json.MarshalIndent(got, "", " ")
	t.Fatalf("List\n%v\ndoes not contain\n%v", string(bList), string(dList))
}

func compareSources(expected models.Source, got models.Source) bool {
	if expected.UID != got.UID {
		return false
	}
	if expected.OwnerId != got.OwnerId {
		return false
	}
	if expected.Description != got.Description {
		return false
	}
	if expected.Name != got.Name {
		return false
	}
	if expected.Description != got.Description {
		return false
	}

	delta := expected.CreatedAt.Sub(got.CreatedAt)

	if delta > time.Second*10 {
		return false
	}

	delta = expected.UpdatedAt.Sub(got.UpdatedAt)
	return delta <= time.Second*10
}
