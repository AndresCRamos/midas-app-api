package models

type User struct {
	UID      string `firestore:"uid" json:"uid" binding:"required"`
	Alias    string `firestore:"alias" json:"alias" binding:"required_without=Name"`
	Name     string `firestore:"name" json:"name" binding:"required_without=Alias"`
	LastName string `firestore:"last_name" json:"last_name" binding:"depends_on=Name"`
}
