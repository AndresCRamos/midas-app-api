package test

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

type TestInterface interface {
	FuncTest() string
}

type StructTest struct {
	testVal string
}

func (st *StructTest) FuncTest() string {
	return st.testVal
}

func InterfaceTest(t TestInterface) string {
	return t.FuncTest()
}

func Test_getFromMap(t *testing.T) {
	type args struct {
		sourceMap  map[string]interface{}
		name       string
		targetType interface{}
	}

	testArgs := map[string]interface{}{
		"string": "string",
		"int":    1,
		"float":  1.1,
		"bool":   true,
		"struct": struct {
			Name string
			Age  int
		}{
			Name: "test",
			Age:  1,
		},
		"slice": []int{1, 2, 3},
		"map": map[string]interface{}{
			"name": "test",
		},
		"interface": &StructTest{testVal: "test"},
		"error":     errors.New("test"),
	}
	tests := []struct {
		name    string
		args    args
		want    interface{}
		wantErr bool
	}{
		{
			name: "type string",
			args: args{
				sourceMap:  testArgs,
				name:       "string",
				targetType: "",
			},
			want: "string",
		},
		{
			name: "type int",
			args: args{
				sourceMap:  testArgs,
				name:       "int",
				targetType: 0,
			},
			want: 1,
		},
		{
			name: "type float",
			args: args{
				sourceMap:  testArgs,
				name:       "float",
				targetType: 0.0,
			},
			want: 1.1,
		},
		{
			name: "type bool",
			args: args{
				sourceMap:  testArgs,
				name:       "bool",
				targetType: false,
			},
			want: true,
		},
		{
			name: "type struct",
			args: args{
				sourceMap: testArgs,
				name:      "struct",
				targetType: struct {
					Name string
					Age  int
				}{},
			},
			want: struct {
				Name string
				Age  int
			}{
				Name: "test",
				Age:  1,
			},
		},
		{
			name: "type slice",
			args: args{
				sourceMap:  testArgs,
				name:       "slice",
				targetType: []int{},
			},
			want: []int{1, 2, 3},
		},
		{
			name: "type map",
			args: args{
				sourceMap:  testArgs,
				name:       "map",
				targetType: map[string]interface{}{},
			},
			want: map[string]interface{}{
				"name": "test",
			},
		},
		{
			name: "type interface",
			args: args{
				sourceMap:  testArgs,
				name:       "interface",
				targetType: new(TestInterface),
			},
			want: &StructTest{testVal: "test"},
		},
		{
			name: "type error",
			args: args{
				sourceMap:  testArgs,
				name:       "error",
				targetType: new(error),
			},
			wantErr: false,
			want:    errors.New("test"),
		},
		{
			name: "error",
			args: args{
				sourceMap:  testArgs,
				name:       "error",
				targetType: "",
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := getFromMap(tt.args.sourceMap, tt.args.name, tt.args.targetType)
			if tt.wantErr {
				assert.Error(t, err)
				return
			}
			assert.NoError(t, err)
			assert.Equal(t, tt.want, got)
		})
	}
}
