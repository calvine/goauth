package nullable

import (
	"encoding/json"
	"fmt"

	"github.com/calvine/goauth/core/errors"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/bsontype"
)

const defaultBoolValue = false

type NullableBool struct {
	HasValue bool
	Value    bool
}

// GetPointerCopy
func (nb *NullableBool) GetPointerCopy() *bool {
	if nb.HasValue {
		b := nb.Value
		return &b
	} else {
		return nil
	}
}

func (nb *NullableBool) Set(value bool) {
	nb.HasValue = true
	nb.Value = value
}

func (nb *NullableBool) Unset() {
	nb.HasValue = false
	nb.Value = defaultBoolValue
}

func (nb *NullableBool) MarshalJSON() ([]byte, error) {
	if !nb.HasValue {
		return []byte("null"), nil
	}
	return json.Marshal(nb.Value)
}

func (nb *NullableBool) UnmarshalJSON(data []byte) error {
	if string(data) == "null" {
		nb.HasValue = false
		nb.Value = false
		return nil
	}
	var value bool
	err := json.Unmarshal(data, &value)
	nb.HasValue = err == nil
	nb.Value = value
	return err
}

func (nb *NullableBool) Scan(value interface{}) error {
	switch t := value.(type) {
	case nil:
		nb.HasValue = false
		nb.Value = false
		return nil
	case bool:
		nb.HasValue = true
		nb.Value = t
		return nil
	default:
		nb.HasValue = false
		nb.Value = false
		err := errors.WrongTypeError{
			Actual:   fmt.Sprintf("%T", t),
			Expected: "bool",
		}
		return err
	}
}

func (nb *NullableBool) MarshalBSONValue() (bsontype.Type, []byte, error) {
	if !nb.HasValue {
		// make temp pointer with null value to marshal
		var b *bool
		return bson.MarshalValue(b)
	}
	return bson.MarshalValue(nb.Value)
}

func (nb *NullableBool) UnmarshalBSONValue(btype bsontype.Type, data []byte) error {
	// TODO: need to handle null value of data...
	var value bool
	err := bson.Unmarshal(data, &value)
	nb.HasValue = err == nil
	nb.Value = value
	return err
}
