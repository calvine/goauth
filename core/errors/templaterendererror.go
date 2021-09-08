package errors

/* WARNING: This is GENERATED CODE Please do not edit. */

import (
	"github.com/calvine/richerror/errors"
)

// ErrCodeTemplateRenderError An error occurred while rendering template
const ErrCodeTemplateRenderError = "TemplateRenderError"

// NewTemplateRenderErrorError creates a new specific error
func NewTemplateRenderErrorError(template string, renderError error, includeStack bool) errors.RichError {
	msg := "An error occurred while rendering template"
	err := errors.NewRichError(ErrCodeTemplateRenderError, msg).AddMetaData("template", template).AddError(renderError)
	if includeStack {
		err = err.WithStack(1)
	}
	return err
}

func IsTemplateRenderErrorError(err errors.ReadOnlyRichError) bool {
	return err.GetErrorCode() == ErrCodeTemplateRenderError
}
