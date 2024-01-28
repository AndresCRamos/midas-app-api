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
	sourceValidationTests = []string{"No UID"}
	mapNameID             = map[string]string{
		"Success":         "0",
		"Fail to connect": "1",
		"Not Found":       "2",
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
				"source":       &models.SourceCreate{Name: "Success", OwnerId: "0"},
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
				"source":       &models.SourceCreate{Name: "CantConnect", OwnerId: "0"},
				"expectedCode": http.StatusInternalServerError,
			},
			WantErr:     true,
			ExpectedErr: error_utils.APIUnknown{},
			PreTest:     nil,
		},
		{
			Name:   "Duplicated Source",
			Fields: fields,
			Args: test_utils.Args{
				"source":       &models.SourceCreate{Name: "Duplicated", OwnerId: "0"},
				"expectedCode": http.StatusBadRequest,
			},
			WantErr:     true,
			ExpectedErr: error_utils.SourceDuplicated{SourceID: "0"},
			PreTest:     nil,
		},
		{
			Name:   "No UID",
			Fields: fields,
			Args: test_utils.Args{
				"source":            &models.SourceCreate{Name: "Bad request", OwnerId: "0"},
				"expectedCode":      http.StatusBadRequest,
				"expectedErrDetail": []string{"field uid is required"},
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
			mockService := test_utils.GetFieldByNameAndType(t, tt.Fields, "mockService", new(services.SourceService))
			h := &sourceHandler{
				s: mockService.(services.SourceService),
			}

			body := getSourceTestBody[models.SourceCreate](t, tt)

			testRouter.POST("/", h.CreateNewSource)
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

				if slices.Contains(sourceValidationTests, tt.Name) {
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
			mockService := test_utils.GetFieldByNameAndType(t, tt.Fields, "mockService", new(services.SourceService))
			h := &sourceHandler{
				s: mockService.(services.SourceService),
			}

			sourceID := test_utils.GetArgByNameAndType(t, tt.Args, "sourceID", "").(string)
			url := fmt.Sprintf("/%s", sourceID)

			testRouter.GET("/:id", h.GetSourceByID)
			req, _ := http.NewRequest("GET", url, bytes.NewBuffer(body))
			testRouter.ServeHTTP(w, req)
			expectedCode := test_utils.GetArgByNameAndType(t, tt.Args, "expectedCode", 0)
			assert.Equal(t, expectedCode, w.Code)
			if !tt.WantErr {
				var resSource models.Source
				testSource := test_utils.GetArgByNameAndType(t, tt.Args, "expectedSource", models.Source{})
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
				"source":       &models.Source{Name: "Success", UID: "0", OwnerId: "0"},
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
				"source":       &models.Source{Name: "CantConnect", UID: "1", OwnerId: "0"},
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
				"source":       &models.Source{Name: "Duplicated", UID: "2", OwnerId: "0"},
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
		err := validations.AddCustomValidations()
		if err != nil {
			t.Fatal(err)
		}
		w := httptest.NewRecorder()
		t.Run(tt.Name, func(t *testing.T) {
			if tt.PreTest != nil {
				tt.PreTest(t)
			}
			mockService := test_utils.GetFieldByNameAndType(t, tt.Fields, "mockService", new(services.SourceService))
			h := &sourceHandler{
				s: mockService.(services.SourceService),
			}

			body := getSourceTestBody[models.Source](t, tt)

			testRouter.PUT("/:id", h.UpdateSource)
			id := mapNameID[tt.Name]
			req, err := http.NewRequest("PUT", "/"+id, bytes.NewBuffer(body))
			assert.NoError(t, err)
			testRouter.ServeHTTP(w, req)
			expectedCode := test_utils.GetArgByNameAndType(t, tt.Args, "expectedCode", 0)
			assert.Equal(t, expectedCode, w.Code)
			if !tt.WantErr {
				assert.Empty(t, w.Body.String())
			} else {
				var errMessage map[string]interface{}
				err := json.Unmarshal(w.Body.Bytes(), &errMessage)
				assert.NoError(t, err)
				assert.Equal(t, tt.ExpectedErr.Error(), errMessage["error"])

				if slices.Contains(sourceValidationTests, tt.Name) {
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
			mockService := test_utils.GetFieldByNameAndType(t, tt.Fields, "mockService", new(services.SourceService))
			h := &sourceHandler{
				s: mockService.(services.SourceService),
			}

			sourceID := test_utils.GetArgByNameAndType(t, tt.Args, "sourceID", "").(string)
			url := fmt.Sprintf("/%s", sourceID)

			testRouter.DELETE("/:id", h.DeleteSource)
			req, _ := http.NewRequest("DELETE", url, bytes.NewBuffer(body))
			testRouter.ServeHTTP(w, req)
			expectedCode := test_utils.GetArgByNameAndType(t, tt.Args, "expectedCode", 0)
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
			"InvalidUser": "Username",
		})
		return body
	default:
		bodyStruct := test_utils.GetArgByNameAndType(test, testCase.Args, "source", new(T)).(*T)
		body, _ := json.Marshal(bodyStruct)
		return body
	}
}
