package models

import (
	"time"

	"github.com/calvine/goauth/core/nullable"
)

// TODO: Add validator for pre insert / update
// TODO: determine if the confirmation code needs an expiration date? or use redis for these short lived tokens?
// Contact is a model that represents a contact method for a user like phone or email.
type Contact struct {
	ID        string                  `bson:"-"`
	UserID    string                  `bson:"-"`
	Name      nullable.NullableString `bson:"name"`
	Principal string                  `bson:"principal"`
	Type      string                  `bson:"type"`
	IsPrimary bool                    `bson:"isPrimary"`
	// ConfirmationCode nullable.NullableString `bson:"confirmationCode"`
	ConfirmedDate nullable.NullableTime `bson:"confirmedDate"`
	AuditData     auditable             `bson:",inline"`
}

func NewContact(userID, name, principal, contactType string, isPrimary bool) Contact {
	nameIsPopulated := name != ""
	return Contact{
		UserID:    userID,
		Name:      nullable.NullableString{HasValue: nameIsPopulated, Value: name},
		Principal: principal,
		Type:      contactType,
		IsPrimary: isPrimary,
	}
}

func (c *Contact) IsConfirmed() bool {
	now := time.Now()
	return !c.ConfirmedDate.Value.IsZero() &&
		c.ConfirmedDate.Value.Before(now)
}
