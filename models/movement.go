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

type MovementCreate struct {
	SourceID     string    `json:"source" binding:"required"`
	Name         string    `json:"name" binding:"required"`
	Description  string    `json:"description"`
	Amount       int64     `json:"amount" binding:"required"`
	MovementDate time.Time `json:"movement_date" binding:"required"`
	Tags         []string  `json:"tags"`
}

func (mv *MovementCreate) ParseMovement() Movement {
	return Movement{
		SourceID:     mv.SourceID,
		Name:         mv.Name,
		Description:  mv.Description,
		Amount:       mv.Amount,
		MovementDate: mv.MovementDate,
		Tags:         mv.Tags,
	}
}

type MovementRetrieve struct {
	UID          string    `json:"uid"`
	OwnerId      string    `json:"owner"`
	SourceID     string    `json:"source"`
	Name         string    `json:"name"`
	Description  string    `json:"description"`
	Amount       int64     `json:"amount"`
	MovementDate time.Time `json:"movement_date"`
	Tags         []string  `json:"tags"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

func (mv *MovementRetrieve) ParseMovement(movement Movement) {
	mv.UID = movement.UID
	mv.OwnerId = movement.OwnerId
	mv.SourceID = movement.SourceID
	mv.Name = movement.Name
	mv.Description = movement.Description
	mv.Amount = movement.Amount
	mv.MovementDate = movement.MovementDate
	mv.Tags = movement.Tags
	mv.CreatedAt = movement.CreatedAt
	mv.UpdatedAt = movement.UpdatedAt
}

type MovementUpdate struct {
	SourceID     string    `json:"source"`
	Name         string    `json:"name"`
	Description  string    `json:"description"`
	Amount       int64     `json:"amount"`
	MovementDate time.Time `json:"movement_date"`
	Tags         []string  `json:"tags"`
}

func (mu *MovementUpdate) ParseMovement() Movement {
	return Movement{
		SourceID:     mu.SourceID,
		Name:         mu.Name,
		Description:  mu.Description,
		Amount:       mu.Amount,
		MovementDate: mu.MovementDate,
		Tags:         mu.Tags,
	}
}
