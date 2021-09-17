package normalization

import (
	"reflect"

	coreerrors "github.com/calvine/goauth/core/errors"
	"github.com/calvine/richerror/errors"
)

func NormalizeFloatValue(intValue interface{}) (float64, errors.RichError) {
	switch ivt := intValue.(type) {
	case *float32, *float64:
		floatValue := reflect.ValueOf(ivt).Elem().Interface()
		return NormalizeFloatValue(floatValue)
	case float32:
		return float64(ivt), nil
	case float64:
		return ivt, nil
	default:
		return 0, coreerrors.NewInvalidTypeError(reflect.TypeOf(intValue).String(), true)
	}
}
