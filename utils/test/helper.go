package test

import (
	"errors"
	"fmt"
	"reflect"
	"testing"

	error_utils "github.com/AndresCRamos/midas-app-api/utils/errors"
)

type TestCase struct {
	Name        string
	Fields      Fields
	Args        Args
	WantErr     bool
	ExpectedErr error
	PreTest     PreTestFunc
}

type PreTestFunc func(t *testing.T)

type Fields map[string]interface{}

type Args map[string]interface{}

type ExpectedValues map[string]interface{}

func ShouldGetArgByNameAndType(args Args, name string, targetType any) (any, error) {
	value, err := getFromMap(args, name, targetType)
	if err != nil {
		if errors.Is(err, error_utils.TestMapInterfaceNotFoundError{}) {
			return nil, fmt.Errorf("Cant find arg %s", name)
		}
		if errors.Is(err, error_utils.TestMapInterfaceCantAssertError{}) {
			return nil, fmt.Errorf("Cant assert %s arg to type %T", name, targetType)
		}
	}
	return value, nil
}

func GetArgByNameAndType(t *testing.T, args Args, name string, targetType any) any {
	value, err := ShouldGetArgByNameAndType(args, name, targetType)
	if err != nil {
		t.Fatal(err.Error())
	}
	return value
}

func ShouldGetFieldByNameAndType(fields Fields, name string, targetType any) (any, error) {
	value, err := getFromMap(fields, name, targetType)
	if err != nil {
		if errors.Is(err, error_utils.TestMapInterfaceNotFoundError{}) {
			return nil, fmt.Errorf("Cant find field %s", name)
		}
		if errors.Is(err, error_utils.TestMapInterfaceCantAssertError{}) {
			return nil, fmt.Errorf("Cant assert %s field to type %T", name, targetType)
		}
	}
	return value, nil
}

func GetFieldByNameAndType(t *testing.T, fields Fields, name string, targetType any) any {
	value, err := ShouldGetFieldByNameAndType(fields, name, targetType)
	if err != nil {
		t.Fatal(err.Error())
	}
	return value
}

func getFromMap(sourceMap map[string]interface{}, name string, targetType any) (any, error) {
	source, ok := sourceMap[name]
	if !ok {
		return nil, error_utils.TestMapInterfaceNotFoundError{}
	}

	sourceVal := reflect.ValueOf(source)
	targetTypeVal := reflect.TypeOf(targetType)

	// Check if targetType is an interface
	if targetTypeVal.Kind() == reflect.Ptr && targetTypeVal.Elem().Kind() == reflect.Interface {
		// If true, check if val map implements the interface
		if sourceVal.Type().Implements(targetTypeVal.Elem()) {
			return source, nil
		}
	} else {
		// Check if the types match directly
		if sourceVal.Type().AssignableTo(targetTypeVal) {
			return source, nil
		}
	}

	return nil, error_utils.TestMapInterfaceCantAssertError{}
}
