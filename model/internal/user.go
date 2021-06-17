package internal

import "time"

// User represents a user in the system.
type User struct {
	ID                             string    `json:"id"`
	PasswordHash                   string    `json:"-"`
	Salt                           string    `json:"-"`
	ConsecutiveFailedLoginAttempts int       `json:"consecutiveFailedLoginAttempts"`
	LockedOutUntil                 time.Time `json:"lockedOutUntil"`
	CreatedByID                    string    `json:"-"`
	CreatedByDate                  time.Time `json:"-"`
	ModifiedByID                   string    `json:"-"`
	ModifiedByDate                 time.Time `json:"-"`
}
