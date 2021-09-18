package models

import (
	"github.com/calvine/goauth/core/nullable"
)

// TODO: Add validator for pre insert / update
// User represents a user in the system.
type User struct {
	ID                             string                `bson:"-"`
	PasswordHash                   string                `bson:"passwordHash"`
	ConsecutiveFailedLoginAttempts int                   `bson:"consecutiveFailedLoginAttempts"`
	LockedOutUntil                 nullable.NullableTime `bson:"lockedOutUntil"`
	LastLoginDate                  nullable.NullableTime `bson:"lastLoginDate"`
	// PasswordResetToken             nullable.NullableString `bson:"passwordResetToken"`
	// PasswordResetTokenExpiration   nullable.NullableTime   `bson:"passwordResetTokenExpiration"`
	AuditData auditable `bson:",inline"`
}

func NewUser() User {
	return User{}
}
