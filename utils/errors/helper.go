package errors

import "fmt"

type ArgNotFoundError struct {
	Name string
}

func (ane ArgNotFoundError) Error() string {
	return fmt.Sprintf("Cant find arg %s", ane.Name)
}

type ArgTypeAssertionError[T any] struct {
	Name  string
	Value T
}

func (ate ArgTypeAssertionError[t]) Error() string {
	return fmt.Sprintf("Cant assert %s arg to type %T", ate.Name, ate.Value)
}

type FieldNotFoundError struct {
	Name string
}

func (fne FieldNotFoundError) Error() string {
	return fmt.Sprintf("Cant find field %s", fne.Name)
}

type FieldTypeAssertionError[T any] struct {
	Name  string
	Value T
}

func (fte FieldTypeAssertionError[T]) Error() string {
	return fmt.Sprintf("Cant assert %s field to type %T", fte.Name, fte.Value)
}
