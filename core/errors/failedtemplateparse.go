package errors

/* WARNING: This is GENERATED CODE Please do not edit. */

import (
	"github.com/calvine/richerror/errors"
)

// ErrCodeFailedTemplateParse An error occurred while parsing template
const ErrCodeFailedTemplateParse = "FailedTemplateParse"

// NewFailedTemplateParseError creates a new specific error
func NewFailedTemplateParseError(template string, parseError error, includeStack bool) errors.RichError {
	msg := "An error occurred while parsing template"
	err := errors.NewRichError(ErrCodeFailedTemplateParse, msg).AddMetaData("template", template).AddError(parseError).WithTags([]string{"template"})
	if includeStack {
		err = err.WithStack(1)
	}
	return err
}

func IsFailedTemplateParseError(err errors.ReadOnlyRichError) bool {
	return err.GetErrorCode() == ErrCodeFailedTemplateParse
}
