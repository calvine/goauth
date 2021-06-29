package nullable

import (
	"encoding/json"
	"fmt"

	"github.com/calvine/goauth/core/errors"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/bsontype"
)

const defaultFloat64Value = float64(0)

type NullableFloat64 struct {
	HasValue bool
	Value    float64
}

// GetPointerCopy
func (nf *NullableFloat64) GetPointerCopy() *float64 {
	if nf.HasValue {
		f := nf.Value
		return &f
	} else {
		return nil
	}
}

func (nf *NullableFloat64) Set(value float64) {
	nf.HasValue = true
	nf.Value = value
}

func (nf *NullableFloat64) Unset() {
	nf.HasValue = false
	nf.Value = defaultFloat64Value
}

func (nf *NullableFloat64) MarshalJSON() ([]byte, error) {
	if !nf.HasValue {
		return []byte("null"), nil
	}
	return json.Marshal(nf.Value)
}

func (nf *NullableFloat64) UnmarshalJSON(data []byte) error {
	if string(data) == "null" {
		nf.HasValue = false
		nf.Value = 0
		return nil
	}
	var value float64
	err := json.Unmarshal(data, &value)
	nf.HasValue = err == nil
	nf.Value = value
	return err
}

func (nf *NullableFloat64) Scan(value interface{}) error {
	switch t := value.(type) {
	case nil:
		nf.HasValue = false
		nf.Value = 0
		return nil
	case float64:
		nf.HasValue = true
		nf.Value = t
		return nil
	default:
		nf.HasValue = false
		nf.Value = 0
		err := errors.WrongTypeError{
			Actual:   fmt.Sprintf("%T", t),
			Expected: "float64",
		}
		return err
	}
}

func (nf *NullableFloat64) MarshalBSONValue() (bsontype.Type, []byte, error) {
	if !nf.HasValue {
		// make temp pointer with null value to marshal
		var f *float64
		return bson.MarshalValue(f)
	}
	return bson.MarshalValue(nf.Value)
}

func (nf *NullableFloat64) UnmarshalBSONValue(btype bsontype.Type, data []byte) error {
	switch btype {
	case bsontype.Null:
		nf.Unset()
		return nil
	case bsontype.Double:
		var value float64
		err := bson.Unmarshal(data, &value)
		if err != nil {
			return err
		}
		nf.Set(value)
		return nil
	default:
		return errors.WrongTypeError{Expected: bsontype.Double.String(), Actual: btype.String()}
	}
}
