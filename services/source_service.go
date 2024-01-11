package services

import (
	"github.com/AndresCRamos/midas-app-api/models"
	"github.com/AndresCRamos/midas-app-api/repository"
	error_utils "github.com/AndresCRamos/midas-app-api/utils/errors"
)

type SourceService interface {
	CreateNewSource(Source models.Source) error
	GetSourceByID(id string) (models.Source, error)
	UpdateNewSource(id string, Source models.Source) error
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

func (s *sourceServiceImplementation) CreateNewSource(source models.Source) error {
	err := s.r.CreateNewSource(source)
	if err != nil {
		sourceServiceErr := error_utils.SourceServiceError{Err: err, Method: "Create"}
		return sourceServiceErr
	}
	return nil
}

func (s *sourceServiceImplementation) GetSourceByID(id string) (models.Source, error) {
	res, err := s.r.GetSourceByID(id)
	if err != nil {
		sourceServiceErr := error_utils.SourceServiceError{Err: err, Method: "Retrieve"}
		return models.Source{}, sourceServiceErr
	}
	return res, nil
}

func (s *sourceServiceImplementation) UpdateNewSource(id string, source models.Source) error {
	err := s.r.UpdateNewSource(source)
	if err != nil {
		sourceServiceErr := error_utils.SourceServiceError{Err: err, Method: "Update"}
		return sourceServiceErr
	}
	return nil
}
func (s *sourceServiceImplementation) DeleteSource(id string) error {
	err := s.r.DeleteSource(id)
	if err != nil {
		sourceServiceErr := error_utils.SourceServiceError{Err: err, Method: "Delete"}
		return sourceServiceErr
	}
	return nil
}
