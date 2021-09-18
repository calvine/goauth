package models

import (
	"github.com/calvine/goauth/core/nullable"
)

// TODO: Add validator for pre insert / update
// Address is a physical address.
type Address struct {
	ID         string                  `bson:"-"`
	UserID     string                  `bson:"-"`
	Name       nullable.NullableString `bson:"name"`
	Line1      string                  `bson:"line1"`
	Line2      nullable.NullableString `bson:"line2"`
	City       string                  `bson:"city"`
	State      string                  `bson:"state"`
	PostalCode string                  `bson:"postalCode"`
	IsPrimary  bool                    `bson:"isPrimary"`
	auditable
}

func NewAddress(userID, name, line1, line2, city, state, postalCode string, isPiramry bool) Address {
	nameIsPopulated := name != ""
	line2IsPopulated := line2 != ""
	return Address{
		UserID:     userID,
		Name:       nullable.NullableString{HasValue: nameIsPopulated, Value: name},
		Line1:      line1,
		Line2:      nullable.NullableString{HasValue: line2IsPopulated, Value: line2},
		City:       city,
		State:      state,
		PostalCode: postalCode,
		IsPrimary:  isPiramry,
	}
}
