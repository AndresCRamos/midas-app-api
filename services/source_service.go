package services

import (
	"github.com/AndresCRamos/midas-app-api/models"
	"github.com/AndresCRamos/midas-app-api/repository"
	error_utils "github.com/AndresCRamos/midas-app-api/utils/errors"
)

type SourceService interface {
	CreateNewSource(Source models.Source) (models.Source, error)
	GetSourceByID(id string, userID string) (models.Source, error)
	UpdateSource(Source models.Source) (models.Source, error)
	DeleteSource(id string) error
}

type sourceServiceImplementation struct {
	r repository.SourceRepository
}

func NewSourceService(r repository.SourceRepository) *sourceServiceImplementation {
	return &sourceServiceImplementation{
		r: r,
	}
}

func (s *sourceServiceImplementation) CreateNewSource(source models.Source) (models.Source, error) {
	source, err := s.r.CreateNewSource(source)
	if err != nil {
		sourceServiceErr := error_utils.SourceServiceError{Err: err, Method: "Create"}
		return models.Source{}, sourceServiceErr
	}
	return source, nil
}

func (s *sourceServiceImplementation) GetSourceByID(id string, userID string) (models.Source, error) {
	res, err := s.r.GetSourceByID(id, userID)
	if err != nil {
		sourceServiceErr := error_utils.SourceServiceError{Err: err, Method: "Retrieve"}
		return models.Source{}, sourceServiceErr
	}
	return res, nil
}

func (s *sourceServiceImplementation) UpdateSource(source models.Source) (models.Source, error) {
	source, err := s.r.UpdateSource(source)
	if err != nil {
		sourceServiceErr := error_utils.SourceServiceError{Err: err, Method: "Update"}
		return models.Source{}, sourceServiceErr
	}
	return source, nil
}
func (s *sourceServiceImplementation) DeleteSource(id string) error {
	err := s.r.DeleteSource(id)
	if err != nil {
		sourceServiceErr := error_utils.SourceServiceError{Err: err, Method: "Delete"}
		return sourceServiceErr
	}
	return nil
}
