package models

import (
	"time"

	"github.com/calvine/goauth/core/nullable"
)

// TODO: Add validator for pre insert / update
// Profile represents personal profile data for a given user.
type Profile struct {
	ID          string                  `bson:"-"`
	UserID      string                  `bson:"-"`
	FirstName   nullable.NullableString `bson:"firstName"`
	MiddleName  nullable.NullableString `bson:"middleName"`
	LastName    nullable.NullableString `bson:"lastName"`
	DateOfBirth nullable.NullableTime   `bson:"dateOfBirth"`
	auditable
}

func NewProfile(userID, firstName, middleName, lastName string, dateOfBirth time.Time) Profile {
	firstNameIsPopulated := firstName != ""
	middleNameIsPopulated := middleName != ""
	lastNameIsPopulated := lastName != ""
	dobIsPopulated := !dateOfBirth.IsZero()
	return Profile{
		UserID:      userID,
		FirstName:   nullable.NullableString{HasValue: firstNameIsPopulated, Value: firstName},
		MiddleName:  nullable.NullableString{HasValue: middleNameIsPopulated, Value: middleName},
		LastName:    nullable.NullableString{HasValue: lastNameIsPopulated, Value: lastName},
		DateOfBirth: nullable.NullableTime{HasValue: dobIsPopulated, Value: dateOfBirth},
	}

}
