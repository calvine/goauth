package templates

const (
	ErrorConstructorTemplate = `
package errors

/* WARNING: This is GENERATED CODE Please do not edit. */

import (
	"github.com/calvine/goauth/core/errors/codes"
)

// New{{ .Code }}Error creates a new specific error
func New{{ .Code }}Error({{ range .MetaData }}{{ .Name }} {{ .DataType }}, {{ end }}{{ if .IncludeMap }}fields map[string]interface{}, {{ end }}includeStack bool) RichError {
	msg := "{{ .Message }}"
	err := NewRichError(codes.ErrCode{{ .Code }}, msg){{ if .IncludeMap }}.WithMetaData(fields){{ end }}{{ range .MetaData }}.AddMetaData("{{ .Name }}", {{ .Name }}){{ end }}
	if includeStack {
		err = err.WithStack(1)
	}
	return err
}

`

	ErrorCodeTemplate = `
package codes

/* WARNING: This is GENERATED CODE Please do not edit. */

// ErrCode{{ .Code }} {{ .Message }}
const ErrCode{{ .Code }} = "{{ .Code }}"

`
)
