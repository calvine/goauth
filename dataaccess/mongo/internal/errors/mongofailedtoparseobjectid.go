package errors

import (
	"fmt"

	coreerrors "github.com/calvine/goauth/core/errors"
	internalerrors "github.com/calvine/goauth/dataaccess/mongo/internal/errors/codes"
)

func NewMongoFailedToParseObjectID(objectID interface{}, includeStack bool) coreerrors.RichError {
	msg := fmt.Sprintf("failed to parse mongo object id: %s", objectID)
	err := coreerrors.NewRichError(internalerrors.ErrCodeFailedToParseObjectID, msg).AddMetaData("objectid", objectID)
	if includeStack {
		err = err.WithStack(1)
	}
	return err
}
