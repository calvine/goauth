package testutilities

import (
	"context"
	"fmt"
	"reflect"
	"testing"

	"github.com/calvine/goauth/core/models"
	"github.com/calvine/goauth/core/services"
	"github.com/calvine/richerror/errors"
)

const (
	reasonTypesDoNotMatch    = "types do not match"
	reasonKindsDoNotMatch    = "kinds do not match"
	reasonDifferentLengths   = "lengths of items not equal"
	reasonFieldOrderMismatch = "fields not in the same order"
	reasonDifferentNumFields = "number of fields do not match"
	reasonValuesDoNotMatch   = "values do not match"
	// reasonKindNotSupported    = "the kind provided is not supported in equality check: %s"
	reasonValuesNotComparable = "the values provided are not of comparable types: type = %s - kind = %s"
)

func PerformErrorCheck(t *testing.T, testCase BaseTestCase, err errors.RichError) {
	if err != nil {
		if !testCase.ExpectedError {
			t.Log(err.Error())
			t.Fatalf("unexpeced error occurred: %s", err.GetErrorCode())
		} else if testCase.ExpectedErrorCode != err.GetErrorCode() {
			t.Log(err.Error())
			t.Fatalf("expeced error code does not match expected error type: got: %s - expected %s", err.GetErrorCode(), testCase.ExpectedErrorCode)
		}
	}
}

func ValidateExpectedAppEqualToStoredAppWithAppService(t *testing.T, appService services.AppService, expectedApp models.App) {
	app, err := appService.GetAppByID(context.TODO(), expectedApp.ID, "ValidateExpectedAppEqualToStoredAppWithAppService")
	if err != nil {
		t.Fatalf("failed to get app with id %s for comparison from underlying data source: %s", expectedApp.ID, err.Error())
	}
	equalityCheck := Equals(app, expectedApp)
	if !equalityCheck.AreEqual {
		t.Fatalf("app not equal to expected app: %v", equalityCheck.Failures)
	}
}

func ValidateExpectedScopeEqualToStoredScopeWithAppService(t *testing.T, appService services.AppService, expectedScope models.Scope) {
	scope, err := appService.GetScopeByID(context.TODO(), expectedScope.ID, "ValidateExpectedScopeEqualToStoredScopeWithAppService")
	if err != nil {
		t.Log(err.Error())
		t.Fatalf("failed to get scope with id %s for comparison from underlying data source: %s", expectedScope.ID, err.GetErrorCode())
	}
	equalityCheck := Equals(scope, expectedScope)
	if !equalityCheck.AreEqual {
		t.Fatalf("scope not equal to expected scope: %v", equalityCheck.Failures)
	}
}

type equalityCheckResult struct {
	AreEqual bool
	Failures []equlaityCheckFailure
}

type equlaityCheckFailure struct {
	FailurePath string
	Reason      string
	Value1      interface{}
	Value2      interface{}
}

func (ecf equlaityCheckFailure) ToString(includeValues bool) string {
	var failurePath string
	if ecf.FailurePath == "" {
		failurePath = "value"
	}
	if includeValues {
		return fmt.Sprintf("inequality found at %s - %s: Values: %v - %v", ecf.FailurePath, ecf.Reason, ecf.Value1, ecf.Value2)
	}
	return fmt.Sprintf("inequality found at %s - %s", failurePath, ecf.Reason)
}

// TODO: have options for reference checks for pointer maps arrays slices.

// Equals is a deep equals check against two things passed in as parameters
func Equals(thing1, thing2 interface{}) equalityCheckResult {
	result := equalityCheckResult{
		Failures: make([]equlaityCheckFailure, 0),
	}
	innerEquals(thing1, thing2, "", &result)
	result.AreEqual = len(result.Failures) == 0
	return result
}

