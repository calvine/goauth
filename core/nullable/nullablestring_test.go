package nullable

import (
	"fmt"
	"testing"

	coreerrors "github.com/calvine/goauth/core/errors"
	"github.com/calvine/richerror/errors"
)

func TestNullableStringGetPointerCopy(t *testing.T) {
	ns := NullableString{}
	ns.Set("test string")
	nsCopy := ns.GetPointerCopy()
	if *nsCopy != ns.Value {
		t.Error("nsCopy value should be the same as ns Value", ns, nsCopy)
	}
	if &ns.Value == nsCopy {
		t.Error("the address of ns.Value and nsCopy should be different", &ns.Value, &nsCopy)
	}
	ns.Unset()
	nsCopy = ns.GetPointerCopy()
	if nsCopy != nil {
		t.Error("nsCopy should be nil because ns HasValue is false", ns, nsCopy)
	}
}

func TestNullableStringSetUnset(t *testing.T) {
	ns := NullableString{}
	testValue := "Hello Test"
	ns.Set(testValue)
	if ns.HasValue != true || ns.Value != testValue {
		t.Error("nullable struct in invalid state after Set call", ns)
	}
	ns.Unset()
	if ns.HasValue || ns.Value != defaultStringValue {
		t.Error("nullable struct in invalid state after Unset call", ns)
	}
}

func TestNullableStringScan(t *testing.T) {
	ns := NullableString{}
	err := ns.Scan(nil)
	if err != nil {
		t.Error("Failed to scan nil into NullableString", err, ns)
	}
	if ns.Value != "" || ns.HasValue != false {
		t.Error("Nullable string has wrong value after scanning nil", ns)
	}
	testString := "Test"
	err = ns.Scan(testString)
	if err != nil {
		t.Error("Failed to scan nil into NullableString", err, ns)
	}
	if ns.Value != testString || ns.HasValue != true {
		errMsg := fmt.Sprintf("Nullable string has wrong value after scanning %s", testString)
		t.Error(errMsg, ns)
	}
	testNumber := 3
	err = ns.Scan(testNumber)
	if err != nil && err.(errors.RichError).GetErrorCode() != coreerrors.ErrCodeWrongType {
		t.Error("Expected error to be of type WrongTypeError", err)
	}
	if ns.Value != "" || ns.HasValue != false {
		errMsg := fmt.Sprintf("Nullable string has wrong value after scanning %d", testNumber)
		t.Error(errMsg, ns)
	}
}

func TestNullableStringMarshalJson(t *testing.T) {
	ns := NullableString{
		Value:    "",
		HasValue: false,
	}
	data, err := ns.MarshalJSON()
	if err != nil {
		t.Error("Failed to marshal null to JSON.", err)
	}
	if string(data) != "null" {
		t.Error("data from marshal was not null when underlaying nullable string was nil", data)
	}
	ns = NullableString{
		Value:    "Test",
		HasValue: true,
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
	if ns.HasValue != false || ns.Value != "" {
		t.Error("Unmarshaling null should result in a nullable string with an empty value and is null being true", ns)
	}
	testString = "\"Test\""
	err = ns.UnmarshalJSON([]byte(testString))
	if err != nil {
		t.Error("Failed to unmarshal \"Test\"", err)
	}
	if ns.HasValue != true || ns.Value != "Test" {
		t.Error("Unmarshaling \"Test\" should result in a nullable string with a value of Test and is null being false", ns)
	}
	testString = "3"
	err = ns.UnmarshalJSON([]byte(testString))
	if err == nil {
		t.Error("expected a WrongTypeError", err)
	}
	if ns.HasValue != false || ns.Value != "" {
		t.Error("Unmarshaling 3 should result in a nullable string with an empty value and is null being true", ns)
	}
}
