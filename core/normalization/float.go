package normalization

import (
	"errors"
	"reflect"
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
		// TODO: specific error here?
		return 0, errors.New("")
	}
}
