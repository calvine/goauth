// +build ignore

package main

import (
	"strings"
	"text/template"

	"github.com/calvine/goauth/core/utilities"
)

type errorData struct {
	// Code is expected to be Pascal Case
	Code     string
	Message  string
	MetaData []string
}

func main() {
	funcMap := template.FuncMap{
		"ToUpper":            strings.ToUpper,
		"ToLower":            strings.ToLower,
		"UpperCaseFirstChar": utilities.UpperCaseFirstChar,
		"LowerCaseFirstChar": utilities.LowerCaseFirstChar,
	}
	errConstructorTemplate := template.Must(template.New("Error constructor template").Parse(errorConstructorTemplate)).Funcs(funcMap)
	errCodeTemplate := template.Must(template.New("Error code template").Parse(errorCodeTemplate)).Funcs(funcMap)
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
