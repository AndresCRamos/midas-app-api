package handlers

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"slices"
	"strings"
	"testing"

	"github.com/AndresCRamos/midas-app-api/models"
	"github.com/AndresCRamos/midas-app-api/services"
	util_models "github.com/AndresCRamos/midas-app-api/utils/api/models"
	error_utils "github.com/AndresCRamos/midas-app-api/utils/errors"
	test_utils "github.com/AndresCRamos/midas-app-api/utils/test"
	"github.com/AndresCRamos/midas-app-api/utils/test/mocks"
	"github.com/AndresCRamos/midas-app-api/utils/validations"
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
			mockService := test_utils.GetFieldByNameAndType[services.SourceService](t, tt.Fields, "mockService")
			h := &sourceHandler{
				s: mockService,
			}

			body := getSourceTestBody[models.SourceCreate](t, tt)

			testRouter.Use(testMiddleware("123"))
			testRouter.POST("/", h.CreateNewSource)
			req, _ := http.NewRequest("POST", "/", bytes.NewBuffer(body))
			testRouter.ServeHTTP(w, req)
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
		gin.SetMode(gin.ReleaseMode)
		testRouter := gin.Default()
		t.Run(tt.Name, func(t *testing.T) {
			w := httptest.NewRecorder()
			if tt.PreTest != nil {
				tt.PreTest(t)
			}
			mockService := test_utils.GetFieldByNameAndType[services.SourceService](t, tt.Fields, "mockService")
			h := &sourceHandler{
				s: mockService,
			}

			userID := test_utils.GetArgByNameAndType[string](t, tt.Args, "userID")

			testRouter.Use(testMiddleware(userID))
			testRouter.GET("/", h.GetSourcesByUser)

			req, _ := http.NewRequest("GET", "/", bytes.NewBuffer([]byte{}))

			page, err := test_utils.ShouldGetArgByNameAndType[string](tt.Args, "page")
			if err == nil {
				query := req.URL.Query()
				query.Add("page", page)
				req.URL.RawQuery = query.Encode()
			}
			testRouter.ServeHTTP(w, req)
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
		gin.SetMode(gin.ReleaseMode)
		testRouter := gin.Default()
		w := httptest.NewRecorder()
		t.Run(tt.Name, func(t *testing.T) {
			var body []byte
			if tt.PreTest != nil {
				tt.PreTest(t)
			}
			mockService := test_utils.GetFieldByNameAndType[services.SourceService](t, tt.Fields, "mockService")
			h := &sourceHandler{
				s: mockService,
			}

			sourceID := test_utils.GetArgByNameAndType[string](t, tt.Args, "sourceID")
			url := fmt.Sprintf("/%s", sourceID)

			testRouter.Use(testMiddleware("123"))
			testRouter.GET("/:id", h.GetSourceByID)
			req, _ := http.NewRequest("GET", url, bytes.NewBuffer(body))
			testRouter.ServeHTTP(w, req)
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
			mockService := test_utils.GetFieldByNameAndType[services.SourceService](t, tt.Fields, "mockService")
			h := &sourceHandler{
				s: mockService,
			}

			body := getSourceTestBody[models.SourceUpdate](t, tt)

			testRouter.Use(testMiddleware("123"))
			testRouter.PUT("/:id", h.UpdateSource)

			id := mapNameID[tt.Name]
			req, err := http.NewRequest("PUT", "/"+id, bytes.NewBuffer(body))
			assert.NoError(t, err)
			testRouter.ServeHTTP(w, req)
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
		gin.SetMode(gin.ReleaseMode)
		testRouter := gin.Default()
		w := httptest.NewRecorder()
		t.Run(tt.Name, func(t *testing.T) {
			var body []byte
			if tt.PreTest != nil {
				tt.PreTest(t)
			}
			mockService := test_utils.GetFieldByNameAndType[services.SourceService](t, tt.Fields, "mockService")
			h := &sourceHandler{
				s: mockService,
			}

			sourceID := test_utils.GetArgByNameAndType[string](t, tt.Args, "sourceID")
			url := fmt.Sprintf("/%s", sourceID)

			testRouter.Use(testMiddleware("123"))
			testRouter.DELETE("/:id", h.DeleteSource)
			req, _ := http.NewRequest("DELETE", url, bytes.NewBuffer(body))
			testRouter.ServeHTTP(w, req)
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

func getSourceTestBody[T any](test *testing.T, testCase test_utils.TestCase) []byte {
	testName := strings.Split(test.Name(), "/")[1]

	switch testName {
	case "Bad_request":
		body, _ := json.Marshal(map[string]any{
			"InvalidBody": "Invalid",
		})
		return body
	default:
		bodyStruct := test_utils.GetArgByNameAndType[T](test, testCase.Args, "source")
		body, _ := json.Marshal(bodyStruct)
		return body
	}
}

func testMiddleware(id string) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Set("user", id)
	}
}
