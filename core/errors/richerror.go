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
	Area        string                 `json:"area,omitempty"`
	Location    string                 `json:"location,omitempty"`
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
			e.Source = nextFrame.File
			e.Area = nextFrame.Function
			e.Location = strconv.Itoa(nextFrame.Line)
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

func (e RichError) AddArea(area string) RichError {
	e.Area = area
	return e
}

func (e RichError) AddLocation(location string) RichError {
	e.Location = location
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
	if e.Area != "" {
		areaSection := fmt.Sprintf("%sAREA: %s", partSeperator, e.Area)
		messageBuffer.WriteString(areaSection)
	}
	if e.Location != "" {
		locationSection := fmt.Sprintf("%sLOCATION: %s", partSeperator, e.Location)
		messageBuffer.WriteString(locationSection)
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
			metaDataMsg := fmt.Sprintf("%s%s%s: %v,", partSeperator, indentString, key, value)
			messageBuffer.WriteString(metaDataMsg)
		}
	}
	return messageBuffer.String()
}
