package services

import (
	"testing"

	"github.com/AndresCRamos/midas-app-api/models"
	"github.com/AndresCRamos/midas-app-api/repository"
	error_utils "github.com/AndresCRamos/midas-app-api/utils/errors"
	test_utils "github.com/AndresCRamos/midas-app-api/utils/test"
	"github.com/AndresCRamos/midas-app-api/utils/test/mocks"
	"github.com/stretchr/testify/assert"
)

func Test_userServiceImplementation_CreateNewUser(t *testing.T) {

	mockRepo := mocks.UserRepositoryMock{}

	fields := test_utils.Fields{
		"mockRepo": mockRepo,
	}

	tests := []test_utils.TestCase{
		{
			Name:   "Success",
			Fields: fields,
			Args: test_utils.Args{
				"user": models.User{Name: "Success"},
			},
			WantErr:     false,
			ExpectedErr: nil,
			PreTest:     nil,
		},
		{
			Name:   "Fail to connect",
			Fields: fields,
			Args: test_utils.Args{
				"user": models.User{Name: "CantConnect"},
			},
			WantErr:     true,
			ExpectedErr: error_utils.UNKNOWN,
			PreTest:     nil,
		},
		{
			Name:   "Duplicated User",
			Fields: fields,
			Args: test_utils.Args{
				"user": models.User{Name: "Duplicated", UID: "0"},
			},
			WantErr:     true,
			ExpectedErr: error_utils.AlreadyExistsError{},
			PreTest:     nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.Name, func(t *testing.T) {
			if tt.PreTest != nil {
				tt.PreTest(t)
			}
			s := &userServiceImplementation{
				r: tt.Fields["mockRepo"].(repository.UserRepository),
			}
			userTest := tt.Args["user"].(models.User)
			err := s.CreateNewUser(userTest)
			if !tt.WantErr {
				assert.NoError(t, err)
			} else {
				assert.Error(t, err)
				assert.ErrorAs(t, err, &tt.ExpectedErr)
			}
		})

	}
}

func Test_userServiceImplementation_GetUserByID(t *testing.T) {

	mockRepo := mocks.UserRepositoryMock{}

	fields := test_utils.Fields{
		"mockRepo": mockRepo,
	}

	tests := []test_utils.TestCase{
		{
			Name:   "Success",
			Fields: fields,
			Args: test_utils.Args{
				"id": "0",
			},
			WantErr:     false,
			ExpectedErr: nil,
			PreTest:     nil,
		},
		{
			Name:   "Fail to connect",
			Fields: fields,
			Args: test_utils.Args{
				"id": "1",
			},
			WantErr:     true,
			ExpectedErr: error_utils.UNKNOWN,
			PreTest:     nil,
		},
		{
			Name:   "Not Found",
			Fields: fields,
			Args: test_utils.Args{
				"id": "2",
			},
			WantErr:     true,
			ExpectedErr: error_utils.FirestoreNotFoundError{},
			PreTest:     nil,
		},
		{
			Name:   "Parsing error",
			Fields: fields,
			Args: test_utils.Args{
				"id": "3",
			},
			WantErr:     true,
			ExpectedErr: error_utils.ParsingError{},
			PreTest:     nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.Name, func(t *testing.T) {
			if tt.PreTest != nil {
				tt.PreTest(t)
			}
			s := &userServiceImplementation{
				r: tt.Fields["mockRepo"].(repository.UserRepository),
			}
			id := tt.Args["id"].(string)
			got, err := s.GetUserByID(id)
			if !tt.WantErr {
				assert.Equal(t, got, mocks.TestUser)
				assert.NoError(t, err)
			} else {
				assert.Error(t, err)
				assert.ErrorAs(t, err, &tt.ExpectedErr)
			}
		})

	}
}
