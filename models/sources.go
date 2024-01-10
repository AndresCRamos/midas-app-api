package models

import "time"

type Source struct {
	UID         string
	Name        string
	OwnerId     string
	Description string
	createdAt   time.Time
	updatedAt   time.Time
}

func (s *Source) NewCreationAtDate() {
	s.createdAt = time.Now()
}

func (s *Source) NewUpdatedAtDate() {
	s.updatedAt = time.Now()
}
