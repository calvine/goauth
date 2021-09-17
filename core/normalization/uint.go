package normalization

import (
	"reflect"

	coreerrors "github.com/calvine/goauth/core/errors"
	"github.com/calvine/richerror/errors"
)

func NormalizeUintValue(uintValue interface{}) (uint64, errors.RichError) {
	switch ivt := uintValue.(type) {
	case *uint8, *uint16, *uint, *uint32, *uint64:
		uintValue := reflect.ValueOf(ivt).Elem().Interface()
		return NormalizeUintValue(uintValue)
	case uint8:
		return uint64(ivt), nil
	case uint16:
		return uint64(ivt), nil
	case uint32:
		return uint64(ivt), nil
	case uint:
		return uint64(ivt), nil
	case uint64:
		return ivt, nil
	default:
		return 0, coreerrors.NewInvalidTypeError(reflect.TypeOf(uintValue).String(), true)
	}
}
