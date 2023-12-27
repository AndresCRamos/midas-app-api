package test

import (
	"reflect"
	"testing"
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

func GetValueByNameAndType(t *testing.T, args Args, name string, targetType reflect.Type) interface{} {
	value, ok := args[name]
	if !ok {
		t.Fatalf("name '%s' not found in Args", name)
	}

	valueType := reflect.TypeOf(value)

	if valueType.AssignableTo(targetType) {
		return value
	}

	convertedValue := reflect.New(targetType).Elem()
	if !convertedValue.Type().ConvertibleTo(valueType) {
		t.Fatalf("unable to convert value for name '%s' to type '%s'", name, targetType)
	}

	convertedValue.Set(reflect.ValueOf(value).Convert(targetType))
	return convertedValue.Interface()
}
