package models

type User struct {
	UID      string `firestore:"uid" json:"uid"`
	Alias    string `firestore:"alias" json:"alias"`
	Name     string `firestore:"name" json:"name"`
	LastName string `firestore:"last_name" json:"last_name"`
}

type UserCreate struct {
	Alias    string `firestore:"alias" json:"alias" binding:"required_without=Name"`
	Name     string `firestore:"name" json:"name" binding:"required_without=Alias"`
	LastName string `firestore:"last_name" json:"last_name" binding:"depends_on=Name"`
}

func (uc *UserCreate) ParseUser() User {
	return User{
		Alias:    uc.Alias,
		Name:     uc.Name,
		LastName: uc.LastName,
	}
}
