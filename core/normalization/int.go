package normalization

import (
	"reflect"

	coreerrors "github.com/calvine/goauth/core/errors"
	"github.com/calvine/richerror/errors"
)

func NormalizeIntValue(intValue interface{}) (int64, errors.RichError) {
	switch ivt := intValue.(type) {
	case *int8, *int16, *int, *int32, *int64:
		intValue := reflect.ValueOf(ivt).Elem().Interface()
		return NormalizeIntValue(intValue)
	case int8:
		return int64(ivt), nil
	case int16:
		return int64(ivt), nil
	case int32:
		return int64(ivt), nil
	case int:
		return int64(ivt), nil
	case int64:
		return ivt, nil
	default:
		return 0, coreerrors.NewInvalidTypeError(reflect.TypeOf(intValue).String(), true)
	}
}
