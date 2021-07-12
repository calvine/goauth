package errors

import (
	"fmt"

	"github.com/calvine/goauth/core/errors/codes"
)

func NewRepoQueryFailed(queryError error, includeStack bool) RichError {
	msg := fmt.Sprintf("repo query failed with error: %s", queryError.Error())
	err := NewRichError(codes.ErrCodeRepoQueryFailed, msg).AddError(queryError)
	if includeStack {
		err = err.WithStack(1)
	}
	return err
}
