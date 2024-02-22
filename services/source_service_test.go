package services

import (
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
				"source": models.Source{Name: "Success"},
			},
			WantErr:     false,
			ExpectedErr: nil,
			PreTest:     nil,
		},
		{
			Name:   "Fail to connect",
			Fields: fields,
			Args: test_utils.Args{
				"source": models.Source{Name: "CantConnect"},
			},
			WantErr:     true,
			ExpectedErr: error_utils.FirebaseUnknownError{},
			PreTest:     nil,
		},
		{
			Name:   "Duplicated Source",
			Fields: fields,
			Args: test_utils.Args{
				"source": models.Source{Name: "Duplicated", UID: "0"},
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
			mockRepo := test_utils.GetFieldByNameAndType[repository.SourceRepository](t, tt.Fields, "mockRepo")
			s := &sourceServiceImplementation{
				r: mockRepo,
			}
			sourceTest := test_utils.GetArgByNameAndType[models.Source](t, tt.Args, "source")
			_, err := s.CreateNewSource(sourceTest)
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
			mockRepo := test_utils.GetFieldByNameAndType[repository.SourceRepository](t, tt.Fields, "mockRepo")
			s := &sourceServiceImplementation{
				r: mockRepo,
			}
			userId := test_utils.GetArgByNameAndType[string](t, tt.Args, "userID")
			got, err := s.GetSourcesByUser(userId, 1)
			if !tt.WantErr {
				expected := test_utils.GetArgByNameAndType[util_models.PaginatedSearch[models.Source]](t, tt.Args, "expectedData")
				assert.NoError(t, err)
				assert.Equal(t, expected, got)
			} else {
				assert.Error(t, err)
				assert.ErrorAs(t, err, &tt.ExpectedErr)
			}
		})

	}
}

func Test_sourceServiceImplementation_GetMovementsBySourceAndDate(t *testing.T) {

	mockRepo := mocks.SourceRepositoryMock{}

	fields := test_utils.Fields{
		"mockRepo": mockRepo,
	}

	tests := []test_utils.TestCase{
		{
			Name:   "Success",
			Fields: fields,
			Args: test_utils.Args{
				"userID":           "0",
				"sourceID":         "0",
				"date_from":        time.Now().UTC().Add(-60 * 24 * time.Hour),
				"date_to":          time.Now().UTC(),
				"page":             1,
				"expectedPageSize": 50,
			},
			WantErr:     false,
			ExpectedErr: nil,
			PreTest:     nil,
		},
		{
			Name:   "Fail to connect",
			Fields: fields,
			Args: test_utils.Args{
				"sourceID":         "0",
				"userID":           "1",
				"date_from":        time.Now().UTC().Add(-60 * 24 * time.Hour),
				"date_to":          time.Now().UTC(),
				"page":             1,
				"expectedPageSize": 50,
			},
			WantErr:     true,
			ExpectedErr: error_utils.FirebaseUnknownError{},
			PreTest:     nil,
		},
		{
			Name:   "Not Found",
			Fields: fields,
			Args: test_utils.Args{
				"sourceID":         "0",
				"userID":           "2",
				"date_from":        time.Now().UTC().Add(-60 * 24 * time.Hour),
				"date_to":          time.Now().UTC(),
				"page":             1,
				"expectedPageSize": 50,
			},
			WantErr:     true,
			ExpectedErr: error_utils.FirestoreNotFoundError{},
			PreTest:     nil,
		},
		{
			Name:   "Different owner",
			Fields: fields,
			Args: test_utils.Args{
				"sourceID":         "0",
				"userID":           "3",
				"date_from":        time.Now().UTC().Add(-60 * 24 * time.Hour),
				"date_to":          time.Now().UTC(),
				"page":             1,
				"expectedPageSize": 50,
			},
			WantErr:     true,
			ExpectedErr: error_utils.SourceDifferentOwner{},
			PreTest:     nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.Name, func(t *testing.T) {
			if tt.PreTest != nil {
				tt.PreTest(t)
			}
			mockRepo := test_utils.GetFieldByNameAndType[repository.SourceRepository](t, tt.Fields, "mockRepo")
			s := &sourceServiceImplementation{
				r: mockRepo,
			}
			id := test_utils.GetArgByNameAndType[string](t, tt.Args, "sourceID")
			userID := test_utils.GetArgByNameAndType[string](t, tt.Args, "userID")
			page := test_utils.GetArgByNameAndType[int](t, tt.Args, "page")
			date_from := test_utils.GetArgByNameAndType[time.Time](t, tt.Args, "date_from")
			date_to := test_utils.GetArgByNameAndType[time.Time](t, tt.Args, "date_to")

			got, err := s.GetMovementsBySourceAndDate(id, userID, page, date_from, date_to)
			if !tt.WantErr {
				assert.NoError(t, err)
				assert.NotEmpty(t, got)
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
			mockRepo := test_utils.GetFieldByNameAndType[repository.SourceRepository](t, tt.Fields, "mockRepo")
			s := &sourceServiceImplementation{
				r: mockRepo,
			}
			id := test_utils.GetArgByNameAndType[string](t, tt.Args, "id")
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
				"source": models.Source{Name: "Success", UID: "0"},
			},
			WantErr:     false,
			ExpectedErr: nil,
			PreTest:     nil,
		},
		{
			Name:   "Fail to connect",
			Fields: fields,
			Args: test_utils.Args{
				"source": models.Source{Name: "Cant update", UID: "1"},
			},
			WantErr:     true,
			ExpectedErr: error_utils.FirebaseUnknownError{},
			PreTest:     nil,
		},
		{
			Name:   "Cant find",
			Fields: fields,
			Args: test_utils.Args{
				"source": models.Source{Name: "Cant update", UID: "2"},
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
			mockRepo := test_utils.GetFieldByNameAndType[repository.SourceRepository](t, tt.Fields, "mockRepo")
			s := &sourceServiceImplementation{
				r: mockRepo,
			}
			testSource := test_utils.GetArgByNameAndType[models.Source](t, tt.Args, "source")
			_, err := s.UpdateSource(testSource)
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
			mockRepo := test_utils.GetFieldByNameAndType[repository.SourceRepository](t, tt.Fields, "mockRepo")
			s := &sourceServiceImplementation{
				r: mockRepo,
			}
			deleteIDSource := test_utils.GetArgByNameAndType[string](t, tt.Args, "id")
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
