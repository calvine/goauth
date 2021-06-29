package nullable

import (
	"encoding/binary"
	"fmt"
	"time"

	coreErrors "github.com/calvine/goauth/core/errors"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/bsontype"
)

var defaultTimeValue = time.Time{}

type NullableTime struct {
	HasValue bool
	Value    time.Time
}

// GetPointerCopy
func (nt *NullableTime) GetPointerCopy() *time.Time {
	if nt.HasValue {
		t := nt.Value
		return &t
	} else {
		return nil
	}
}

func (nt *NullableTime) Set(value time.Time) {
	nt.HasValue = true
	nt.Value = value
}

func (nt *NullableTime) Unset() {
	nt.HasValue = false
	nt.Value = defaultTimeValue
}

func (nt *NullableTime) MarshalJSON() ([]byte, error) {
	if !nt.HasValue {
		return []byte("null"), nil
	}
	return nt.Value.MarshalJSON()
}

func (nt *NullableTime) UnmarshalJSON(data []byte) error {
	if string(data) == "null" {
		nt.HasValue = false
		nt.Value = time.Time{}
		return nil
	}
	var value time.Time
	err := value.UnmarshalJSON(data)
	nt.HasValue = err == nil
	nt.Value = value
	return err
}

func (nt *NullableTime) Scan(value interface{}) error {
	// TODO: How is time sent from the database? Do we need a type switch case? Does it vary by database?
	switch t := value.(type) {
	case nil:
		nt.HasValue = false
		nt.Value = time.Time{}
		return nil
	case time.Time:
		nt.HasValue = true
		nt.Value = t
		return nil
	default:
		nt.HasValue = false
		nt.Value = time.Time{}
		err := coreErrors.WrongTypeError{
			Actual:   fmt.Sprintf("%T", t),
			Expected: "time.Time",
		}
		return err
	}
}

func (nt *NullableTime) MarshalBSONValue() (bsontype.Type, []byte, error) {
	if !nt.HasValue {
		// make temp pointer with null value to marshal
		var t *time.Time
		return bson.MarshalValue(t)
	}
	return bson.MarshalValue(nt.Value)
}

func (nt *NullableTime) UnmarshalBSONValue(btype bsontype.Type, data []byte) error {
	switch btype {
	case bsontype.Null:
		nt.Unset()
		return nil
	case bsontype.DateTime:
		// According to mongodb dates are stored a unix timestamps with millisecond resolution per:
		// https://docs.mongodb.com/manual/reference/bson-types/#date
		int64Value := binary.LittleEndian.Uint64(data)
		// Go time.Time has nano second resolution.
		// Mongo returns the data with millisecons resolution
		// time.Unix takes seconds so we have to clip off the milliseconds for it to work properly.
		unixSeconds := int64Value / uint64(1000)
		// Get the number of milliseconds from the data from mongo
		milliSeconds := int64Value % 1000
		// convert the millisecond component to nano seconds.
		nanoSeconds := milliSeconds * uint64(time.Millisecond)
		timeValue := time.Unix(int64(unixSeconds), int64(nanoSeconds)).UTC()
		nt.Set(timeValue)
		return nil
	default:
		return coreErrors.WrongTypeError{Expected: bsontype.DateTime.String(), Actual: btype.String()}
	}
}
