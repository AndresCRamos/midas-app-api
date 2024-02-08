package test

import (
	"errors"
	"fmt"
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

func ShouldGetArgByNameAndType[T any](args Args, name string) (T, error) {
	var value T
	value, err := getFromMap[T](args, name)
	if err != nil {
		if errors.Is(err, error_utils.TestMapInterfaceNotFoundError{}) {
			return value, fmt.Errorf("Cant find arg %s", name)
		}
		if errors.Is(err, error_utils.TestMapInterfaceCantAssertError{}) {
			return value, fmt.Errorf("Cant assert %s arg to type %T", name, value)
		}
	}

	return value, nil
}

func GetArgByNameAndType[T any](t *testing.T, args Args, name string, targetType any) T {
	value, err := ShouldGetArgByNameAndType[T](args, name)
	if err != nil {
		t.Fatal(err.Error())
	}
	return value
}

func ShouldGetFieldByNameAndType[T any](fields Fields, name string) (T, error) {
	var value T
	value, err := getFromMap[T](fields, name)
	if err != nil {
		if errors.Is(err, error_utils.TestMapInterfaceNotFoundError{}) {
			return value, fmt.Errorf("Cant find field %s", name)
		}
		if errors.Is(err, error_utils.TestMapInterfaceCantAssertError{}) {
			return value, fmt.Errorf("Cant assert %s field to type %T", name, value)
		}
	}
	return value, nil
}

func GetFieldByNameAndType[T any](t *testing.T, fields Fields, name string) any {
	value, err := ShouldGetFieldByNameAndType[T](fields, name)
	if err != nil {
		t.Fatal(err.Error())
	}
	return value
}

func getFromMap[T any](sourceMap map[string]interface{}, name string) (T, error) {
	var result T
	source, ok := sourceMap[name]
	if !ok {
		return result, error_utils.TestMapInterfaceNotFoundError{}
	}

	result, ok = source.(T)
	if ok {
		return result, nil
	}

	return result, error_utils.TestMapInterfaceCantAssertError{}
}
