package models

import (
	"time"
)

type Movement struct {
	UID          string    `firebase:"uid"`
	OwnerId      string    `firebase:"owner"`
	Name         string    `firebase:"name"`
	Description  string    `firebase:"description"`
	Amount       int64     `firebase:"amount"`
	MovementDate time.Time `firebase:"movement_date"`
	Tags         []string  `firebase:"tags"`
	CreatedAt    time.Time `firebase:"created_at"`
	UpdatedAt    time.Time `firebase:"updated_at"`
}

func (m *Movement) NewCreationAtDate() {
	m.CreatedAt = time.Now()
}

func (m *Movement) NewUpdatedAtDate() {
	m.UpdatedAt = time.Now()
}
