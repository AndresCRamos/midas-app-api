package models

import (
	"time"
)

type Movement struct {
	UID          string
	OwnerId      string
	Name         string
	Description  string
	Amount       int64
	MovementDate time.Time
	CreatedAt    time.Time
	UpdatedAt    time.Time
}
