package models

import "time"

type Source struct {
	UID         string    `json:"uid" firestore:"uid" binding:"required"`
	Name        string    `json:"name" firestore:"name" binding:"required"`
	OwnerId     string    `json:"ownerId" firestore:"owner" binding:"required"`
	Description string    `json:"description,omitempty" firestore:"description"`
	CreatedAt   time.Time `json:"-" firestore:"created_at"`
	UpdatedAt   time.Time `json:"-" firestore:"updated_at"`
}

func (s *Source) NewCreationAtDate() {
	s.CreatedAt = time.Now()
}

func (s *Source) NewUpdatedAtDate() {
	s.UpdatedAt = time.Now()
}
