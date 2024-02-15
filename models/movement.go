package models

import (
	"time"
)

type Movement struct {
	UID          string    `firestore:"uid"`
	OwnerId      string    `firestore:"owner"`
	SourceID     string    `firestore:"source"`
	Name         string    `firestore:"name"`
	Description  string    `firestore:"description"`
	Amount       int64     `firestore:"amount"`
	MovementDate time.Time `firestore:"movement_date"`
	Tags         []string  `firestore:"tags"`
	CreatedAt    time.Time `firestore:"created_at"`
	UpdatedAt    time.Time `firestore:"updated_at"`
}

func (m *Movement) NewCreationAtDate() {
	m.CreatedAt = time.Now().UTC()
}

func (m *Movement) NewUpdatedAtDate() {
	m.UpdatedAt = time.Now().UTC()
}
