package models

import (
	"encoding/json"
	"fmt"

	"github.com/calvine/goauth/models/errors"
	"go.mongodb.org/mongo-driver/bson"
)

type NullableString struct {
	IsNull bool
	Value  string
}

func (ns *NullableString) MarshalJSON() ([]byte, error) {
	if ns.IsNull {
		return []byte("null"), nil
	}
	return json.Marshal(ns.Value)
}

func (ns *NullableString) UnmarshalJSON(data []byte) error {
	// TODO: Is there a need here to make sure data is a quoted string?
	if string(data) == "null" {
		ns.IsNull = true
		ns.Value = ""
		return nil
	}
	var value string
	err := json.Unmarshal(data, &value)
	ns.IsNull = err != nil
	ns.Value = value
	return err
}

func (ns *NullableString) Scan(value interface{}) error {
	switch t := value.(type) {
	case nil:
		ns.IsNull = true
		ns.Value = ""
		return nil
	case string:
		ns.IsNull = false
		ns.Value = t
		return nil
	default:
		ns.IsNull = true
		ns.Value = ""
		err := errors.WrongTypeError{
			Actual:   fmt.Sprintf("%T", t),
			Expected: "string",
		}
		return err
	}
}

func (ns *NullableString) MarshalBSON() ([]byte, error) {
	if ns.IsNull {
		return nil, nil
	}
	return bson.Marshal(ns.Value)
}

func (ns *NullableString) UnmarshalBSON(data []byte) error {
	// TODO: need to handle null value of data...
	// TODO: Is there a need here to make sure data is a quoted string?
	var value string
	err := bson.Unmarshal(data, &value)
	ns.IsNull = err != nil
	ns.Value = value
	return err
}