func innerEquals(thing1, thing2 interface{}, currentPath string, result *equalityCheckResult) {
	t1Type := reflect.TypeOf(thing1)
	t2Type := reflect.TypeOf(thing2)

	if t1Type != t2Type {
		reason := equlaityCheckFailure{
			FailurePath: currentPath,
			Reason:      reasonTypesDoNotMatch,
		}
		result.Failures = append(result.Failures, reason)
		return
	}

	if thing1 == nil && thing2 == nil {
		result.AreEqual = true
		return
	}

	t1Kind := t1Type.Kind()
	t2Kind := t2Type.Kind()

	if t1Kind != t2Kind {
		reason := equlaityCheckFailure{
			FailurePath: currentPath,
			Reason:      reasonKindsDoNotMatch,
		}
		result.Failures = append(result.Failures, reason)
		return
	}

	t1Value := reflect.ValueOf(thing1)
	// t1IsNil := t1Value.IsNil()
	t2Value := reflect.ValueOf(thing2)
	// t2IsNil := t2Value.IsNil()

	switch t1Kind {
	case reflect.Struct:
		t1NumFields := t1Type.NumField()
		t2NumFields := t2Type.NumField()
		if t1NumFields != t2NumFields {
			// structs do not have the same number of fields, so they cannot be equal
			reason := equlaityCheckFailure{
				FailurePath: currentPath,
				Reason:      reasonDifferentNumFields,
			}
			result.Failures = append(result.Failures, reason)
			return
		}
		for i := 0; i < t1NumFields; i++ {
			t1FieldValue := t1Value.FieldByIndex([]int{i})
			t2FieldValue := t2Value.FieldByIndex([]int{i})
			t1Field := t1Type.Field(i)
			t2Field := t2Type.Field(i)
			path := fmt.Sprintf("%s.%s", currentPath, t1Field.Name)
			if t1Field.Name != t2Field.Name {
				// struct fields misaligned some how? Not sure this is possible given checks above...
				reason := equlaityCheckFailure{
					FailurePath: path,
					Reason:      reasonFieldOrderMismatch,
				}
				result.Failures = append(result.Failures, reason)
			} else if t1FieldValue.CanInterface() && t2FieldValue.CanInterface() { // TODO: is this the best way to tell if a struct field is exported?
				// if the fields we are looking at are not exported we cannot interface them to check...
				innerEquals(t1FieldValue.Interface(), t2FieldValue.Interface(), path, result)
			}
		}
	case reflect.Ptr:
		dereferencedT1 := t1Value.Elem().Interface()
		dereferencedT2 := t2Value.Elem().Interface()
		innerEquals(dereferencedT1, dereferencedT2, currentPath, result)
	case reflect.Map:
		// iterate over children and evaluate values and keys exist in both maps
		t1Len := t1Value.Len()
		t2Len := t2Value.Len()
		if t1Len != t2Len {
			// lengths are not equal so the maps are not equal
			reason := equlaityCheckFailure{
				FailurePath: currentPath,
				Reason:      reasonDifferentNumFields,
			}
			result.Failures = append(result.Failures, reason)
			return
		}
		mapIter := t1Value.MapRange()
		for mapIter.Next() {
			k := mapIter.Key()
			v := mapIter.Value()

			v2 := t2Value.MapIndex(k)
			path := fmt.Sprintf("%s[%v]", currentPath, k)
			innerEquals(v.Interface(), v2.Interface(), path, result)
		}
	case reflect.Array:
	case reflect.Slice:
		// for now this will do an ordinal match of children, possibly allow a param to allow for non ordinal match?
		// iterate over array and check children
		t1Len := t1Value.Len()
		t2Len := t2Value.Len()
		if t1Len != t2Len {
			// lengths are not equal so the arrays / slices are not equal
			reason := equlaityCheckFailure{
				FailurePath: currentPath,
				Reason:      reasonDifferentLengths,
			}
			result.Failures = append(result.Failures, reason)
			return
		}
		for i := 0; i < t1Len; i++ {
			t1Item := t1Value.Index(i).Interface()
			t2Item := t2Value.Index(i).Interface()
			path := fmt.Sprintf("%s[%d]", currentPath, i)
			innerEquals(t1Item, t2Item, path, result)
		}
	case reflect.Interface:
		// interested to see how this might work?
	case reflect.Chan:
	case reflect.Func:
	case reflect.Invalid:
	case reflect.UnsafePointer:
		// I dont think we should support these?
		break
	default:
		// should be basic type thats is comparable?
		if !t1Type.Comparable() || !t2Type.Comparable() {
			// if these are not comparable, thats a problem and we need to add a case to check for it...
			reason := equlaityCheckFailure{
				FailurePath: currentPath,
				Reason:      fmt.Sprintf(reasonValuesNotComparable, t1Type.Name(), t1Kind.String()),
			}
			result.Failures = append(result.Failures, reason)
			return
		}
		if t1Value.Interface() != t2Value.Interface() {
			reason := equlaityCheckFailure{
				FailurePath: currentPath,
				Reason:      reasonValuesDoNotMatch,
				Value1:      t1Value.Interface(),
				Value2:      t2Value.Interface(),
			}
			result.Failures = append(result.Failures, reason)
		}
	}
}
