package internal

import "time"

// Contact is a model that represents a contact method for a user like phone or email.
type Contact struct {
	ID             string    `json:"id"`
	UserID         string    `json:"userId"`
	Name           string    `json:"name"`
	Principal      string    `json:"principal"`
	Type           string    `json:"type"`
	IsPrimary      bool      `json:"isPrimary"`
	CreatedByID    string    `json:"-"`
	CreatedByDate  time.Time `json:"-"`
	ModifiedByID   string    `json:"-"`
	ModifiedByDate time.Time `json:"-"`
}
