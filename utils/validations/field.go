package validations

import (
	"reflect"

	"github.com/go-playground/validator/v10"
)

func dependsOn(fl validator.FieldLevel) bool {
	// Get the current field's value and type
	currentFieldValue := fl.Field().Interface()

	// Get the name of the dependent field
	dependentFieldName := fl.Param()

	// Get the dependent field's value and type
	dependentField := fl.Parent().FieldByName(dependentFieldName).Interface()

	isCurrentPresent := isFieldPresent(currentFieldValue)
	isDependentPresent := isFieldPresent(dependentField)

	// Check if the dependent field is present, then the current field must also be present
	return !isCurrentPresent || (isDependentPresent && isCurrentPresent)
}

// isFieldPresent checks if the field is present (not equal to the zero value of its type)
func isFieldPresent(fieldValue interface{}) bool {
	return !reflect.DeepEqual(fieldValue, reflect.Zero(reflect.TypeOf(fieldValue)).Interface())
}
