package nullable

import (
	"errors"
	"fmt"
	"testing"

	goautherrors "github.com/calvine/goauth/models/errors"
)

func TestNullableIntScan(t *testing.T) {
	ns := NullableInt{}
	err := ns.Scan(nil)
	if err != nil {
		t.Error("Failed to scan nil into NullableInt", err, ns)
	}
	if ns.Value != 0 || ns.IsNull != true {
		t.Error("Nullable int has wrong value after scanning nil", ns)
	}
	testValue := 2
	err = ns.Scan(testValue)
	if err != nil {
		t.Error("Failed to scan nil into NullableInt", err, ns)
	}
	if ns.Value != testValue || ns.IsNull != false {
		errMsg := fmt.Sprintf("Nullable int has wrong value after scanning %v", testValue)
		t.Error(errMsg, ns)
	}
	testString := "abc"
	err = ns.Scan(testString)
	emptyErr := &goautherrors.WrongTypeError{}
	if !errors.As(err, emptyErr) {
		t.Error("Expected error to be of type WrongTypeError", err)
	}
	if ns.Value != 0 || ns.IsNull != true {
		errMsg := fmt.Sprintf("Nullable int has wrong value after scanning %v", testString)
		t.Error(errMsg, ns)
	}
}

func TestNullableIntMarshalJson(t *testing.T) {
	ns := NullableInt{
		Value:  0,
		IsNull: true,
	}
	data, err := ns.MarshalJSON()
	if err != nil {
		t.Error("Failed to marshal null to JSON.", err)
	}
	if string(data) != "null" {
		t.Error("data from marshal was not null when underlaying nullable int was nil", data)
	}
	ns = NullableInt{
		Value:  -2,
		IsNull: false,
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
	if ns.IsNull != true || ns.Value != 0 {
		t.Error("Unmarshaling null should result in a nullable int with an empty value and is null being true", ns)
	}
	testString = "5"
	err = ns.UnmarshalJSON([]byte(testString))
	if err != nil {
		t.Error("Failed to unmarshal \"Test\"", err)
	}
	if ns.IsNull != false || ns.Value != 5 {
		t.Error("Unmarshaling 1.2 should result in a nullable int with a value of 1.2 and is null being false", ns)
	}
	testString = "false"
	err = ns.UnmarshalJSON([]byte(testString))
	if err == nil {
		t.Error("expected an error", err)
	}
	if ns.IsNull != true || ns.Value != 0 {
		t.Error("Unmarshaling false should result in a nullable int with an empty value and is null being true", ns)
	}
}
