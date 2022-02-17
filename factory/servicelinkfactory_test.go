package factory

import (
	"fmt"
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
		queryParams       map[string]string
		expectedFullURL   string
		expectedErrorCode string
	}
	testCases := []testCase{
		{
			name:            "GIVEN a call without query params EXPECT success",
			linkPath:        "/path/to/resource",
			expectedFullURL: fmt.Sprintf("%s/path/to/resource", testPublicBaseURL),
		},
		{
			name:     "GIVEN a call with query params EXPECT success",
			linkPath: "/path/to/resource",
			queryParams: map[string]string{
				"v": "1234",
			},
			expectedFullURL: fmt.Sprintf("%s/path/to/resource?v=1234", testPublicBaseURL),
		},
	}
	serviceLinkFactory, err := NewServiceLinkFactory(testPublicBaseURL)
	if err != nil {
		t.Errorf("\tfailed to create test serviceLinkFactory with error: %s", err.Error())
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			url, err := serviceLinkFactory.CreateLink(tc.linkPath, tc.queryParams)
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
	testCases := []testCase{
		{
			name:               "GIVEN EXPECT success",
			passwordResetToken: "1234",
			expectedFullURL:    fmt.Sprintf("%s/user/resetpassword/1234", testPublicBaseURL),
		},
	}
	serviceLinkFactory, err := NewServiceLinkFactory(testPublicBaseURL)
	if err != nil {
		t.Errorf("\tfailed to create test serviceLinkFactory with error: %s", err.Error())
	}
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
	testCases := []testCase{
		{
			name:                "GIVEN EXPECT success",
			confirmContactToken: "1234",
			expectedFullURL:     fmt.Sprintf("%s/user/confirmcontact/1234", testPublicBaseURL),
		},
	}
	serviceLinkFactory, err := NewServiceLinkFactory(testPublicBaseURL)
	if err != nil {
		t.Errorf("\tfailed to create test serviceLinkFactory with error: %s", err.Error())
	}
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

func TestCreateLoginLink(t *testing.T) {
	type testCase struct {
		name              string
		expectedFullURL   string
		expectedErrorCode string
	}
	testCases := []testCase{
		{
			name:            "GIVEN EXPECT success",
			expectedFullURL: fmt.Sprintf("%s/auth/login", testPublicBaseURL),
		},
	}
	serviceLinkFactory, err := NewServiceLinkFactory(testPublicBaseURL)
	if err != nil {
		t.Errorf("\tfailed to create test serviceLinkFactory with error: %s", err.Error())
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			url, err := serviceLinkFactory.CreateLoginLink()
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

func TestCreateMagicLoginLink(t *testing.T) {
	type testCase struct {
		name              string
		magicLoginToken   string
		expectedFullURL   string
		expectedErrorCode string
	}
	testCases := []testCase{
		{
			name:            "GIVEN EXPECT success",
			magicLoginToken: "1234",
			expectedFullURL: fmt.Sprintf("%s/auth/magiclogin?m=1234", testPublicBaseURL),
		},
	}
	serviceLinkFactory, err := NewServiceLinkFactory(testPublicBaseURL)
	if err != nil {
		t.Errorf("\tfailed to create test serviceLinkFactory with error: %s", err.Error())
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			url, err := serviceLinkFactory.CreateMagicLoginLink(tc.magicLoginToken)
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

func TestCreateUserRegisterLink(t *testing.T) {
	type testCase struct {
		name              string
		expectedFullURL   string
		expectedErrorCode string
	}
	testCases := []testCase{
		{
			name:            "GIVEN EXPECT success",
			expectedFullURL: fmt.Sprintf("%s/user/register", testPublicBaseURL),
		},
	}
	serviceLinkFactory, err := NewServiceLinkFactory(testPublicBaseURL)
	if err != nil {
		t.Errorf("\tfailed to create test serviceLinkFactory with error: %s", err.Error())
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			url, err := serviceLinkFactory.CreateUserRegisterLink()
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
