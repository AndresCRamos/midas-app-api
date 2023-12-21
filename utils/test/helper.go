package test

import (
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
