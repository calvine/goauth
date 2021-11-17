package models

import (
	"strings"
	"time"

	"github.com/calvine/goauth/core"
	"github.com/calvine/goauth/core/nullable"
)

// TODO: Add validator for pre insert / update
// TODO: determine if the confirmation code needs an expiration date? or use redis for these short lived tokens?
// Contact is a model that represents a contact method for a user like phone or email.
type Contact struct {
	ID           string                  `bson:"-"`
	UserID       string                  `bson:"-"`
	Name         nullable.NullableString `bson:"name"`
	RawPrincipal string                  `bson:"rawPrincipal"`
	Principal    string                  `bson:"principal"`
	Type         string                  `bson:"type"`
	IsPrimary    bool                    `bson:"isPrimary"`
	// ConfirmationCode nullable.NullableString `bson:"confirmationCode"`
	ConfirmedDate nullable.NullableTime `bson:"confirmedDate"`
	AuditData     auditable             `bson:",inline"`
}

// TODO: write unit tests
func NewContact(userID, name, principal, contactType string, isPrimary bool) Contact {
	nameIsPopulated := name != ""
	normalizedPrincipal := NormalizeContactPrincipal(contactType, principal)
	// TODO: validate contact type is valid and contact principal is populate or valid
	return Contact{
		UserID:       userID,
		Name:         nullable.NullableString{HasValue: nameIsPopulated, Value: name},
		RawPrincipal: principal,
		Principal:    normalizedPrincipal,
		Type:         contactType,
		IsPrimary:    isPrimary,
	}
}

func (c *Contact) IsConfirmed() bool {
	now := time.Now()
	return c.ConfirmedDate.HasValue &&
		c.ConfirmedDate.Value.Before(now)
}

func NormalizeContactPrincipal(contactType, contactPrincipal string) string {
	var normalizedPrincipal string
	switch contactType {
	case core.CONTACT_TYPE_MOBILE:
		// remove dashes
		normalizedPrincipal = strings.ReplaceAll(contactPrincipal, "-", "")
	default:
		// lower case the contact
		normalizedPrincipal = strings.ToLower(contactPrincipal)
	}
	return normalizedPrincipal
}
