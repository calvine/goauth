package nullable

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/calvine/goauth/core/errors"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/bsontype"
)

const defaultDurationValue = time.Duration(0)

type NullableDuration struct {
	HasValue bool
	Value    time.Duration
}

// GetPointerCopy
func (ni *NullableDuration) GetPointerCopy() *time.Duration {
	if ni.HasValue {
		i := ni.Value
		return &i
	} else {
		return nil
	}
}

func (ni *NullableDuration) Set(value time.Duration) {
	ni.HasValue = true
	ni.Value = value
}

func (ni *NullableDuration) Unset() {
	ni.HasValue = false
	ni.Value = defaultDurationValue
}

func (ni *NullableDuration) MarshalJSON() ([]byte, error) {
	if !ni.HasValue {
		return []byte("null"), nil
	}
	return json.Marshal(ni.Value)
}

func (ni *NullableDuration) UnmarshalJSON(data []byte) error {
	if string(data) == "null" {
		ni.HasValue = false
		ni.Value = 0
		return nil
	}
	var value time.Duration
	err := json.Unmarshal(data, &value)
	ni.HasValue = err == nil
	ni.Value = value
	return err
}

func (ni *NullableDuration) Scan(value interface{}) error {
	switch t := value.(type) {
	case nil:
		ni.HasValue = false
		ni.Value = defaultDurationValue
		return nil
	case int64:
		ni.HasValue = true
		ni.Value = time.Duration(t)
		return nil
	case time.Duration:
		ni.HasValue = true
		ni.Value = t
		return nil
	default:
		ni.HasValue = false
		ni.Value = defaultDurationValue
		err := errors.NewWrongTypeError(fmt.Sprintf("%T", t), "int64", true)
		return err
	}
}

func (ni *NullableDuration) MarshalBSONValue() (bsontype.Type, []byte, error) {
	if !ni.HasValue {
		// make temp pointer with null value to marshal
		var i *int64
		return bson.MarshalValue(i)
	}
	return bson.MarshalValue(int64(ni.Value))
}

func (ni *NullableDuration) UnmarshalBSONValue(btype bsontype.Type, data []byte) error {
	switch btype {
	case bsontype.Null:
		ni.Unset()
		return nil
	case bsontype.Int64:
		var value time.Duration
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
