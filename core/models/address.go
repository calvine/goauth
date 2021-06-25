package models

import (
	"github.com/calvine/goauth/core/nullable"
)

// Address is a physical address.
type Address struct {
	Id         string                  `bson:"-"`
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
