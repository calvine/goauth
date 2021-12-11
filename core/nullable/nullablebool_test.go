package nullable

import (
	"fmt"
	"testing"

	coreerrors "github.com/calvine/goauth/core/errors"
	"github.com/calvine/richerror/errors"
)

func TestNullableBoolGetPointerCopy(t *testing.T) {
	nb := NullableBool{}
	nb.Set(true)
	nbCopy := nb.GetPointerCopy()
	if *nbCopy != nb.Value {
		t.Error("\tnbCopy value should be the same as nb Value", nb, nbCopy)
	}
	if &nb.Value == nbCopy {
		t.Error("\tthe address of nb.Value and nbCopy should be different", &nb.Value, &nbCopy)
	}
	nb.Unset()
	nbCopy = nb.GetPointerCopy()
	if nbCopy != nil {
		t.Error("\tnbCopy should be nil because nb HasValue is false", nb, nbCopy)
	}
}

func TestNullableBoolSetUnset(t *testing.T) {
	nb := NullableBool{}
	testValue := true
	nb.Set(testValue)
	if nb.HasValue != true || !nb.Value {
		t.Error("\tnullable struct in invalid state after Set call", nb)
	}
	nb.Unset()
	if nb.HasValue || nb.Value != defaultBoolValue {
		t.Error("\tnullable struct in invalid state after Unset call", nb)
	}
}

func TestNullableBoolScan(t *testing.T) {
	ns := NullableBool{}
	err := ns.Scan(nil)
	if err != nil {
		t.Error("\tFailed to scan nil into NullableBool", err, ns)
	}
	if ns.Value != false || ns.HasValue != false {
		t.Error("\tNullable Bool has wrong value after scanning nil", ns)
	}
	testValue := true
	err = ns.Scan(testValue)
	if err != nil {
		t.Error("\tFailed to scan nil into NullableBool", err, ns)
	}
	if ns.Value != testValue || ns.HasValue != true {
		errMsg := fmt.Sprintf("Nullable Bool has wrong value after scanning %v", testValue)
		t.Error(errMsg, ns)
	}
	testNumber := 3
	err = ns.Scan(testNumber)
	if err != nil && err.(errors.RichError).GetErrorCode() != coreerrors.ErrCodeWrongType {
		t.Error("\tExpected error to be of type WrongTypeError", err)
	}
	if ns.Value != false || ns.HasValue != false {
		errMsg := fmt.Sprintf("Nullable Bool has wrong value after scanning %d", testNumber)
		t.Error(errMsg, ns)
	}
}

func TestNullableBoolMarshalJson(t *testing.T) {
	ns := NullableBool{
		Value:    false,
		HasValue: false,
	}
	data, err := ns.MarshalJSON()
	if err != nil {
		t.Error("\tFailed to marshal null to JSON.", err)
	}
	if string(data) != "null" {
		t.Error("\tdata from marshal was not null when underlaying nullable Bool was nil", data)
	}
	ns = NullableBool{
		Value:    false,
		HasValue: true,
	}
	data, err = ns.MarshalJSON()
	if err != nil {
		t.Error("\tFailed to marshal null to JSON.", err)
	}
	if string(data) != "false" {
		t.Error("\tdata from marshal was not what was expected", data, ns)
	}
}

func TestNullableBoolUnmarshalJson(t *testing.T) {
	testString := "null"
	ns := NullableBool{}
	err := ns.UnmarshalJSON([]byte(testString))
	if err != nil {
		t.Error("\tFailed to unmarshal null", err)
	}
	if ns.HasValue != false || ns.Value != false {
		t.Error("\tUnmarshaling null should result in a nullable Bool with an empty value and is null being true", ns)
	}
	testString = "true"
	err = ns.UnmarshalJSON([]byte(testString))
	if err != nil {
		t.Error("\tFailed to unmarshal \"Test\"", err)
	}
	if ns.HasValue != true || ns.Value != true {
		t.Error("\tUnmarshaling 1.2 should result in a nullable Bool with a value of 1.2 and is null being false", ns)
	}
	testString = "11"
	err = ns.UnmarshalJSON([]byte(testString))
	if err == nil {
		t.Error("\texpected an error", err)
	}
	if ns.HasValue != false || ns.Value != false {
		t.Error("\tUnmarshaling false should result in a nullable Bool with an empty value and is null being true", ns)
	}
}
