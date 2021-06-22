package core

import (
	"time"

	"github.com/calvine/goauth/models/nullable"
)

// Address is a physical address.
type Address struct {
	ID             string                  `bson:"id"`
	UserID         string                  `bson:"userId"`
	Name           nullable.NullableString `bson:"name"`
	Line1          string                  `bson:"line1"`
	Line2          nullable.NullableString `bson:"line2"`
	City           string                  `bson:"city"`
	State          string                  `bson:"state"`
	PostalCode     string                  `bson:"postalCode"`
	IsPrimary      bool                    `bson:"isPrimary"`
	CreatedByID    string                  `bson:"createdById"`
	CreatedOnDate  time.Time               `bson:"createdOnDate"`
	ModifiedByID   nullable.NullableString `bson:"modifiedById"`
	ModifiedOnDate nullable.NullableTime   `bson:"modifiedOnDate"`
}
