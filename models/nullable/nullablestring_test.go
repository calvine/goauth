package nullable

import (
	"errors"
	"fmt"
	"testing"

	goautherrors "github.com/calvine/goauth/models/errors"
)

func TestNullableStringSetUnset(t *testing.T) {
	ns := NullableString{}
	testValue := "Hello Test"
	ns.Set(testValue)
	if ns.IsNull != false || ns.Value != testValue {
		t.Error("nullable struct in invalid state after Set call", ns)
	}
	ns.Unset()
	if !ns.IsNull || ns.Value != defaultStringValue {
		t.Error("nullable struct in invalid state after Unset call", ns)
	}
}

func TestNullableStringScan(t *testing.T) {
	ns := NullableString{}
	err := ns.Scan(nil)
	if err != nil {
		t.Error("Failed to scan nil into NullableString", err, ns)
	}
	if ns.Value != "" || ns.IsNull != true {
		t.Error("Nullable string has wrong value after scanning nil", ns)
	}
	testString := "Test"
	err = ns.Scan(testString)
	if err != nil {
		t.Error("Failed to scan nil into NullableString", err, ns)
	}
	if ns.Value != testString || ns.IsNull != false {
		errMsg := fmt.Sprintf("Nullable string has wrong value after scanning %s", testString)
		t.Error(errMsg, ns)
	}
	testNumber := 3
	err = ns.Scan(testNumber)
	emptyErr := &goautherrors.WrongTypeError{}
	if !errors.As(err, emptyErr) {
		t.Error("Expected error to be of type WrongTypeError", err)
	}
	if ns.Value != "" || ns.IsNull != true {
		errMsg := fmt.Sprintf("Nullable string has wrong value after scanning %d", testNumber)
		t.Error(errMsg, ns)
	}
}

func TestNullableStringMarshalJson(t *testing.T) {
	ns := NullableString{
		Value:  "",
		IsNull: true,
	}
	data, err := ns.MarshalJSON()
	if err != nil {
		t.Error("Failed to marshal null to JSON.", err)
	}
	if string(data) != "null" {
		t.Error("data from marshal was not null when underlaying nullable string was nil", data)
	}
	ns = NullableString{
		Value:  "Test",
		IsNull: false,
	}
	data, err = ns.MarshalJSON()
	if err != nil {
		t.Error("Failed to marshal null to JSON.", err)
	}
	if string(data) != "\"Test\"" {
		t.Error("data from marshal was not what was expected", data, ns)
	}
}

func TestNullableStringUnmarshalJson(t *testing.T) {
	testString := "null"
	ns := NullableString{}
	err := ns.UnmarshalJSON([]byte(testString))
	if err != nil {
		t.Error("Failed to unmarshal null", err)
	}
	if ns.IsNull != true || ns.Value != "" {
		t.Error("Unmarshaling null should result in a nullable string with an empty value and is null being true", ns)
	}
	testString = "\"Test\""
	err = ns.UnmarshalJSON([]byte(testString))
	if err != nil {
		t.Error("Failed to unmarshal \"Test\"", err)
	}
	if ns.IsNull != false || ns.Value != "Test" {
		t.Error("Unmarshaling \"Test\" should result in a nullable string with a value of Test and is null being false", ns)
	}
	testString = "3"
	err = ns.UnmarshalJSON([]byte(testString))
	if err == nil {
		t.Error("expected a WrongTypeError", err)
	}
	if ns.IsNull != true || ns.Value != "" {
		t.Error("Unmarshaling 3 should result in a nullable string with an empty value and is null being true", ns)
	}
}
