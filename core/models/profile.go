package models

import (
	"time"

	"github.com/calvine/goauth/core/nullable"
)

// Profile represents personal profile data for a given user.
type Profile struct {
	ID          string                  `bson:"id"`
	UserID      string                  `bson:"userId"`
	FirstName   nullable.NullableString `bson:"firstName"`
	MiddleName  nullable.NullableString `bson:"middleName"`
	LastName    nullable.NullableString `bson:"lastName"`
	DateOfBirth time.Time               `bson:"dateOfBirth"`
	auditable
}
