package nullable

import (
	"encoding/json"
	"fmt"

	"github.com/calvine/goauth/core/errors"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/bsontype"
)

const defaultFloat32Value = float32(0)

type NullableFloat32 struct {
	HasValue bool
	Value    float32
}

// GetPointerCopy
func (nf *NullableFloat32) GetPointerCopy() *float32 {
	if nf.HasValue {
		f := nf.Value
		return &f
	} else {
		return nil
	}
}

func (nf *NullableFloat32) Set(value float32) {
	nf.HasValue = true
	nf.Value = value
}

func (nf *NullableFloat32) Unset() {
	nf.HasValue = false
	nf.Value = defaultFloat32Value
}

func (nf *NullableFloat32) MarshalJSON() ([]byte, error) {
	if !nf.HasValue {
		return []byte("null"), nil
	}
	return json.Marshal(nf.Value)
}

func (nf *NullableFloat32) UnmarshalJSON(data []byte) error {
	if string(data) == "null" {
		nf.HasValue = false
		nf.Value = 0
		return nil
	}
	var value float32
	err := json.Unmarshal(data, &value)
	nf.HasValue = err == nil
	nf.Value = value
	return err
}

func (nf *NullableFloat32) Scan(value interface{}) error {
	switch t := value.(type) {
	case nil:
		nf.HasValue = false
		nf.Value = 0
		return nil
	case float32:
		nf.HasValue = true
		nf.Value = t
		return nil
	default:
		nf.HasValue = false
		nf.Value = 0
		err := errors.NewWrongTypeError(fmt.Sprintf("%T", t), "float32")
		return err
	}
}

func (nf *NullableFloat32) MarshalBSONValue() (bsontype.Type, []byte, error) {
	if !nf.HasValue {
		// make temp pointer with null value to marshal
		var f *float32
		return bson.MarshalValue(f)
	}
	return bson.MarshalValue(nf.Value)
}

func (nf *NullableFloat32) UnmarshalBSONValue(btype bsontype.Type, data []byte) error {
	switch btype {
	case bsontype.Null:
		nf.Unset()
		return nil
	case bsontype.Double:
		var value float32
		err := bson.Unmarshal(data, &value)
		if err != nil {
			return err
		}
		nf.Set(value)
		return err
	default:
		return errors.NewWrongTypeError(btype.String(), bsontype.Double.String())
	}
}
