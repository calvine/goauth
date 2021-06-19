package models

import (
	"errors"
	"fmt"
	"testing"

	goautherrors "github.com/calvine/goauth/models/errors"
)

func TestNullableBoolScan(t *testing.T) {
	ns := NullableBool{}
	err := ns.Scan(nil)
	if err != nil {
		t.Error("Failed to scan nil into NullableBool", err, ns)
	}
	if ns.Value != false || ns.IsNull != true {
		t.Error("Nullable Bool has wrong value after scanning nil", ns)
	}
	testValue := true
	err = ns.Scan(testValue)
	if err != nil {
		t.Error("Failed to scan nil into NullableBool", err, ns)
	}
	if ns.Value != testValue || ns.IsNull != false {
		errMsg := fmt.Sprintf("Nullable Bool has wrong value after scanning %v", testValue)
		t.Error(errMsg, ns)
	}
	testNumber := 3
	err = ns.Scan(testNumber)
	emptyErr := &goautherrors.WrongTypeError{}
	if !errors.As(err, emptyErr) {
		t.Error("Expected error to be of type WrongTypeError", err)
	}
	if ns.Value != false || ns.IsNull != true {
		errMsg := fmt.Sprintf("Nullable Bool has wrong value after scanning %d", testNumber)
		t.Error(errMsg, ns)
	}
}

func TestNullableBoolMarshalJson(t *testing.T) {
	ns := NullableBool{
		Value:  false,
		IsNull: true,
	}
	data, err := ns.MarshalJSON()
	if err != nil {
		t.Error("Failed to marshal null to JSON.", err)
	}
	if string(data) != "null" {
		t.Error("data from marshal was not null when underlaying nullable Bool was nil", data)
	}
	ns = NullableBool{
		Value:  false,
		IsNull: false,
	}
	data, err = ns.MarshalJSON()
	if err != nil {
		t.Error("Failed to marshal null to JSON.", err)
	}
	if string(data) != "false" {
		t.Error("data from marshal was not what was expected", data, ns)
	}
}

func TestNullableBoolUnmarshalJson(t *testing.T) {
	testString := "null"
	ns := NullableBool{}
	err := ns.UnmarshalJSON([]byte(testString))
	if err != nil {
		t.Error("Failed to unmarshal null", err)
	}
	if ns.IsNull != true || ns.Value != false {
		t.Error("Unmarshaling null should result in a nullable Bool with an empty value and is null being true", ns)
	}
	testString = "true"
	err = ns.UnmarshalJSON([]byte(testString))
	if err != nil {
		t.Error("Failed to unmarshal \"Test\"", err)
	}
	if ns.IsNull != false || ns.Value != true {
		t.Error("Unmarshaling 1.2 should result in a nullable Bool with a value of 1.2 and is null being false", ns)
	}
	testString = "11"
	err = ns.UnmarshalJSON([]byte(testString))
	if err == nil {
		t.Error("expected an error", err)
	}
	if ns.IsNull != true || ns.Value != false {
		t.Error("Unmarshaling false should result in a nullable Bool with an empty value and is null being true", ns)
	}
}
