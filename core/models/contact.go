package models

import (
	"net/mail"
	"strings"
	"time"

	"github.com/calvine/goauth/core"
	coreerrors "github.com/calvine/goauth/core/errors"
	"github.com/calvine/goauth/core/nullable"
	"github.com/calvine/richerror/errors"
)

// TODO: Add validator for pre insert / update
// TODO: determine if the confirmation code needs an expiration date? or use redis for these short lived tokens?
// Contact is a model that represents a contact method for a user like phone or email.
type Contact struct {
	ID            string                  `bson:"-"`
	UserID        string                  `bson:"-"`
	Name          nullable.NullableString `bson:"name"`
	RawPrincipal  string                  `bson:"rawPrincipal"`
	Principal     string                  `bson:"principal"`
	Type          string                  `bson:"type"`
	IsPrimary     bool                    `bson:"isPrimary"`
	ConfirmedDate nullable.NullableTime   `bson:"confirmedDate"`
	AuditData     auditable               `bson:",inline"`
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

func IsValidContactType(contactType string) errors.RichError {
	switch contactType {
	case core.CONTACT_TYPE_EMAIL:
		return nil
	case core.CONTACT_TYPE_MOBILE:
		return nil
	}
	return coreerrors.NewInvalidContactTypeError(contactType, true)
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

func IsValidNormalizedContactPrincipal(contactType, normalizedContactPrincipal string) errors.RichError {
	if normalizedContactPrincipal == "" {
		// an empty string is never valid...
		return coreerrors.NewInvalidContactPrincipalError(normalizedContactPrincipal, contactType, true)
	}
	// TODO: implement this for mobile...
	switch contactType {
	case core.CONTACT_TYPE_EMAIL:
		_, err := mail.ParseAddress(normalizedContactPrincipal)
		if err != nil {
			return coreerrors.NewInvalidContactPrincipalError(normalizedContactPrincipal, contactType, true)
		}
		return nil
	default:
		return nil
	}
}
