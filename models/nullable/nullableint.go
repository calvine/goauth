package nullable

import (
	"encoding/json"
	"fmt"

	"github.com/calvine/goauth/models/errors"
	"go.mongodb.org/mongo-driver/bson"
)

const defaultIntValue = int(0)

type NullableInt struct {
	IsNull bool
	Value  int
}

func (ni *NullableInt) Set(value int) {
	ni.IsNull = false
	ni.Value = value
}

func (ni *NullableInt) Unset() {
	ni.IsNull = true
	ni.Value = defaultIntValue
}

func (ni *NullableInt) MarshalJSON() ([]byte, error) {
	if ni.IsNull {
		return []byte("null"), nil
	}
	return json.Marshal(ni.Value)
}

func (ni *NullableInt) UnmarshalJSON(data []byte) error {
	if string(data) == "null" {
		ni.IsNull = true
		ni.Value = 0
		return nil
	}
	var value int
	err := json.Unmarshal(data, &value)
	ni.IsNull = err != nil
	ni.Value = value
	return err
}

func (ni *NullableInt) Scan(value interface{}) error {
	switch t := value.(type) {
	case nil:
		ni.IsNull = true
		ni.Value = 0
		return nil
	case int:
		ni.IsNull = false
		ni.Value = t
		return nil
	default:
		ni.IsNull = true
		ni.Value = 0
		err := errors.WrongTypeError{
			Actual:   fmt.Sprintf("%T", t),
			Expected: "int",
		}
		return err
	}
}

func (ni *NullableInt) MarshalBSON() ([]byte, error) {
	if ni.IsNull {
		return nil, nil
	}
	return bson.Marshal(ni.Value)
}

func (ni *NullableInt) UnmarshalBSON(data []byte) error {
	// TODO: need to handle null value of data...
	var value int
	err := bson.Unmarshal(data, &value)
	ni.IsNull = err != nil
	ni.Value = value
	return err
}
