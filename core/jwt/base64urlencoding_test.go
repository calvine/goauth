package jwt

import (
	"testing"

	"github.com/calvine/goauth/internal/testutils"
)

func TestBase64UrlEncode(t *testing.T) {
	type testCase struct {
		name           string
		inputString    string
		expectedOutput string
	}
	testCases := []testCase{
		{
			name:           "GIVEN a string EXPECT the base64 url encoded string",
			inputString:    `{"alg":"HS256","typ":"JWT"}`,
			expectedOutput: "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9",
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			output := Base64UrlEncode([]byte(tc.inputString))
			if tc.expectedOutput != output {
				t.Errorf("\texpected output was incorrect: got - %s expected - %s", output, tc.expectedOutput)
			}
		})
	}
}

func TestBase64UrlDecode(t *testing.T) {
	type testCase struct {
		name              string
		inputString       string
		expectedOutput    string
		expectedErrorCode string
	}
	testCases := []testCase{
		{
			name:           "GIVEN a valid encoded string EXPECT the original string to be returned",
			inputString:    "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9",
			expectedOutput: `{"alg":"HS256","typ":"JWT"}`,
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			decoded, err := Base64UrlDecode(tc.inputString)
			if err != nil {
				testutils.HandleTestError(t, err, tc.expectedErrorCode)
			} else if tc.expectedErrorCode != "" {
				t.Errorf("\texpected an error to occurr: %s", tc.expectedErrorCode)
			} else {
				if string(decoded) != tc.expectedOutput {
					t.Errorf("\texpected output was incorrect: got - %s expected - %s", string(decoded), tc.expectedOutput)
				}
			}
		})
	}
}
