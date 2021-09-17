package normalization

import (
	"reflect"
	"strings"

	coreerrors "github.com/calvine/goauth/core/errors"
	"github.com/calvine/richerror/errors"
)

func NormalizeStringValue(value interface{}) (string, errors.RichError) {
	switch svt := value.(type) {
	case *string:
		stringValue := reflect.ValueOf(svt).Elem().Interface()
		return NormalizeStringValue(stringValue)
	case string:
		return strings.ToUpper(svt), nil
	case []byte:
		sValue := string(svt)
		return strings.ToUpper(sValue), nil
	default:
		return "", coreerrors.NewInvalidTypeError(reflect.TypeOf(value).String(), true)
	}
}
