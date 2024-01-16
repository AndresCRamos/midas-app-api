package models

import "time"

type Source struct {
	UID         string    `firestore:"uid"`
	Name        string    `firestore:"name"`
	OwnerId     string    `firestore:"owner"`
	Description string    `firestore:"description"`
	CreatedAt   time.Time `firestore:"created_at"`
	UpdatedAt   time.Time `firestore:"updated_at"`
}

type SourceCreate struct {
	Name        string `json:"name" binding:"required"`
	OwnerId     string `json:"ownerId" binding:"required"`
	Description string `json:"description,omitempty"`
}

func (sc SourceCreate) ParseSource() Source {
	return Source{
		Name:        sc.Name,
		OwnerId:     sc.OwnerId,
		Description: sc.Description,
	}
}

type SourceRetrieve struct {
	UID         string    `json:"uid"`
	Name        string    `json:"name"`
	Description string    `json:"description,omitempty"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

func (sr *SourceRetrieve) ParseSource(src Source) {
	sr.UID = src.UID
	sr.Name = src.Name
	sr.Description = src.Description
	sr.CreatedAt = src.CreatedAt
	sr.UpdatedAt = src.UpdatedAt
}

func (s *Source) NewCreationAtDate() {
	s.CreatedAt = time.Now()
}

func (s *Source) NewUpdatedAtDate() {
	s.UpdatedAt = time.Now()
}
