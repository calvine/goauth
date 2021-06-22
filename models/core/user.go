package core

import (
	"time"

	"github.com/calvine/goauth/models/nullable"
)

// User represents a user in the system.
type User struct {
	ID                             string                  `bson:"id"`
	PasswordHash                   string                  `bson:"passwordHash"`
	Salt                           string                  `bson:"salt"`
	ConsecutiveFailedLoginAttempts int                     `bson:"consecutiveFailedLoginAttempts"`
	LockedOutUntil                 nullable.NullableTime   `bson:"lockedOutUntil"`
	CreatedByID                    string                  `bson:"createdBy"`
	CreatedOnDate                  time.Time               `bson:"createdOnDate"`
	ModifiedByID                   nullable.NullableString `bson:"modifiedById"`
	ModifiedOnDate                 nullable.NullableTime   `bson:"modifiedOnDate"`
}
