package models

import (
	"time"

	"github.com/calvine/goauth/core/nullable"
)

// TODO: Add validator for pre insert / update
// Profile represents personal profile data for a given user.
type Profile struct {
	Id          string                  `bson:"-"`
	UserId      string                  `bson:"-"`
	FirstName   nullable.NullableString `bson:"firstName"`
	MiddleName  nullable.NullableString `bson:"middleName"`
	LastName    nullable.NullableString `bson:"lastName"`
	DateOfBirth time.Time               `bson:"dateOfBirth"`
	auditable
}
