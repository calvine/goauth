package factory

import (
	"testing"

	"github.com/calvine/goauth/internal/testutils"
)

/*
	Some quick thoughts on these tests. They serve a few purposes:

	1. Ensure the functions work properly
	2. Hard code expected paths (for now) of what links should look like so that
*/

const (
	testPublicBaseURL = "http://localhost:8080"
)

func TestCreateLink(t *testing.T) {
	type testCase struct {
		name              string
		linkPath          string
		expectedFullURL   string
		expectedErrorCode string
	}
	testCases := []testCase{}
	serviceLinkFactory, err := NewServiceLinkFactory(testPublicBaseURL)
	if err != nil {
		t.Errorf("\tfailed to create test serviceLinkFactory with error: %s", err.Error())
	}
	t.Error("Test not impolemented")
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			url, err := serviceLinkFactory.CreateLink(tc.linkPath)
			if err != nil {
				testutils.HandleTestError(t, err, tc.expectedErrorCode)
			} else if tc.expectedErrorCode != "" {
				t.Errorf("\texpected an error to occurr: %s", tc.expectedErrorCode)
			} else {
				if url != tc.expectedFullURL {
					t.Errorf("\tvalue of url was not expected: got - %s expected - %s", url, tc.expectedFullURL)
				}
			}
		})
	}
}

func TestCreatePasswordResetLink(t *testing.T) {
	type testCase struct {
		name               string
		passwordResetToken string
		expectedFullURL    string
		expectedErrorCode  string
	}
	testCases := []testCase{}
	serviceLinkFactory, err := NewServiceLinkFactory(testPublicBaseURL)
	if err != nil {
		t.Errorf("\tfailed to create test serviceLinkFactory with error: %s", err.Error())
	}
	t.Error("Test not impolemented")
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			url, err := serviceLinkFactory.CreatePasswordResetLink(tc.passwordResetToken)
			if err != nil {
				testutils.HandleTestError(t, err, tc.expectedErrorCode)
			} else if tc.expectedErrorCode != "" {
				t.Errorf("\texpected an error to occurr: %s", tc.expectedErrorCode)
			} else {
				if url != tc.expectedFullURL {
					t.Errorf("\tvalue of url was not expected: got - %s expected - %s", url, tc.expectedFullURL)
				}
			}
		})
	}
}

func TestCreateConfirmContactLink(t *testing.T) {
	type testCase struct {
		name                string
		confirmContactToken string
		expectedFullURL     string
		expectedErrorCode   string
	}
	testCases := []testCase{}
	serviceLinkFactory, err := NewServiceLinkFactory(testPublicBaseURL)
	if err != nil {
		t.Errorf("\tfailed to create test serviceLinkFactory with error: %s", err.Error())
	}
	t.Error("Test not impolemented")
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			url, err := serviceLinkFactory.CreateConfirmContactLink(tc.confirmContactToken)
			if err != nil {
				testutils.HandleTestError(t, err, tc.expectedErrorCode)
			} else if tc.expectedErrorCode != "" {
				t.Errorf("\texpected an error to occurr: %s", tc.expectedErrorCode)
			} else {
				if url != tc.expectedFullURL {
					t.Errorf("\tvalue of url was not expected: got - %s expected - %s", url, tc.expectedFullURL)
				}
			}
		})
	}
}
