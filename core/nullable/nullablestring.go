package nullable

import (
	"encoding/json"
	"fmt"

	"github.com/calvine/goauth/core/errors"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/bsontype"
)

const defaultStringValue = ""

type NullableString struct {
	HasValue bool
	Value    string
}

func (ns *NullableString) Set(value string) {
	ns.HasValue = true
	ns.Value = value
}

func (ns *NullableString) Unset() {
	ns.HasValue = false
	ns.Value = defaultStringValue
}

func (ns *NullableString) MarshalJSON() ([]byte, error) {
	if !ns.HasValue {
		return []byte("null"), nil
	}
	return json.Marshal(ns.Value)
}

func (ns *NullableString) UnmarshalJSON(data []byte) error {
	// TODO: Is there a need here to make sure data is a quoted string?
	if string(data) == "null" {
		ns.HasValue = false
		ns.Value = ""
		return nil
	}
	var value string
	err := json.Unmarshal(data, &value)
	ns.HasValue = err == nil
	ns.Value = value
	return err
}

func (ns *NullableString) Scan(value interface{}) error {
	switch t := value.(type) {
	case nil:
		ns.HasValue = false
		ns.Value = ""
		return nil
	case string:
		ns.HasValue = true
		ns.Value = t
		return nil
	default:
		ns.HasValue = false
		ns.Value = ""
		err := errors.WrongTypeError{
			Actual:   fmt.Sprintf("%T", t),
			Expected: "string",
		}
		return err
	}
}

func (ns *NullableString) MarshalBSONValue() (bsontype.Type, []byte, error) {
	if !ns.HasValue {
		// make temp pointer with null value to marshal
		var s *string
		return bson.MarshalValue(s)
	}
	return bson.MarshalValue(ns.Value)
}

func (ns *NullableString) UnmarshalBSONValue(btype bsontype.Type, data []byte) error {
	// TODO: need to handle null value of data...
	// TODO: Is there a need here to make sure data is a quoted string?
	var value string
	err := bson.Unmarshal(data, &value)
	ns.HasValue = err == nil
	ns.Value = value
	return err
}
