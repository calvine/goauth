package models

import (
	"time"

	"github.com/calvine/goauth/core/nullable"
)

type auditable struct {
	CreatedById    string                  `bson:"createdById"`
	CreatedOnDate  time.Time               `bson:"createdOnDate"`
	ModifiedByID   nullable.NullableString `bson:"modifiedById"`
	ModifiedOnDate nullable.NullableTime   `bson:"modifiedOnDate"`
}
