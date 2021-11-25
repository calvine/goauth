package errors

/* WARNING: This is GENERATED CODE Please do not edit. */

import (
	"github.com/calvine/richerror/errors"
)

// ErrCodeFailedTemplateRender An error occurred while rendering template
const ErrCodeFailedTemplateRender = "FailedTemplateRender"

// NewFailedTemplateRenderError creates a new specific error
func NewFailedTemplateRenderError(templateName string, renderError error, includeStack bool) errors.RichError {
	msg := "An error occurred while rendering template"
	err := errors.NewRichError(ErrCodeFailedTemplateRender, msg).AddMetaData("templateName", templateName).AddError(renderError).WithTags([]string{"template"})
	if includeStack {
		err = err.WithStack(1)
	}
	return err
}

func IsFailedTemplateRenderError(err errors.ReadOnlyRichError) bool {
	return err.GetErrorCode() == ErrCodeFailedTemplateRender
}
