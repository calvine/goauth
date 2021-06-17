package internal

import "time"

// Address is a physical address.
type Address struct {
	ID             string    `json:"id"`
	UserID         string    `json:"userId"`
	Name           string    `json:"name"`
	Line1          string    `json:"line1"`
	Line2          string    `json:"line2"`
	City           string    `json:"city"`
	State          string    `json:"state"`
	PostalCode     string    `json:"postalCode"`
	IsPrimary      bool      `json:"isPrimary"`
	CreatedByID    string    `json:"-"`
	CreatedByDate  time.Time `json:"-"`
	ModifiedByID   string    `json:"-"`
	ModifiedByDate time.Time `json:"-"`
}
