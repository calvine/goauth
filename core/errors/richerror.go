package errors

import (
	"bytes"
	"fmt"
	"runtime"
	"strconv"
	"strings"
	"time"
)

const (
	innerErrorMessageTemplate = " - INNER ERROR #%d: %s"
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
	return RichError{
		ErrCode:    errCode,
		Message:    message,
		OccurredAt: occurredAt,
	}
}

func NewRichErrorWithData(errCode, message, source, area, location string) RichError {
	err := NewRichError(errCode, message)
	err.Source = source
	err.Area = area
	err.Location = location
	return err
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

func (e RichError) WithStack() RichError {
	// Here we initialize the slice to 10 because the runtime.Callers
	// function will not grow the slice as needed.
	var callerData []uintptr = make([]uintptr, 10)
	// Here we use 2 to remove the runtime.Callers call
	// and the call to the RichError.WithStack call.
	// This should leave only the relevant stack pieces
	numFrames := runtime.Callers(2, callerData)
	data := runtime.CallersFrames(callerData)
	stackBuffer := bytes.Buffer{}
	for i := 0; i < numFrames; i++ {
		nextFrame, _ := data.Next()
		if i == 0 {
			e.Source = nextFrame.File
			e.Area = nextFrame.Function
			e.Location = strconv.Itoa(nextFrame.Line)
		}
		stackFrame := fmt.Sprintf("%s%v - %s:%d - %s\n", strings.Repeat("\t", i), nextFrame.Entry, nextFrame.File, nextFrame.Line, nextFrame.Function)
		stackBuffer.WriteString(stackFrame)
	}
	e.Stack = string(stackBuffer.String())
	return e
}

func (e RichError) WithMetaData(metaData map[string]interface{}) RichError {
	e.MetaData = metaData
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

func (e RichError) WithErrors(errs []error) RichError {
	if errs != nil {
		if e.InnerErrors == nil {
			e.InnerErrors = make([]error, len(errs))
		}
		e.InnerErrors = append(e.InnerErrors, errs...)
	}
	return e
}

func (e RichError) Error() string {
	timeString := e.OccurredAt.String()
	var messageBuffer bytes.Buffer
	messageBuffer.WriteString(timeString)
	if e.Source != "" {
		sourceSection := fmt.Sprintf(" - SOURCE: %s", e.Source)
		messageBuffer.WriteString(sourceSection)
	}
	if e.Area != "" {
		areaSection := fmt.Sprintf(" - AREA: %s", e.Area)
		messageBuffer.WriteString(areaSection)
	}
	if e.Location != "" {
		locationSection := fmt.Sprintf(" - LOCATION: %s", e.Location)
		messageBuffer.WriteString(locationSection)
	}
	if e.ErrCode != "" {
		errCodeSection := fmt.Sprintf(" - ERRCODE: %s", e.ErrCode)
		messageBuffer.WriteString(errCodeSection)
	}
	if e.Message != "" {
		messageSection := fmt.Sprintf(" - MESSAGE: %s", e.Message)
		messageBuffer.WriteString(messageSection)
	}
	if e.Stack != "" {
		stackSection := fmt.Sprintf(" - STACK: %s", e.Stack)
		messageBuffer.WriteString(stackSection)
	}
	messageBuffer.WriteString(e.Message)
	for i, err := range e.InnerErrors {
		innerErrMessage := fmt.Sprintf(innerErrorMessageTemplate, i, err.Error())
		messageBuffer.WriteString(innerErrMessage)
	}
	return messageBuffer.String()
}
