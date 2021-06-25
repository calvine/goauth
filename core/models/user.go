package models

import (
	"github.com/calvine/goauth/core/nullable"
)

// User represents a user in the system.
type User struct {
	Id                             string                `bson:"-"`
	PasswordHash                   string                `bson:"passwordHash"`
	Salt                           string                `bson:"salt"`
	ConsecutiveFailedLoginAttempts int                   `bson:"consecutiveFailedLoginAttempts"`
	LockedOutUntil                 nullable.NullableTime `bson:"lockedOutUntil"`
	LastLoginDate                  nullable.NullableTime `bson:"lastLoginDate"`
	AuditData                      auditable             `bson:",inline"`
}