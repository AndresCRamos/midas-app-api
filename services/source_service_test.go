package services

import (
	"testing"

	"github.com/AndresCRamos/midas-app-api/models"
	"github.com/AndresCRamos/midas-app-api/repository"
	util_models "github.com/AndresCRamos/midas-app-api/utils/api/models"
	error_utils "github.com/AndresCRamos/midas-app-api/utils/errors"
	test_utils "github.com/AndresCRamos/midas-app-api/utils/test"
	"github.com/AndresCRamos/midas-app-api/utils/test/mocks"
	"github.com/stretchr/testify/assert"
)

func Test_sourceServiceImplementation_CreateNewSource(t *testing.T) {

	mockRepo := mocks.SourceRepositoryMock{}

	fields := test_utils.Fields{
		"mockRepo": mockRepo,
	}

	tests := []test_utils.TestCase{
		{
			Name:   "Success",
			Fields: fields,
			Args: test_utils.Args{
				"source": &models.Source{Name: "Success"},
			},
			WantErr:     false,
			ExpectedErr: nil,
			PreTest:     nil,
		},
		{
			Name:   "Fail to connect",
			Fields: fields,
			Args: test_utils.Args{
				"source": &models.Source{Name: "CantConnect"},
			},
			WantErr:     true,
			ExpectedErr: error_utils.FirebaseUnknownError{},
			PreTest:     nil,
		},
		{
			Name:   "Duplicated Source",
			Fields: fields,
			Args: test_utils.Args{
				"source": &models.Source{Name: "Duplicated", UID: "0"},
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
			mockRepo := test_utils.GetFieldByNameAndType(t, tt.Fields, "mockRepo", new(repository.SourceRepository))
			s := &sourceServiceImplementation{
				r: mockRepo.(repository.SourceRepository),
			}
			sourceTest := test_utils.GetArgByNameAndType(t, tt.Args, "source", new(models.Source)).(*models.Source)
			_, err := s.CreateNewSource(*sourceTest)
			if !tt.WantErr {
				assert.NoError(t, err)
			} else {
				assert.Error(t, err)
				assert.ErrorAs(t, err, &tt.ExpectedErr)
			}
		})

	}
}

func Test_sourceServiceImplementation_GetSourcesByUser(t *testing.T) {

	mockRepo := mocks.SourceRepositoryMock{}

	fields := test_utils.Fields{
		"mockRepo": mockRepo,
	}

	tests := []test_utils.TestCase{
		{
			Name:   "Success",
			Fields: fields,
			Args: test_utils.Args{
				"userID": "0",
				"expectedData": util_models.PaginatedSearch[models.Source]{
					CurrentPage: 1,
					TotalData:   1,
					PageSize:    1,
					Data: []models.Source{
						mocks.TestSource,
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
				"userID": "1",
			},
			WantErr:     true,
			ExpectedErr: error_utils.FirebaseUnknownError{},
			PreTest:     nil,
		},
		{
			Name:   "Not enough data",
			Fields: fields,
			Args: test_utils.Args{
				"userID": "2",
			},
			WantErr:     true,
			ExpectedErr: error_utils.SourceNotEnoughData{},
			PreTest:     nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.Name, func(t *testing.T) {
			if tt.PreTest != nil {
				tt.PreTest(t)
			}
			mockRepo := test_utils.GetFieldByNameAndType(t, tt.Fields, "mockRepo", new(repository.SourceRepository))
			s := &sourceServiceImplementation{
				r: mockRepo.(repository.SourceRepository),
			}
			userId := test_utils.GetArgByNameAndType(t, tt.Args, "userID", "").(string)
			got, err := s.GetSourcesByUser(userId, 1)
			if !tt.WantErr {
				expected := test_utils.GetArgByNameAndType(t, tt.Args, "expectedData", util_models.PaginatedSearch[models.Source]{})
				assert.NoError(t, err)
				assert.Equal(t, expected, got)
			} else {
				assert.Error(t, err)
				assert.ErrorAs(t, err, &tt.ExpectedErr)
			}
		})

	}
}

func Test_sourceServiceImplementation_GetSourceByID(t *testing.T) {

	mockRepo := mocks.SourceRepositoryMock{}

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
			mockRepo := test_utils.GetFieldByNameAndType(t, tt.Fields, "mockRepo", new(repository.SourceRepository))
			s := &sourceServiceImplementation{
				r: mockRepo.(repository.SourceRepository),
			}
			id := test_utils.GetArgByNameAndType(t, tt.Args, "id", "").(string)
			got, err := s.GetSourceByID(id, "123")
			if !tt.WantErr {
				assert.Equal(t, got, mocks.TestSource)
				assert.NoError(t, err)
			} else {
				assert.Error(t, err)
				assert.ErrorAs(t, err, &tt.ExpectedErr)
			}
		})

	}
}

func Test_sourceServiceImplementation_UpdateSource(t *testing.T) {

	mockRepo := mocks.SourceRepositoryMock{}

	fields := test_utils.Fields{
		"mockRepo": mockRepo,
	}

	tests := []test_utils.TestCase{
		{
			Name:   "Success",
			Fields: fields,
			Args: test_utils.Args{
				"source": &models.Source{Name: "Success", UID: "0"},
			},
			WantErr:     false,
			ExpectedErr: nil,
			PreTest:     nil,
		},
		{
			Name:   "Fail to connect",
			Fields: fields,
			Args: test_utils.Args{
				"source": &models.Source{Name: "Cant update", UID: "1"},
			},
			WantErr:     true,
			ExpectedErr: error_utils.FirebaseUnknownError{},
			PreTest:     nil,
		},
		{
			Name:   "Cant find",
			Fields: fields,
			Args: test_utils.Args{
				"source": &models.Source{Name: "Cant update", UID: "2"},
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
			mockRepo := test_utils.GetFieldByNameAndType(t, tt.Fields, "mockRepo", new(repository.SourceRepository))
			s := &sourceServiceImplementation{
				r: mockRepo.(repository.SourceRepository),
			}
			testSource := test_utils.GetArgByNameAndType(t, tt.Args, "source", new(models.Source)).(*models.Source)
			_, err := s.UpdateSource(*testSource)
			if !tt.WantErr {
				assert.NoError(t, err)
			} else {
				assert.Error(t, err)
				assert.ErrorAs(t, err, &tt.ExpectedErr)
			}
		})

	}
}

func Test_sourceServiceImplementation_DeleteSource(t *testing.T) {

	mockRepo := mocks.SourceRepositoryMock{}

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
			mockRepo := test_utils.GetFieldByNameAndType(t, tt.Fields, "mockRepo", new(repository.SourceRepository))
			s := &sourceServiceImplementation{
				r: mockRepo.(repository.SourceRepository),
			}
			deleteIDSource := test_utils.GetArgByNameAndType(t, tt.Args, "id", "").(string)
			err := s.DeleteSource(deleteIDSource, "123")
			if !tt.WantErr {
				assert.NoError(t, err)
			} else {
				assert.Error(t, err)
				assert.ErrorAs(t, err, &tt.ExpectedErr)
			}
		})

	}
}
