package services

import (
	"encoding/json"
	"testing"
	"time"

	"github.com/AndresCRamos/midas-app-api/models"
	"github.com/AndresCRamos/midas-app-api/repository"
	util_models "github.com/AndresCRamos/midas-app-api/utils/api/models"
	error_utils "github.com/AndresCRamos/midas-app-api/utils/errors"
	test_utils "github.com/AndresCRamos/midas-app-api/utils/test"
	"github.com/AndresCRamos/midas-app-api/utils/test/mocks"
	"github.com/stretchr/testify/assert"
)

func Test_movementServiceImplementation_CreateNewMovement(t *testing.T) {

	mockRepo := mocks.MovementRepositoryMock{}

	fields := test_utils.Fields{
		"mockRepo": mockRepo,
	}

	tests := []test_utils.TestCase{
		{
			Name:   "Success",
			Fields: fields,
			Args: test_utils.Args{
				"movement": models.Movement{Name: "Success"},
			},
			WantErr:     false,
			ExpectedErr: nil,
			PreTest:     nil,
		},
		{
			Name:   "Fail to connect",
			Fields: fields,
			Args: test_utils.Args{
				"movement": models.Movement{Name: "CantConnect"},
			},
			WantErr:     true,
			ExpectedErr: error_utils.FirebaseUnknownError{},
			PreTest:     nil,
		},
		{
			Name:   "Duplicated Movement",
			Fields: fields,
			Args: test_utils.Args{
				"movement": models.Movement{Name: "Duplicated", UID: "0"},
			},
			WantErr:     true,
			ExpectedErr: error_utils.FirestoreAlreadyExistsError{},
			PreTest:     nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.Name, func(t *testing.T) {
			if tt.PreTest != nil {
				tt.PreTest(t)
			}
			mockRepo := test_utils.GetFieldByNameAndType[repository.MovementRepository](t, tt.Fields, "mockRepo")
			s := &movementServiceImplementation{
				r: mockRepo,
			}
			movementTest := test_utils.GetArgByNameAndType[models.Movement](t, tt.Args, "movement")
			_, err := s.CreateNewMovement(movementTest)
			if !tt.WantErr {
				assert.NoError(t, err)
			} else {
				assert.Error(t, err)
				assert.ErrorAs(t, err, &tt.ExpectedErr)
			}
		})

	}
}

func Test_movementServiceImplementation_GetMovementsByUserAndDate(t *testing.T) {

	mockRepo := mocks.MovementRepositoryMock{}

	fields := test_utils.Fields{
		"mockRepo": mockRepo,
	}

	tests := []test_utils.TestCase{
		{
			Name:   "Success",
			Fields: fields,
			Args: test_utils.Args{
				"userID":    "0",
				"page":      1,
				"date_from": time.Now().UTC().Add(-10 * time.Second),
				"date_to":   time.Now().UTC(),
				"expectedData": util_models.PaginatedSearch[models.Movement]{
					CurrentPage: 1,
					TotalData:   1,
					PageSize:    1,
					Data: []models.Movement{
						mocks.TestMovement,
					},
				},
			},
			WantErr:     false,
			ExpectedErr: nil,
			PreTest:     nil,
		},
		{
			Name:   "Fail to connect",
			Fields: fields,
			Args: test_utils.Args{
				"userID":    "1",
				"page":      1,
				"date_from": time.Now().UTC(),
				"date_to":   time.Now().UTC().Add(-10 * time.Second),
			},
			WantErr:     true,
			ExpectedErr: error_utils.FirebaseUnknownError{},
			PreTest:     nil,
		},
		{
			Name:   "Not enough data",
			Fields: fields,
			Args: test_utils.Args{
				"userID":    "2",
				"page":      1,
				"date_from": time.Now().UTC(),
				"date_to":   time.Now().UTC().Add(-10 * time.Second),
			},
			WantErr:     true,
			ExpectedErr: error_utils.MovementNotEnoughData{},
			PreTest:     nil,
		},
		{
			Name:   "Bad dates",
			Fields: fields,
			Args: test_utils.Args{
				"userID":    "1",
				"page":      1,
				"date_from": time.Now().UTC().Add(-10 * time.Second),
				"date_to":   time.Now().UTC(),
			},
			WantErr:     true,
			ExpectedErr: error_utils.MovementBadDates{},
			PreTest:     nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.Name, func(t *testing.T) {
			if tt.PreTest != nil {
				tt.PreTest(t)
			}
			mockRepo := test_utils.GetFieldByNameAndType[repository.MovementRepository](t, tt.Fields, "mockRepo")
			s := &movementServiceImplementation{
				r: mockRepo,
			}
			userId := test_utils.GetArgByNameAndType[string](t, tt.Args, "userID")
			page := test_utils.GetArgByNameAndType[int](t, tt.Args, "page")
			date_from := test_utils.GetArgByNameAndType[time.Time](t, tt.Args, "date_from")
			date_to := test_utils.GetArgByNameAndType[time.Time](t, tt.Args, "date_to")
			got, err := s.GetMovementsByUserAndDate(userId, page, date_from, date_to)
			if !tt.WantErr {
				expected := test_utils.GetArgByNameAndType[util_models.PaginatedSearch[models.Movement]](t, tt.Args, "expectedData")
				if !assert.NoError(t, err) {
					expectedByte, _ := json.MarshalIndent(expected, "", " ")
					gotByte, _ := json.MarshalIndent(expected, "", " ")
					assert.Equalf(t, expected.Data, got.Data, "Wanted:\n%s\nGot:\n%s\n", expectedByte, gotByte)
				}

			} else {
				if !assert.Errorf(t, err, "Wanted:\n%s\nGot:\n%s\n", tt.ExpectedErr, err) {
					assert.ErrorAs(t, err, &tt.ExpectedErr)
				}
			}
		})

	}
}

