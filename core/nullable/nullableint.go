package nullable

import (
	"encoding/json"
	"fmt"

	"github.com/calvine/goauth/core/errors"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/bsontype"
)

const defaultIntValue = int(0)

type NullableInt struct {
	HasValue bool
	Value    int
}

// GetPointerCopy
func (ni *NullableInt) GetPointerCopy() *int {
	if ni.HasValue {
		i := ni.Value
		return &i
	} else {
		return nil
	}
}

func (ni *NullableInt) Set(value int) {
	ni.HasValue = true
	ni.Value = value
}

func (ni *NullableInt) Unset() {
	ni.HasValue = false
	ni.Value = defaultIntValue
}

func (ni *NullableInt) MarshalJSON() ([]byte, error) {
	if !ni.HasValue {
		return []byte("null"), nil
	}
	return json.Marshal(ni.Value)
}

func (ni *NullableInt) UnmarshalJSON(data []byte) error {
	if string(data) == "null" {
		ni.HasValue = false
		ni.Value = 0
		return nil
	}
	var value int
	err := json.Unmarshal(data, &value)
	ni.HasValue = err == nil
	ni.Value = value
	return err
}

func (ni *NullableInt) Scan(value interface{}) error {
	switch t := value.(type) {
	case nil:
		ni.HasValue = false
		ni.Value = 0
		return nil
	case int:
		ni.HasValue = true
		ni.Value = t
		return nil
	default:
		ni.HasValue = false
		ni.Value = 0
		err := errors.NewWrongTypeError(fmt.Sprintf("%T", t), "int", true)
		return err
	}
}

func (ni *NullableInt) MarshalBSONValue() (bsontype.Type, []byte, error) {
	if !ni.HasValue {
		// make temp pointer with null value to marshal
		var i *int
		return bson.MarshalValue(i)
	}
	return bson.MarshalValue(ni.Value)
}

func (ni *NullableInt) UnmarshalBSONValue(btype bsontype.Type, data []byte) error {
	switch btype {
	case bsontype.Null:
		ni.Unset()
		return nil
	case bsontype.Int32:
		var value int
		err := bson.Unmarshal(data, &value)
		if err != nil {
			return err
		}
		ni.Set(value)
		return nil
	default:
		return errors.NewWrongTypeError(btype.String(), bsontype.Int32.String(), true)
	}
}
