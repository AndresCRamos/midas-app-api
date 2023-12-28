package test

import (
	"errors"
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

func GetArgByNameAndType(t *testing.T, args Args, name string, targetType any) any {
	value, err := getFromMap(args, name, targetType)
	if err != nil {
		if errors.Is(err, error_utils.MAP_INTERFACE_NOT_FOUND) {
			t.Fatalf("Cant find arg %s", name)
		}
		if errors.Is(err, error_utils.MAP_INTERFACE_CANT_ASSERT) {
			t.Fatalf("Cant assert %s arg to type %T", name, targetType)
		}
	}
	return value
}

func GetFieldByNameAndType(t *testing.T, fields Fields, name string, targetType any) any {
	value, err := getFromMap(fields, name, targetType)
	if err != nil {
		if errors.Is(err, error_utils.MAP_INTERFACE_NOT_FOUND) {
			t.Fatalf("Cant find field %s", name)
		}
		if errors.Is(err, error_utils.MAP_INTERFACE_CANT_ASSERT) {
			t.Fatalf("Cant assert %s field to type %T", name, targetType)
		}
	}
	return value
}

func getFromMap(sourceMap map[string]interface{}, name string, targetType any) (any, error) {
	value, ok := sourceMap[name]
	if !ok {
		return nil, error_utils.MAP_INTERFACE_NOT_FOUND
	}

	switch v := value.(type) {
	case nil:
		return nil, nil
	default:
		if reflect.TypeOf(v).AssignableTo(reflect.TypeOf(targetType)) {
			return v, nil
		}
	}

	return nil, error_utils.MAP_INTERFACE_CANT_ASSERT

}
