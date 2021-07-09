package normalization

import (
	"reflect"

	"github.com/calvine/goauth/core/errors"
)

func NormalizeFloatValue(intValue interface{}) (float64, error) {
	switch ivt := intValue.(type) {
	case *float32, *float64:
		floatValue := reflect.ValueOf(ivt).Elem().Interface()
		return NormalizeFloatValue(floatValue)
	case float32:
		return float64(ivt), nil
	case float64:
		return ivt, nil
	default:
		return 0, errors.NewInvalidTypeError(reflect.TypeOf(intValue).String(), true)
	}
}
