package nullable

import (
	"encoding/json"
	"fmt"

	"github.com/calvine/goauth/core/errors"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/bsontype"
)

const (
	stringLengthByteLength = 4
	defaultStringValue     = ""
)

type NullableString struct {
	HasValue bool
	Value    string
}

// GetPointerCopy
func (ns *NullableString) GetPointerCopy() *string {
	if ns.HasValue {
		s := ns.Value
		return &s
	} else {
		return nil
	}
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
	// TODO: Is there a need here to make sure data is a quoted string?
	switch btype {
	case bsontype.Null:
		ns.Unset()
		return nil
	case bsontype.String:
		var value string
		// got some weird errors using bson.Unmarshal here, so coded a solution per the bson spec
		// https://bsonspec.org/spec.html

		// stringLength := binary.LittleEndian.Uint32(data[:stringLengthByteLength])
		// fmt.Println(stringLength)

		// start of string data begins after the uint32 string length value.
		startString := stringLengthByteLength
		// we want to clip off the null terminator.
		endString := len(data) - 1
		value = string(data[startString:endString])
		ns.Set(value)
		return nil
	default:
		return errors.WrongTypeError{Expected: bsontype.Array.String(), Actual: btype.String()}
	}
}
