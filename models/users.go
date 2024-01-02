package models

type User struct {
	UID      string `firestore:"uid" json:"uid" binding:"required"`
	Alias    string `firestore:"alias" json:"alias"`
	Name     string `firestore:"name" json:"name"`
	LastName string `firestore:"last_name" json:"last_name"`
}
