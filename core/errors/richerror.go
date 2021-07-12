package errors

import (
	"bytes"
	"fmt"
	"runtime"
	"strconv"
	"strings"
	"time"
)

var (
	indentString  = "\t" //""
	partSeperator = "\n" //" - "
)

type RichError struct {
	ErrCode     string                 `json:"code"`
	Message     string                 `json:"message"`
	Source      string                 `json:"source,omitempty"`
	Function    string                 `json:"function,omitempty"`
	LineNumber  string                 `json:"lineNumber,omitempty"`
	Stack       string                 `json:"stack,omitempty"`
	OccurredAt  time.Time              `json:"occurredAt"`
	InnerErrors []error                `json:"innerErrors"`
	MetaData    map[string]interface{} `json:"metaData"`
}

func NewRichError(errCode, message string) RichError {
	occurredAt := time.Now().UTC()
	err := RichError{
		ErrCode:    errCode,
		Message:    message,
		OccurredAt: occurredAt,
	}
	return err

}

func NewRichErrorWithStack(errCode, message string, stackOffset int) RichError {
	occurredAt := time.Now().UTC()
	err := RichError{
		ErrCode:    errCode,
		Message:    message,
		OccurredAt: occurredAt,
	}
	err = err.WithStack(stackOffset)
	return err

}

func (e RichError) WithStack(stackOffset int) RichError {
	baseStackOffset := 2
	// Here we initialize the slice to 10 because the runtime.Callers
	// function will not grow the slice as needed.
	var callerData []uintptr = make([]uintptr, 10)
	// Here we use 2 to remove the runtime.Callers call
	// and the call to the RichError.WithStack call.
	// This should leave only the relevant stack pieces
	numFrames := runtime.Callers(baseStackOffset+stackOffset, callerData)
	data := runtime.CallersFrames(callerData)
	stackBuffer := bytes.Buffer{}
	for i := 0; i < numFrames; i++ {
		nextFrame, _ := data.Next()
		if i == 0 {
			source := nextFrame.File
			// if len(source) > 0 {
			// 	sourceLastIndex := strings.LastIndex(nextFrame.File, string(os.PathSeparator))
			// 	source = source[sourceLastIndex+1:]
			// }

			functionName := nextFrame.Function
			if len(functionName) > 0 {
				functionNameLastIndex := strings.LastIndex(functionName, ".")
				functionName = functionName[functionNameLastIndex+1:]
			}
			e.Source = source
			e.Function = functionName
			e.LineNumber = strconv.Itoa(nextFrame.Line)
		}
		stackFrame := fmt.Sprintf("%sL:%d %v - %s:%d - %s%s", strings.Repeat(indentString, i), i+1, nextFrame.Entry, nextFrame.File, nextFrame.Line, nextFrame.Function, partSeperator)
		stackBuffer.WriteString(stackFrame)
	}
	e.Stack = string(stackBuffer.String())
	return e
}

func (e RichError) WithMetaData(metaData map[string]interface{}) RichError {
	e.MetaData = metaData
	return e
}

func (e RichError) WithErrors(errs []error) RichError {
	if errs != nil {
		if e.InnerErrors == nil {
			e.InnerErrors = make([]error, len(errs))
		}
		e.InnerErrors = append(e.InnerErrors, errs...)
	}
	return e
}

func (e RichError) AddSource(source string) RichError {
	e.Source = source
	return e
}

func (e RichError) AddFunction(function string) RichError {
	e.Function = function
	return e
}

func (e RichError) AddLineNumber(lineNumber string) RichError {
	e.LineNumber = lineNumber
	return e
}

func (e RichError) AddMetaData(key string, value interface{}) RichError {
	if e.MetaData == nil {
		e.MetaData = make(map[string]interface{}, 1)
	}
	e.MetaData[key] = value
	return e
}

func (e RichError) AddError(err error) RichError {
	if err != nil {
		if e.InnerErrors == nil {
			e.InnerErrors = make([]error, 1)
		}
		e.InnerErrors = append(e.InnerErrors, err)
	}
	return e
}

func (e RichError) Error() string {
	var messageBuffer bytes.Buffer
	timeStampMsg := fmt.Sprintf("TIMESTAMP: %s", e.OccurredAt.String())
	messageBuffer.WriteString(timeStampMsg)
	if e.Source != "" {
		sourceSection := fmt.Sprintf("%sSOURCE: %s", partSeperator, e.Source)
		messageBuffer.WriteString(sourceSection)
	}
	if e.Function != "" {
		functionSection := fmt.Sprintf("%sFUNCTION: %s", partSeperator, e.Function)
		messageBuffer.WriteString(functionSection)
	}
	if e.LineNumber != "" {
		LineNumberSection := fmt.Sprintf("%sLINE_NUM: %s", partSeperator, e.LineNumber)
		messageBuffer.WriteString(LineNumberSection)
	}
	if e.ErrCode != "" {
		errCodeSection := fmt.Sprintf("%sERRCODE: %s", partSeperator, e.ErrCode)
		messageBuffer.WriteString(errCodeSection)
	}
	if e.Message != "" {
		messageSection := fmt.Sprintf("%sMESSAGE: %s", partSeperator, e.Message)
		messageBuffer.WriteString(messageSection)
	}
	if e.Stack != "" {
		stackSection := fmt.Sprintf("%sSTACK: %s", partSeperator, e.Stack)
		messageBuffer.WriteString(stackSection)
	}
	if len(e.InnerErrors) > 0 {
		messageBuffer.WriteString("INNER ERRORS:")
		for i, err := range e.InnerErrors {
			innerErrMessage := fmt.Sprintf("%s%sERROR #%d: %s", partSeperator, strings.Repeat(indentString, i+1), i+1, err.Error())
			messageBuffer.WriteString(innerErrMessage)
		}
		messageBuffer.WriteString(partSeperator)
	}
	if len(e.MetaData) > 0 {
		messageBuffer.WriteString("METADATA:")
		for key, value := range e.MetaData {
			metaDataMsg := fmt.Sprintf("%s%s%s: %v", partSeperator, indentString, key, value)
			messageBuffer.WriteString(metaDataMsg)
		}
	}
	return messageBuffer.String()
}
