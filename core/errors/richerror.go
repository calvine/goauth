package errors

import (
	"bytes"
	"fmt"
	"runtime/debug"
	"time"
)

const (
	innerErrorMessageTemplate = " - INNER ERROR #%d: %s"
)

type RichError struct {
	ErrCode     string                 `json:"code"`
	Message     string                 `json:"message"`
	Source      string                 `json:"source,omitempty"`
	Stack       string                 `json:"stack,omitempty"`
	Area        string                 `json:"area,omitempty"`
	Location    string                 `json:"location,omitempty"`
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

func (e RichError) WithSource(source string) RichError {
	e.Source = source
	return e
}

func (e RichError) WithStach() RichError {
	stack := debug.Stack()
	e.Stack = string(stack)
	return e
}

func (e RichError) AddStack(stack string) RichError {
	e.Stack = string(stack)
	return e
}

func (e RichError) WithArea(area string) RichError {
	e.Area = area
	return e
}

func (e RichError) WithLocation(location string) RichError {
	e.Location = location
	return e
}

func (e RichError) WithMetaData(metaData map[string]interface{}) RichError {
	e.MetaData = metaData
	return e
}

func (e RichError) AddMetaData(key string, value interface{}) RichError {
	e.MetaData[key] = value
	return e
}

func (e RichError) WithError(err error) RichError {
	if err != nil {
		e.InnerErrors = append(e.InnerErrors, err)
	}
	return e
}

func (e RichError) WithErrors(errs []error) RichError {
	if errs != nil {
		e.InnerErrors = append(e.InnerErrors, errs...)
	}
	return e
}

func (e RichError) Error() string {
	timeString := e.OccurredAt.String()
	var messageBuffer bytes.Buffer
	messageBuffer.WriteString(timeString)
	if e.Source != "" {
		sourceSection := fmt.Sprintf("- %s ", e.Source)
		messageBuffer.WriteString(sourceSection)
	}
	if e.Area != "" {
		areaSection := fmt.Sprintf("- %s ", e.Area)
		messageBuffer.WriteString(areaSection)
	}
	if e.ErrCode != "" {
		errCodeSection := fmt.Sprintf("- %s ", e.ErrCode)
		messageBuffer.WriteString(errCodeSection)
	}
	if e.Location != "" {
		locationSection := fmt.Sprintf("- %s ", e.Location)
		messageBuffer.WriteString(locationSection)
	}
	if e.Message != "" {
		messageSection := fmt.Sprintf(": %s", e.Message)
		messageBuffer.WriteString(messageSection)
	}
	messageBuffer.WriteString(e.Message)
	for i, err := range e.InnerErrors {
		innerErrMessage := fmt.Sprintf(innerErrorMessageTemplate, i, err.Error())
		messageBuffer.WriteString(innerErrMessage)
	}
	return messageBuffer.String()
}
