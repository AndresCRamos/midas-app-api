package models

type User struct {
	UID      string `firestore:"uid"`
	Alias    string `firestore:"alias"`
	Name     string `firestore:"name"`
	LastName string `firestore:"last_name"`
}
