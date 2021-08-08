// +build ignore

package main

type errorData struct {
	// Code is expected to be Pascal Case
	Code     string
	Message  string
	MetaData []string
}

func main() {

}

const errorConstructorTemplate = `
package errors

import (
	"fmt"

	"github.com/calvine/goauth/core/errors/codes"
)

// New{{ .Code }}Error creates a new specific error
func New{{ .Code }}Error(actual string, includeStack bool) RichError {
	msg := {{ .Message }}
	err := NewRichError(codes.ErrCode{{ .Code }}, msg){{ range .MetaData }}.AddMetaData("actual", actual){{ end }}
	if includeStack {
		err = err.WithStack(1)
	}
	return err
}

`

const errorCodeTemplate = `
package codes

// ErrCode{{ .Code }} {{ .Message }}
const ErrCode{{ .Code }} = "{{ .Code }}"

`