func Test_movementServiceImplementation_GetMovementByID(t *testing.T) {

	mockRepo := mocks.MovementRepositoryMock{}

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
			ExpectedErr: error_utils.FirebaseUnknownError{},
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
			ExpectedErr: error_utils.FirestoreParsingError{},
			PreTest:     nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.Name, func(t *testing.T) {
			if tt.PreTest != nil {
				tt.PreTest(t)
			}
			mockRepo := test_utils.GetFieldByNameAndType[repository.MovementRepository](t, tt.Fields, "mockRepo")
			s := &movementServiceImplementation{
				r: mockRepo,
			}
			id := test_utils.GetArgByNameAndType[string](t, tt.Args, "id")
			got, err := s.GetMovementByID(id, "123")
			if !tt.WantErr {
				assert.Equal(t, got, mocks.TestMovement)
				assert.NoError(t, err)
			} else {
				assert.Error(t, err)
				assert.ErrorAs(t, err, &tt.ExpectedErr)
			}
		})

	}
}

func Test_movementServiceImplementation_UpdateMovement(t *testing.T) {

	mockRepo := mocks.MovementRepositoryMock{}

	fields := test_utils.Fields{
		"mockRepo": mockRepo,
	}

	tests := []test_utils.TestCase{
		{
			Name:   "Success",
			Fields: fields,
			Args: test_utils.Args{
				"movement": models.Movement{Name: "Success", UID: "0"},
			},
			WantErr:     false,
			ExpectedErr: nil,
			PreTest:     nil,
		},
		{
			Name:   "Fail to connect",
			Fields: fields,
			Args: test_utils.Args{
				"movement": models.Movement{Name: "Cant update", UID: "1"},
			},
			WantErr:     true,
			ExpectedErr: error_utils.FirebaseUnknownError{},
			PreTest:     nil,
		},
		{
			Name:   "Cant find",
			Fields: fields,
			Args: test_utils.Args{
				"movement": models.Movement{Name: "Cant update", UID: "2"},
			},
			WantErr:     true,
			ExpectedErr: error_utils.FirestoreNotFoundError{},
			PreTest:     nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.Name, func(t *testing.T) {
			if tt.PreTest != nil {
				tt.PreTest(t)
			}
			mockRepo := test_utils.GetFieldByNameAndType[repository.MovementRepository](t, tt.Fields, "mockRepo")
			s := &movementServiceImplementation{
				r: mockRepo,
			}
			testMovement := test_utils.GetArgByNameAndType[models.Movement](t, tt.Args, "movement")
			_, err := s.UpdateMovement(testMovement)
			if !tt.WantErr {
				assert.NoError(t, err)
			} else {
				assert.Error(t, err)
				assert.ErrorAs(t, err, &tt.ExpectedErr)
			}
		})

	}
}

func Test_movementServiceImplementation_DeleteMovement(t *testing.T) {

	mockRepo := mocks.MovementRepositoryMock{}

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
			ExpectedErr: error_utils.FirebaseUnknownError{},
			PreTest:     nil,
		},
		{
			Name:   "Cant find",
			Fields: fields,
			Args: test_utils.Args{
				"id": "2",
			},
			WantErr:     true,
			ExpectedErr: error_utils.FirestoreNotFoundError{},
			PreTest:     nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.Name, func(t *testing.T) {
			if tt.PreTest != nil {
				tt.PreTest(t)
			}
			mockRepo := test_utils.GetFieldByNameAndType[repository.MovementRepository](t, tt.Fields, "mockRepo")
			s := &movementServiceImplementation{
				r: mockRepo,
			}
			deleteIDMovement := test_utils.GetArgByNameAndType[string](t, tt.Args, "id")
			err := s.DeleteMovement(deleteIDMovement, "123")
			if !tt.WantErr {
				assert.NoError(t, err)
			} else {
				assert.Error(t, err)
				assert.ErrorAs(t, err, &tt.ExpectedErr)
			}
		})

	}
}
