package internal

import "time"

// Profile represents personal profile data for a given user.
type Profile struct {
	ID             string    `json:"id"`
	UserID         string    `json:"userId"`
	FirstName      string    `json:"firstName"`
	MiddleName     string    `json:"middleName"`
	LastName       string    `json:"lastName"`
	DateOfBirth    time.Time `json:"dateOfBirth"`
	CreatedByID    string    `json:"-"`
	CreatedByDate  time.Time `json:"-"`
	ModifiedByID   string    `json:"-"`
	ModifiedByDate time.Time `json:"-"`
}
