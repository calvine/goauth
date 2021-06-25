package nullable

import (
	"encoding/json"
	"errors"
	"fmt"
	"testing"
	"time"

	goautherrors "github.com/calvine/goauth/core/errors"
)

var emptyTime = time.Time{}

const testTimeString = "2018-05-02T18:07:10-05:00"

func TestNullableTimeSetUnset(t *testing.T) {
	ns := NullableTime{}
	testValue, err := time.Parse(time.RFC3339, testTimeString)
	if err != nil {
		t.Error("error occurred while parsing time string for test.", err)
	}
	ns.Set(testValue)
	if ns.HasValue != true || ns.Value != testValue {
		t.Error("nullable struct in invalid state after Set call", ns)
	}
	ns.Unset()
	if ns.HasValue || ns.Value != defaultTimeValue {
		t.Error("nullable struct in invalid state after Unset call", ns)
	}
}

func TestNullableTimeScan(t *testing.T) {
	ns := NullableTime{}
	err := ns.Scan(nil)
	if err != nil {
		t.Error("Failed to scan nil into NullableTime", err, ns)
	}
	if ns.Value != emptyTime || ns.HasValue != false {
		t.Error("Nullable time has wrong value after scanning nil", ns)
	}
	testValue, err := time.Parse(time.RFC3339, testTimeString)
	if err != nil {
		t.Error("error occurred while parsing time string for test.", err)
	}
	err = ns.Scan(testValue)
	if err != nil {
		t.Error("Failed to scan nil into NullableTime", err, ns)
	}
	if ns.Value != testValue || ns.HasValue != true {
		errMsg := fmt.Sprintf("Nullable time has wrong value after scanning %v", testValue)
		t.Error(errMsg, ns)
	}
	testString := "abc"
	err = ns.Scan(testString)
	emptyErr := &goautherrors.WrongTypeError{}
	if !errors.As(err, emptyErr) {
		t.Error("Expected error to be of type WrongTypeError", err)
	}
	if ns.Value != emptyTime || ns.HasValue != false {
		errMsg := fmt.Sprintf("Nullable time has wrong value after scanning %v", testString)
		t.Error(errMsg, ns)
	}
}

func TestNullableTimeMarshalJson(t *testing.T) {
	ns := NullableTime{
		Value:    emptyTime,
		HasValue: false,
	}
	data, err := ns.MarshalJSON()
	if err != nil {
		t.Error("Failed to marshal null to JSON.", err)
	}
	if string(data) != "null" {
		t.Error("data from marshal was not null when underlaying nullable time was nil", data)
	}
	testValue, err := time.Parse(time.RFC3339, testTimeString)
	if err != nil {
		t.Error("error occurred while parsing time string for test.", err)
	}
	ns = NullableTime{
		Value:    testValue,
		HasValue: true,
	}
	data, err = ns.MarshalJSON()
	if err != nil {
		t.Error("Failed to marshal value to JSON.", ns, err)
	}
	if string(data) != fmt.Sprintf("\"%s\"", testTimeString) {
		t.Error("data from marshal was not what was expected", data, ns)
	}
}

func TestNullableTimeUnmarshalJson(t *testing.T) {
	testString := "null"
	ns := NullableTime{}
	err := ns.UnmarshalJSON([]byte(testString))
	if err != nil {
		t.Error("Failed to unmarshal null", err)
	}
	if ns.HasValue != false || ns.Value != emptyTime {
		t.Error("Unmarshaling null should result in a nullable time with an empty value and is null being true", ns)
	}
	testTime, err := time.Parse(time.RFC3339, testTimeString)
	if err != nil {
		t.Error("error occurred while parsing time string for test.", err)
	}
	testString = testTimeString
	err = ns.UnmarshalJSON([]byte("\"" + testString + "\""))
	if err != nil {
		t.Error("Failed to unmarshal", testString, err)
	}
	if ns.HasValue != true || !ns.Value.Equal(testTime) {
		t.Error("Unmarshaling 1.2 should result in a nullable time with a value of 1.2 and is null being false", ns, testTime)
	}
	testString = "false"
	err = ns.UnmarshalJSON([]byte(testString))
	if err == nil {
		t.Error("expected an error", err)
	}
	if ns.HasValue != false || ns.Value != emptyTime {
		t.Error("unmarshaling false should result in a nullable time with an empty value and is null being true", ns)
	}
}

func TestNullableTimeUnmarshalJsonPayload(t *testing.T) {
	testString := `{"t": "` + testTimeString + `"}`
	testReceiver := struct {
		TestValue NullableTime `json:"t"`
	}{}
	err := json.Unmarshal([]byte(testString), &testReceiver)
	if err != nil {
		t.Error("failed to unmarshal test string", err, testString)
	}
	testTime, err := time.Parse(time.RFC3339, testTimeString)
	if err != nil {
		t.Error("failed to parse test time for test", err)
	}
	if !testReceiver.TestValue.HasValue || !testReceiver.TestValue.Value.Equal(testTime) {
		t.Error("testRecevier TestValue is not equal to testTime", testReceiver, testTime)
	}
}
