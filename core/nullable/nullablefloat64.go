package nullable

import (
	"encoding/json"
	"fmt"

	"github.com/calvine/goauth/core/errors"
	"go.mongodb.org/mongo-driver/bson"
)

const defaultFloat64Value = float64(0)

type NullableFloat64 struct {
	IsNull bool
	Value  float64
}

func (nf *NullableFloat64) Set(value float64) {
	nf.IsNull = false
	nf.Value = value
}

func (nf *NullableFloat64) Unset() {
	nf.IsNull = true
	nf.Value = defaultFloat64Value
}

func (nf *NullableFloat64) MarshalJSON() ([]byte, error) {
	if nf.IsNull {
		return []byte("null"), nil
	}
	return json.Marshal(nf.Value)
}

func (nf *NullableFloat64) UnmarshalJSON(data []byte) error {
	if string(data) == "null" {
		nf.IsNull = true
		nf.Value = 0
		return nil
	}
	var value float64
	err := json.Unmarshal(data, &value)
	nf.IsNull = err != nil
	nf.Value = value
	return err
}

func (nf *NullableFloat64) Scan(value interface{}) error {
	switch t := value.(type) {
	case nil:
		nf.IsNull = true
		nf.Value = 0
		return nil
	case float64:
		nf.IsNull = false
		nf.Value = t
		return nil
	default:
		nf.IsNull = true
		nf.Value = 0
		err := errors.WrongTypeError{
			Actual:   fmt.Sprintf("%T", t),
			Expected: "float64",
		}
		return err
	}
}

func (nf *NullableFloat64) MarshalBSON() ([]byte, error) {
	if nf.IsNull {
		return nil, nil
	}
	return bson.Marshal(nf.Value)
}

func (nf *NullableFloat64) UnmarshalBSON(data []byte) error {
	// TODO: need to handle null value of data...
	var value float64
	err := bson.Unmarshal(data, &value)
	nf.IsNull = err != nil
	nf.Value = value
	return err
}
