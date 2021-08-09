// +build ignore

package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"go/format"
	"io/ioutil"
	"os"
	"strings"
	"text/template"

	"github.com/calvine/goauth/core/errors/generator/templates"
	"github.com/calvine/goauth/core/utilities"
)

const (
	errorsDir = "../"
	codesDir  = "../codes/"
)

type dataItem struct {
	Name     string `json:"name"`
	DataType string `json:"dataType"`
}

type errorData struct {
	// Code is expected to be Pascal Case
	Code     string     `json:"code"`
	Message  string     `json:"message"`
	MetaData []dataItem `json:"metaData"`
}

type jsonErrorFileData struct {
	Data []errorData
}

type options struct {
	emitOutputFiles bool
	jsonErrorData   string
	errorFilePath   string
}

func main() {
	genOptions := options{}
	flag.BoolVar(&genOptions.emitOutputFiles, "emitOutputFiles", false, "if set source files will be created for the error.")
	flag.StringVar(&genOptions.jsonErrorData, "jsonErrorData", "", "a json string that represents the error details to generate code for.")
	flag.StringVar(&genOptions.errorFilePath, "errorFilePath", "", "The path the a json file that contains an array of error objects to have code generated for.")
	flag.Parse()
	fmt.Printf("%v\n\n", genOptions)
	funcMap := template.FuncMap{
		"ToUpper":            strings.ToUpper,
		"ToLower":            strings.ToLower,
		"UpperCaseFirstChar": utilities.UpperCaseFirstChar,
		"LowerCaseFirstChar": utilities.LowerCaseFirstChar,
	}
	errConstructorTemplate := template.Must(template.New("Error constructor template").Parse(templates.ErrorConstructorTemplate)).Funcs(funcMap)
	errCodeTemplate := template.Must(template.New("Error code template").Parse(templates.ErrorCodeTemplate)).Funcs(funcMap)
	errDataSlice := make([]errorData, 0)
	if genOptions.errorFilePath != "" {
		jsonErrorDataFileData, err := ioutil.ReadFile(genOptions.errorFilePath)
		if err != nil {
			errMsg := fmt.Sprintf("failed to open file %s - %s", genOptions.errorFilePath, err.Error())
			panic(errMsg)
		}
		json.Unmarshal(jsonErrorDataFileData, &errDataSlice)
	} else if genOptions.jsonErrorData != "" {
		var errData errorData
		err := json.Unmarshal([]byte(genOptions.jsonErrorData), &errData)
		if err != nil {
			// handle bad json data...
			errMsg := fmt.Sprintf("failed to unmarshal json data: %s", err.Error())
			panic(errMsg)
		}
		errDataSlice = append(errDataSlice, errData)
	} else {
		panic("jsonErrorData or errorFilePath are required.")
	}
	fmt.Printf("%v\n\n", errDataSlice)
	for _, data := range errDataSlice {
		constructorBuffer := bytes.NewBufferString("")
		err := errConstructorTemplate.Execute(constructorBuffer, data)
		if err != nil {
			fmt.Printf("failed to execute error constructor template: %s", err.Error())
			continue
		}
		errConstructorCode, err := format.Source([]byte(constructorBuffer.String()))
		if err != nil {
			fmt.Printf("%s", constructorBuffer)
			fmt.Printf("Failed to run format.Source on error code template: %s", err.Error())
			continue
		}

		codeBuffer := bytes.NewBufferString("")
		err = errCodeTemplate.Execute(codeBuffer, data)
		if err != nil {
			fmt.Printf("failed to execute error code template: %s", err.Error())
			continue
		}
		errCodeCode, err := format.Source([]byte(codeBuffer.String()))
		if err != nil {
			fmt.Printf("%s", codeBuffer)
			fmt.Printf("Failed to run format.Source on error code template: %s", err.Error())
			continue
		}

		if !genOptions.emitOutputFiles {
			fmt.Printf("\n\n************** Error Constructor Code **************\n\n")
			fmt.Fprint(os.Stdout, string(errConstructorCode))
			fmt.Printf("\n\n************** Error Code Code **************\n\n")
			fmt.Fprint(os.Stdout, string(errCodeCode))
		} else {
			// emit files...
		}
	}
}
