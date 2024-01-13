package models

import "time"

type Source struct {
	UID         string `json:"uid" binding:"required"`
	Name        string `json:"name" binding:"required"`
	OwnerId     string `json:"ownerId" binding:"required"`
	Description string `json:"description"`
	createdAt   time.Time
	updatedAt   time.Time
}

func (s *Source) NewCreationAtDate() {
	s.createdAt = time.Now()
}

func (s *Source) NewUpdatedAtDate() {
	s.updatedAt = time.Now()
}
