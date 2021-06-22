package nullable

import (
	"encoding/json"
	"fmt"

	"github.com/calvine/goauth/models/errors"
	"go.mongodb.org/mongo-driver/bson"
)

type NullableBool struct {
	IsNull bool
	Value  bool
}

func (nb *NullableBool) MarshalJSON() ([]byte, error) {
	if nb.IsNull {
		return []byte("null"), nil
	}
	return json.Marshal(nb.Value)
}

func (nb *NullableBool) UnmarshalJSON(data []byte) error {
	if string(data) == "null" {
		nb.IsNull = true
		nb.Value = false
		return nil
	}
	var value bool
	err := json.Unmarshal(data, &value)
	nb.IsNull = err != nil
	nb.Value = value
	return err
}

func (nb *NullableBool) Scan(value interface{}) error {
	switch t := value.(type) {
	case nil:
		nb.IsNull = true
		nb.Value = false
		return nil
	case bool:
		nb.IsNull = false
		nb.Value = t
		return nil
	default:
		nb.IsNull = true
		nb.Value = false
		err := errors.WrongTypeError{
			Actual:   fmt.Sprintf("%T", t),
			Expected: "bool",
		}
		return err
	}
}

func (nb *NullableBool) MarshalBSON() ([]byte, error) {
	if nb.IsNull {
		return nil, nil
	}
	return bson.Marshal(nb.Value)
}

func (nb *NullableBool) UnmarshalBSON(data []byte) error {
	// TODO: need to handle null value of data...
	var value bool
	err := bson.Unmarshal(data, &value)
	nb.IsNull = err != nil
	nb.Value = value
	return err
}
