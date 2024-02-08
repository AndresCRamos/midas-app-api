package test

import (
	"errors"
	"testing"

	error_utils "github.com/AndresCRamos/midas-app-api/utils/errors"
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

var (
	testString = "string"
	testInt    = 1
	testFloat  = 1.1
	testBool   = true
	testStruct = StructTest{testVal: "test"}
	testSlice  = []int{1, 2, 3}
	testMap    = map[string]interface{}{"name": "test"}
	testError  = errors.New("test")
)

func Test_getFromMap(t *testing.T) {
	testArgs := map[string]interface{}{
		"string":    testString,
		"int":       testInt,
		"float":     testFloat,
		"bool":      testBool,
		"struct":    testStruct,
		"slice":     testSlice,
		"map":       testMap,
		"interface": &testStruct,
		"error":     testError,
	}

	{
		// Test for string type
		t.Run("Type string", func(t *testing.T) {
			testGetFromMapString(t, testArgs)
		})

		// Test for int type
		t.Run("Type int", func(t *testing.T) {
			testGetFromMapInt(t, testArgs)
		})

		// Test for float type
		t.Run("Type float", func(t *testing.T) {
			testGetFromMapFloat(t, testArgs)
		})

		// Test for bool type
		t.Run("Type bool", func(t *testing.T) {
			testGetFromMapBool(t, testArgs)
		})

		// Test for struct type
		t.Run("Type struct", func(t *testing.T) {
			testGetFromMapStruct(t, testArgs)
		})

		// Test for slice type
		t.Run("Type slice", func(t *testing.T) {
			testGetFromMapSlice(t, testArgs)
		})

		// Test for map type
		t.Run("Type map", func(t *testing.T) {
			testGetFromMapMap(t, testArgs)
		})

		// Test for interface type
		t.Run("Type interface", func(t *testing.T) {
			testGetFromMapInterface(t, testArgs)
		})

		// Test for error type
		t.Run("Type error", func(t *testing.T) {
			testGetFromMapError(t, testArgs)
		})

		t.Run("Not found", func(t *testing.T) {
			testGetFromMapNotFound(t, testArgs)
		})

		t.Run("Bad type", func(t *testing.T) {
			testGetFromMapBadType(t, testArgs)
		})
	}
}

func testGetFromMapString(t *testing.T, testArgs map[string]interface{}) {
	got, err := getFromMap[string](testArgs, "string")
	assert.NoError(t, err)
	assert.Equal(t, testString, got)
}

func testGetFromMapInt(t *testing.T, testArgs map[string]interface{}) {
	got, err := getFromMap[int](testArgs, "int")
	assert.NoError(t, err)
	assert.Equal(t, testInt, got)
}

func testGetFromMapFloat(t *testing.T, testArgs map[string]interface{}) {
	got, err := getFromMap[float64](testArgs, "float")
	assert.NoError(t, err)
	assert.Equal(t, testFloat, got)
}

func testGetFromMapBool(t *testing.T, testArgs map[string]interface{}) {
	got, err := getFromMap[bool](testArgs, "bool")
	assert.NoError(t, err)
	assert.Equal(t, testBool, got)
}

func testGetFromMapStruct(t *testing.T, testArgs map[string]interface{}) {
	got, err := getFromMap[StructTest](testArgs, "struct")
	assert.NoError(t, err)
	assert.Equal(t, testStruct, got)
}

func testGetFromMapSlice(t *testing.T, testArgs map[string]interface{}) {
	got, err := getFromMap[[]int](testArgs, "slice")
	assert.NoError(t, err)
	assert.Equal(t, testSlice, got)
}

func testGetFromMapMap(t *testing.T, testArgs map[string]interface{}) {
	got, err := getFromMap[map[string]interface{}](testArgs, "map")
	assert.NoError(t, err)
	assert.Equal(t, testMap, got)
}

func testGetFromMapInterface(t *testing.T, testArgs map[string]interface{}) {
	got, err := getFromMap[TestInterface](testArgs, "interface")
	assert.NoError(t, err)
	assert.Equal(t, &testStruct, got)
}

func testGetFromMapError(t *testing.T, testArgs map[string]interface{}) {
	res, err := getFromMap[error](testArgs, "error")
	assert.Error(t, res)
	assert.NoError(t, err)
}

func testGetFromMapNotFound(t *testing.T, testArgs map[string]interface{}) {
	_, err := getFromMap[error](testArgs, "not_found")
	assert.Error(t, err)
	assert.Equal(t, error_utils.TestMapInterfaceNotFoundError{}, err)
}

func testGetFromMapBadType(t *testing.T, testArgs map[string]interface{}) {
	_, err := getFromMap[string](testArgs, "int")
	assert.Error(t, err)
	assert.Equal(t, error_utils.TestMapInterfaceCantAssertError{}, err)
}
