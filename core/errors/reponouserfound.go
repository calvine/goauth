package errors

import (
	"fmt"
	"strings"

	"github.com/calvine/goauth/core/errors/codes"
)

func NewRepoNoUserFoundError(key, value string, includeStack bool) RichError {
	msg := fmt.Sprintf("no user found for given field %s", key)
	err := NewRichError(codes.ErrCodeNoUserFound, msg).AddMetaData("key", key).AddMetaData("value", value)
	if includeStack {
		err = err.WithStack(1)
	}
	return err
}

func NewRepoNoUserFoundErrorWithFields(fields map[string]interface{}, includeStack bool) RichError {
	keys := make([]string, 0, len(fields))
	for k := range fields {
		keys = append(keys, k)
	}

	msg := fmt.Sprintf("no user found for given fields: %s", strings.Join(keys, ", "))
	err := NewRichError(codes.ErrCodeNoUserFound, msg).WithMetaData(fields)
	if includeStack {
		err = err.WithStack(1)
	}
	return err
}
