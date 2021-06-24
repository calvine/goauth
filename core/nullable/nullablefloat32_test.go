package nullable

import (
	"errors"
	"fmt"
	"testing"

	goautherrors "github.com/calvine/goauth/core/errors"
)

func TestNullableFloat32SetUnset(t *testing.T) {
	nf := NullableFloat32{}
	testValue := float32(1.23)
	nf.Set(testValue)
	if nf.IsNull != false || nf.Value != testValue {
		t.Error("nullable struct in invalid state after Set call", nf)
	}
	nf.Unset()
	if !nf.IsNull || nf.Value != defaultFloat32Value {
		t.Error("nullable struct in invalid state after Unset call", nf)
	}
}

func TestNullableFloat32Scan(t *testing.T) {
	ns := NullableFloat32{}
	err := ns.Scan(nil)
	if err != nil {
		t.Error("Failed to scan nil into NullableFloat32", err, ns)
	}
	if ns.Value != 0 || ns.IsNull != true {
		t.Error("Nullable float32 has wrong value after scanning nil", ns)
	}
	testValue := float32(1.2)
	err = ns.Scan(testValue)
	if err != nil {
		t.Error("Failed to scan nil into NullableFloat32", err, ns)
	}
	if ns.Value != testValue || ns.IsNull != false {
		errMsg := fmt.Sprintf("Nullable float32 has wrong value after scanning %f", testValue)
		t.Error(errMsg, ns)
	}
	testNumber := 3
	err = ns.Scan(testNumber)
	emptyErr := &goautherrors.WrongTypeError{}
	if !errors.As(err, emptyErr) {
		t.Error("Expected error to be of type WrongTypeError", err)
	}
	if ns.Value != 0 || ns.IsNull != true {
		errMsg := fmt.Sprintf("Nullable float32 has wrong value after scanning %d", testNumber)
		t.Error(errMsg, ns)
	}
}

func TestNullableFloat32MarshalJson(t *testing.T) {
	ns := NullableFloat32{
		Value:  0,
		IsNull: true,
	}
	data, err := ns.MarshalJSON()
	if err != nil {
		t.Error("Failed to marshal null to JSON.", err)
	}
	if string(data) != "null" {
		t.Error("data from marshal was not null when underlaying nullable float32 was nil", data)
	}
	ns = NullableFloat32{
		Value:  1.2,
		IsNull: false,
	}
	data, err = ns.MarshalJSON()
	if err != nil {
		t.Error("Failed to marshal null to JSON.", err)
	}
	if string(data) != "1.2" {
		t.Error("data from marshal was not what was expected", data, ns)
	}
}

func TestNullableFloat32UnmarshalJson(t *testing.T) {
	testString := "null"
	ns := NullableFloat32{}
	err := ns.UnmarshalJSON([]byte(testString))
	if err != nil {
		t.Error("Failed to unmarshal null", err)
	}
	if ns.IsNull != true || ns.Value != 0 {
		t.Error("Unmarshaling null should result in a nullable float32 with an empty value and is null being true", ns)
	}
	testString = "1.2"
	err = ns.UnmarshalJSON([]byte(testString))
	if err != nil {
		t.Error("Failed to unmarshal \"Test\"", err)
	}
	if ns.IsNull != false || ns.Value != 1.2 {
		t.Error("Unmarshaling 1.2 should result in a nullable float32 with a value of 1.2 and is null being false", ns)
	}
	testString = "false"
	err = ns.UnmarshalJSON([]byte(testString))
	if err == nil {
		t.Error("expected an error", err)
	}
	if ns.IsNull != true || ns.Value != 0 {
		t.Error("Unmarshaling false should result in a nullable float32 with an empty value and is null being true", ns)
	}
}
