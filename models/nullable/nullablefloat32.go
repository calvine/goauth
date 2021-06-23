package nullable

import (
	"encoding/json"
	"fmt"

	"github.com/calvine/goauth/models/errors"
	"go.mongodb.org/mongo-driver/bson"
)

const defaultFloat32Value = float32(0)

type NullableFloat32 struct {
	IsNull bool
	Value  float32
}

func (nf *NullableFloat32) Set(value float32) {
	nf.IsNull = false
	nf.Value = value
}

func (nf *NullableFloat32) Unset() {
	nf.IsNull = true
	nf.Value = defaultFloat32Value
}

func (nf *NullableFloat32) MarshalJSON() ([]byte, error) {
	if nf.IsNull {
		return []byte("null"), nil
	}
	return json.Marshal(nf.Value)
}

func (nf *NullableFloat32) UnmarshalJSON(data []byte) error {
	if string(data) == "null" {
		nf.IsNull = true
		nf.Value = 0
		return nil
	}
	var value float32
	err := json.Unmarshal(data, &value)
	nf.IsNull = err != nil
	nf.Value = value
	return err
}

func (nf *NullableFloat32) Scan(value interface{}) error {
	switch t := value.(type) {
	case nil:
		nf.IsNull = true
		nf.Value = 0
		return nil
	case float32:
		nf.IsNull = false
		nf.Value = t
		return nil
	default:
		nf.IsNull = true
		nf.Value = 0
		err := errors.WrongTypeError{
			Actual:   fmt.Sprintf("%T", t),
			Expected: "float32",
		}
		return err
	}
}

func (nf *NullableFloat32) MarshalBSON() ([]byte, error) {
	if nf.IsNull {
		return nil, nil
	}
	return bson.Marshal(nf.Value)
}

func (nf *NullableFloat32) UnmarshalBSON(data []byte) error {
	// TODO: need to handle null value of data...
	var value float32
	err := bson.Unmarshal(data, &value)
	nf.IsNull = err != nil
	nf.Value = value
	return err
}
