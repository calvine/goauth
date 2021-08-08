// +build ignore

package main

import (
	"bytes"
	"fmt"
	"go/format"
	"os"
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
	example := errorData{
		Code:    "NoUserIdFound",
		Message: "failed to get item by user id",
		MetaData: []string{
			"userId",
		},
	}
	funcMap := template.FuncMap{
		"ToUpper":            strings.ToUpper,
		"ToLower":            strings.ToLower,
		"UpperCaseFirstChar": utilities.UpperCaseFirstChar,
		"LowerCaseFirstChar": utilities.LowerCaseFirstChar,
	}
	errConstructorTemplate := template.Must(template.New("Error constructor template").Parse(errorConstructorTemplate)).Funcs(funcMap)
	errCodeTemplate := template.Must(template.New("Error code template").Parse(errorCodeTemplate)).Funcs(funcMap)
	constructorBuffer := bytes.NewBufferString("")
	codeBuffer := bytes.NewBufferString("")
	errConstructorTemplate.Execute(constructorBuffer, example)
	errCodeTemplate.Execute(codeBuffer, example)

	errConstructorCode, _ := format.Source([]byte(constructorBuffer.String()))
	errCodeCode, _ := format.Source([]byte(codeBuffer.String()))

	fmt.Fprint(os.Stdout, string(errConstructorCode))
	fmt.Fprint(os.Stdout, string(errCodeCode))
}

const errorConstructorTemplate = `
package errors

import (
	"fmt"

	"github.com/calvine/goauth/core/errors/codes"
)

// New{{ .Code }}Error creates a new specific error
func New{{ .Code }}Error({{ range .MetaData }}{{ . }} interface{},{{ end }} includeStack bool) RichError {
	msg := "{{ .Message }}"
	err := NewRichError(codes.ErrCode{{ .Code }}, msg){{ range .MetaData }}.AddMetaData("{{ . }}", {{ . }}){{ end }}
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
