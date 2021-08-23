package nullable

import (
	"fmt"
	"testing"

	coreerrors "github.com/calvine/goauth/core/errors"
	"github.com/calvine/richerror/errors"
)

func TestNullableIntGetPointerCopy(t *testing.T) {
	ni := NullableInt{}
	ni.Set(3)
	niCopy := ni.GetPointerCopy()
	if *niCopy != ni.Value {
		t.Error("niCopy value should be the same as ni Value", ni, niCopy)
	}
	if &ni.Value == niCopy {
		t.Error("the address of ni.Value and niCopy should be different", &ni.Value, &niCopy)
	}
	ni.Unset()
	niCopy = ni.GetPointerCopy()
	if niCopy != nil {
		t.Error("niCopy should be nil because ni HasValue is false", ni, niCopy)
	}
}

func TestNullableIntSetUnset(t *testing.T) {
	ni := NullableInt{}
	testValue := int(1)
	ni.Set(testValue)
	if ni.HasValue != true || ni.Value != testValue {
		t.Error("nullable struct in invalid state after Set call", ni)
	}
	ni.Unset()
	if ni.HasValue || ni.Value != defaultIntValue {
		t.Error("nullable struct in invalid state after Unset call", ni)
	}
}

func TestNullableIntScan(t *testing.T) {
	ns := NullableInt{}
	err := ns.Scan(nil)
	if err != nil {
		t.Error("Failed to scan nil into NullableInt", err, ns)
	}
	if ns.Value != 0 || ns.HasValue != false {
		t.Error("Nullable int has wrong value after scanning nil", ns)
	}
	testValue := 2
	err = ns.Scan(testValue)
	if err != nil {
		t.Error("Failed to scan nil into NullableInt", err, ns)
	}
	if ns.Value != testValue || ns.HasValue != true {
		errMsg := fmt.Sprintf("Nullable int has wrong value after scanning %v", testValue)
		t.Error(errMsg, ns)
	}
	testString := "abc"
	err = ns.Scan(testString)
	if err != nil && err.(errors.RichError).GetErrorCode() != coreerrors.ErrCodeWrongType {
		t.Error("Expected error to be of type WrongTypeError", err)
	}
	if ns.Value != 0 || ns.HasValue != false {
		errMsg := fmt.Sprintf("Nullable int has wrong value after scanning %v", testString)
		t.Error(errMsg, ns)
	}
}

func TestNullableIntMarshalJson(t *testing.T) {
	ns := NullableInt{
		Value:    0,
		HasValue: false,
	}
	data, err := ns.MarshalJSON()
	if err != nil {
		t.Error("Failed to marshal null to JSON.", err)
	}
	if string(data) != "null" {
		t.Error("data from marshal was not null when underlaying nullable int was nil", data)
	}
	ns = NullableInt{
		Value:    -2,
		HasValue: true,
	}
	data, err = ns.MarshalJSON()
	if err != nil {
		t.Error("Failed to marshal null to JSON.", err)
	}
	if string(data) != "-2" {
		t.Error("data from marshal was not what was expected", data, ns)
	}
}

func TestNullableIntUnmarshalJson(t *testing.T) {
	testString := "null"
	ns := NullableInt{}
	err := ns.UnmarshalJSON([]byte(testString))
	if err != nil {
		t.Error("Failed to unmarshal null", err)
	}
	if ns.HasValue != false || ns.Value != 0 {
		t.Error("Unmarshaling null should result in a nullable int with an empty value and is null being true", ns)
	}
	testString = "5"
	err = ns.UnmarshalJSON([]byte(testString))
	if err != nil {
		t.Error("Failed to unmarshal \"Test\"", err)
	}
	if ns.HasValue != true || ns.Value != 5 {
		t.Error("Unmarshaling 1.2 should result in a nullable int with a value of 1.2 and is null being false", ns)
	}
	testString = "false"
	err = ns.UnmarshalJSON([]byte(testString))
	if err == nil {
		t.Error("expected an error", err)
	}
	if ns.HasValue != false || ns.Value != 0 {
		t.Error("Unmarshaling false should result in a nullable int with an empty value and is null being true", ns)
	}
}
