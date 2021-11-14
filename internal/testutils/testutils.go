package testutils

import (
	"testing"

	"github.com/calvine/richerror/errors"
)

// HandleTestError runs common error checks for tests.
func HandleTestError(t *testing.T, err errors.RichError, expectedErrorCode string) {
	if expectedErrorCode == "" {
		t.Errorf("\tunexpected error encountered: %s - %s", err.GetErrorCode(), err.Error())
		t.Fail()
	} else if expectedErrorCode != err.GetErrorCode() {
		t.Errorf("\terror code did not match expected: got - %s expected - %s", err.GetErrorCode(), expectedErrorCode)
		t.Fail()
	}
}
