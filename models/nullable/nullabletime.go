package nullable

import (
	"fmt"
	"time"

	"github.com/calvine/goauth/models/errors"
	"go.mongodb.org/mongo-driver/bson"
)

var defaultTimeValue = time.Time{}

type NullableTime struct {
	IsNull bool
	Value  time.Time
}

func (nt *NullableTime) Set(value time.Time) {
	nt.IsNull = false
	nt.Value = value
}

func (nt *NullableTime) Unset() {
	nt.IsNull = true
	nt.Value = defaultTimeValue
}

func (nt *NullableTime) MarshalJSON() ([]byte, error) {
	if nt.IsNull {
		return []byte("null"), nil
	}
	return nt.Value.MarshalJSON()
}

func (nt *NullableTime) UnmarshalJSON(data []byte) error {
	if string(data) == "null" {
		nt.IsNull = true
		nt.Value = time.Time{}
		return nil
	}
	var value time.Time
	err := value.UnmarshalJSON(data)
	nt.IsNull = err != nil
	nt.Value = value
	return err
}

func (nt *NullableTime) Scan(value interface{}) error {
	// TODO: How is time sent from the database? Do we need a type switch case? Does it vary by database?
	switch t := value.(type) {
	case nil:
		nt.IsNull = true
		nt.Value = time.Time{}
		return nil
	case time.Time:
		nt.IsNull = false
		nt.Value = t
		return nil
	default:
		nt.IsNull = true
		nt.Value = time.Time{}
		err := errors.WrongTypeError{
			Actual:   fmt.Sprintf("%T", t),
			Expected: "time.Time",
		}
		return err
	}
}

func (nt *NullableTime) MarshalBSON() ([]byte, error) {
	if nt.IsNull {
		return nil, nil
	}
	return bson.Marshal(nt.Value)
}

func (nt *NullableTime) UnmarshalBSON(data []byte) error {
	// TODO: need to handle null value of data...
	var value time.Time
	err := bson.Unmarshal(data, &value)
	nt.IsNull = err != nil
	nt.Value = value
	return err
}
