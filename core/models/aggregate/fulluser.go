package aggregate

import "github.com/calvine/goauth/core/models"

type FullUser struct {
	models.User
	Addresses []models.Address `bson:"addresses"`
	Contacts  []models.Contact `bson:"contacts"`
	Profile   models.Profile   `bson:"profile"`
}

func NewFullUser(user models.User) FullUser {
	return FullUser{
		User: user,
	}
}

func NewFullUserWithData(user models.User, addresses []models.Address, contacts []models.Contact, profile *models.Profile) FullUser {
	if addresses == nil {
		addresses = make([]models.Address, 0)
	}
	if contacts == nil {
		contacts = make([]models.Contact, 0)
	}
	return FullUser{
		User:      user,
		Addresses: addresses,
		Contacts:  contacts,
		Profile:   *profile,
	}
}
