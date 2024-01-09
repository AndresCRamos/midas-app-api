package errors

import "fmt"

const invalid_test_case = "Parameter %v is not a valid test case"

// TestInvalidTestCaseError struct
type TestInvalidTestCaseError struct {
	Param interface{}
}

func (tce TestInvalidTestCaseError) Error() string {
	return fmt.Sprintf(invalid_test_case, tce.Param)
}

func (tce TestInvalidTestCaseError) Unwrap() error {
	return nil
}

// Cant find value
type TestMapInterfaceNotFoundError struct{}

func (mnf TestMapInterfaceNotFoundError) Error() string {
	return "Cant find value"
}

func (mnf TestMapInterfaceNotFoundError) Unwrap() error {
	return nil
}

// Cant assert value
type TestMapInterfaceCantAssertError struct{}

func (mca TestMapInterfaceCantAssertError) Error() string {
	return "Cant assert value"
}

func (mca TestMapInterfaceCantAssertError) Unwrap() error {
	return nil
}
