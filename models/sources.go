package models

import "time"

type Source struct {
	UID         string
	Name        string
	OwnerId     string
	Description string
	CreatedAt   time.Time
	UpdatedAt   time.Time
}
