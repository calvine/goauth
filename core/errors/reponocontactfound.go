package errors

import (
	"fmt"
	"strings"

	"github.com/calvine/goauth/core/errors/codes"
)

func NewRepoNoContactFoundError(key, value string, includeStack bool) RichError {
	msg := fmt.Sprintf("no contact found for given field %s", key)
	err := NewRichError(codes.ErrCodeNoContactFound, msg).AddMetaData("key", key).AddMetaData("value", value)
	if includeStack {
		err = err.WithStack(1)
	}
	return err
}

func NewRepoNoContactFoundErrorWithFields(fields map[string]interface{}, includeStack bool) RichError {
	keys := make([]string, 0, len(fields))
	for k := range fields {
		keys = append(keys, k)
	}

	msg := fmt.Sprintf("no contact found for given fields: %s", strings.Join(keys, ", "))
	err := NewRichError(codes.ErrCodeNoContactFound, msg).WithMetaData(fields)
	if includeStack {
		err = err.WithStack(1)
	}
	return err
}
