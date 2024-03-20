package test

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	error_utils "github.com/AndresCRamos/midas-app-api/utils/errors"
	"github.com/AndresCRamos/midas-app-api/utils/validations"
	"github.com/gin-gonic/gin"
)

type TestRequest struct {
	Method      string
	BasePath    string
	RequestPath string
	Handler     gin.HandlerFunc
	Middlewares []gin.HandlerFunc
	Body        *bytes.Buffer
	Headers     map[string]string
	QueryParams map[string]string
}

func (te TestRequest) ServeRequest(t *testing.T) *httptest.ResponseRecorder {

	if te.BasePath == "" {
		t.Fatal("Base path must not be empty")
	}

	gin.SetMode(gin.ReleaseMode)
	testRouter := gin.New()
	err := validations.AddCustomValidations()
	if err != nil {
		t.Fatalf("Cant add custom validations to gin.Engine:\n%v", err)
	}

	if te.Middlewares != nil {
		testRouter.Use(te.Middlewares...)
	}

	testRouter.Handle(te.Method, te.BasePath, te.Handler)

	if te.RequestPath == "" {
		te.RequestPath = te.BasePath
	}

	if te.Body == nil {
		te.Body = bytes.NewBuffer([]byte{})
	}

	req, err := http.NewRequest(te.Method, te.RequestPath, te.Body)
	if err != nil {
		t.Fatalf("An error has ocurred creating the request:\n%v", err)
	}

	for headerKey, headerVal := range te.Headers {
		req.Header.Set(headerKey, headerVal)
	}

	if te.QueryParams != nil {
		q := req.URL.Query()
		for queryParamKey, queryParamVal := range te.QueryParams {
			if queryParamVal != "" {
				q.Add(queryParamKey, queryParamVal)
			}
		}
		req.URL.RawQuery = q.Encode()
	}

	w := httptest.NewRecorder()
	testRouter.ServeHTTP(w, req)

	return w
}

func getTestBodyGeneric[T any](t *testing.T, args Args, searchName string) []byte {

	body := GetArgByNameAndType[T](t, args, searchName)

	bodyBytes, err := json.Marshal(body)

	if err != nil {
		t.Fatalf("Cant parse body into json: %s", err)
	}

	return bodyBytes
}

func GetTestBody[T any](t *testing.T, args Args, searchName string) []byte {
	var bodyBytes []byte

	var isBytes bool
	isBytes, err := ShouldGetArgByNameAndType[bool](args, "isBytes")

	if err != nil {
		if errors.Is(err, error_utils.ArgNotFoundError{Name: "isBytes"}) {
			isBytes = false
		} else {
			t.Fatalf("Cant get isBytes value: %s", err)
		}
	}

	if isBytes {
		bodyBytes = GetArgByNameAndType[[]byte](t, args, searchName)
	} else {
		return getTestBodyGeneric[T](t, args, searchName)
	}

	return bodyBytes
}
