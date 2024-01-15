package models

import "time"

type Source struct {
	UID         string `json:"uid" firestore:"uid" binding:"required"`
	Name        string `json:"name" firestore:"name" binding:"required"`
	OwnerId     string `json:"ownerId" firestore:"owner" binding:"required"`
	Description string `json:"description,omitempty" firestore:"description"`
	createdAt   time.Time
	updatedAt   time.Time
}

func (s *Source) NewCreationAtDate() {
	s.createdAt = time.Now()
}

func (s *Source) NewUpdatedAtDate() {
	s.updatedAt = time.Now()
}
