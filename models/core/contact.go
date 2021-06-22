package core

import (
	"time"

	"github.com/calvine/goauth/models/nullable"
)

// Contact is a model that represents a contact method for a user like phone or email.
type Contact struct {
	ID             string                  `bson:"id"`
	UserID         string                  `bson:"userId"`
	Name           nullable.NullableString `bson:"name"`
	Principal      string                  `bson:"principal"`
	Type           string                  `bson:"type"`
	IsPrimary      bool                    `bson:"isPrimary"`
	CreatedByID    string                  `bson:"createdById"`
	CreatedOnDate  time.Time               `bson:"createdOnDate"`
	ModifiedByID   string                  `bson:"modifiedById"`
	ModifiedOnDate time.Time               `bson:"modifiedOnDate"`
}
