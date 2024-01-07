package models

import "github.com/go-playground/validator/v10"

type User struct {
	UID      string `firestore:"uid" json:"uid" binding:"required"`
	Alias    string `firestore:"alias" json:"alias" binding:"required_without=Name"`
	Name     string `firestore:"name" json:"name" binding:"required_without=Alias"`
	LastName string `firestore:"last_name" json:"last_name" binding:"required_with=Name"`
}

func (u *User) ValidateStructLevel(sl validator.StructLevel) {
	// Check if last_name is present only if name is present
	if u.Name == "" && u.LastName != "" {
		sl.ReportError(u.LastName, "LastName", "last_name", "dependsOn", "")
	}
}
