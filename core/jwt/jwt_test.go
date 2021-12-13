package jwt

import (
	"hash"
	"testing"

	coreerrors "github.com/calvine/goauth/core/errors"
	"github.com/calvine/goauth/internal/testutils"
)

func Test_splitEncodedJWT(t *testing.T) {
	type testCase struct {
		name              string
		encodedJWT        string
		expectedErrorCode string
	}
	testCases := []testCase{
		{
			name:       "GIVEN a valid encoded jwt EXPECT 3 parts to be returned",
			encodedJWT: "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxMjM0NTY3ODkwIiwibmFtZSI6IkpvaG4gRG9lIiwiaWF0IjoxNTE2MjM5MDIyfQ.SflKxwRJSMeKKF2QT4fwpMeJf36POk6yJV_adQssw5c",
		},
		{
			name:              "GIVEN an encoded JWT with no signature EXPECT error code jwt malformed",
			encodedJWT:        "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxMjM0NTY3ODkwIiwibmFtZSI6IkpvaG4gRG9lIiwiaWF0IjoxNTE2MjM5MDIyfQ",
			expectedErrorCode: coreerrors.ErrCodeJWTMalformed,
		},
		{
			name:              "GIVEN an encoded jwt with an empty signature EXPECT error code jwt missing signature",
			encodedJWT:        "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxMjM0NTY3ODkwIiwibmFtZSI6IkpvaG4gRG9lIiwiaWF0IjoxNTE2MjM5MDIyfQ.",
			expectedErrorCode: coreerrors.ErrCodeJWTSignatureMissing,
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			parts, err := splitEncodedJWT(tc.encodedJWT)
			if err != nil {
				testutils.HandleTestError(t, err, tc.expectedErrorCode)
			} else if tc.expectedErrorCode != "" {
				t.Errorf("\texpected an error to occurr: %s", tc.expectedErrorCode)
			} else {
				partsLength := len(parts)
				if partsLength != 3 {
					t.Errorf("\texpected 3 items in parts array but got %d parts", partsLength)
				}
			}
		})
	}
}

func TestDecodeHeader(t *testing.T) {
	type testCase struct {
		name              string
		encodedHeader     string
		expectedHeader    Header
		expectedErrorCode string
	}
	testCases := []testCase{}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			_, err := DecodeHeader(tc.encodedHeader)
			if err != nil {
				testutils.HandleTestError(t, err, tc.expectedErrorCode)
			} else if tc.expectedErrorCode != "" {
				t.Errorf("\texpected an error to occurr: %s", tc.expectedErrorCode)
			} else {

			}
		})
	}
	t.Error("\ttest not implemented yet...")
}

func TestHeaderEncode(t *testing.T) {
	type testCase struct {
		name                  string
		header                Header
		expectedEncodedHeader string
		expectedErrorCode     string
	}
	testCases := []testCase{}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			encodedHeader, err := tc.header.Encode()
			if err != nil {
				testutils.HandleTestError(t, err, tc.expectedErrorCode)
			} else if tc.expectedErrorCode != "" {
				t.Errorf("\texpected an error to occurr: %s", tc.expectedErrorCode)
			} else {
				if encodedHeader != tc.expectedEncodedHeader {
					t.Errorf("encodedHeader not expected: got - %s expected - %s", encodedHeader, tc.expectedEncodedHeader)
				}
			}
		})
	}
	t.Error("\ttest not implemented yet...")
}

func TestDecodeBody(t *testing.T) {
	type testCase struct {
		name              string
		encodedBody       string
		expectedErrorCode string
	}
	testCases := []testCase{}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			_, err := DecodeBody(tc.encodedBody)
			if err != nil {
				testutils.HandleTestError(t, err, tc.expectedErrorCode)
			} else if tc.expectedErrorCode != "" {
				t.Errorf("\texpected an error to occurr: %s", tc.expectedErrorCode)
			} else {

			}
		})
	}
	t.Error("\ttest not implemented yet...")
}

func TestStandardClaimEncode(t *testing.T) {
	type testCase struct {
		name                          string
		body                          StandardClaims
		expectedEncodedStandardClaims string
		expectedErrorCode             string
	}
	testCases := []testCase{}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			encodedStandardClaims, err := tc.body.Encode()
			if err != nil {
				testutils.HandleTestError(t, err, tc.expectedErrorCode)
			} else if tc.expectedErrorCode != "" {
				t.Errorf("\texpected an error to occurr: %s", tc.expectedErrorCode)
			} else {
				if encodedStandardClaims != tc.expectedEncodedStandardClaims {
					t.Errorf("encodedStandardClaims not expected: got - %s expected - %s", encodedStandardClaims, tc.expectedEncodedStandardClaims)
				}
			}
		})
	}
	t.Error("\ttest not implemented yet...")
}

func TestCalculateHMACSignature(t *testing.T) {
	type testCase struct {
		name                 string
		hashFunc             func() hash.Hash
		secret               string
		encodedHeaderAndBody string
		expectedSignature    string
	}
	testCases := []testCase{}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			signature := CalculateHMACSignature(tc.secret, tc.encodedHeaderAndBody, tc.hashFunc)
			if signature != tc.expectedSignature {
				t.Errorf("\tsignature is not expected value: got %s - expected %s", signature, tc.expectedSignature)
			}
		})
	}
	t.Error("\ttest not implemented yet...")
}

func TestBase64UrlEncode(t *testing.T) {
	type testCase struct {
		name           string
		inputString    string
		expectedOutput string
	}
	testCases := []testCase{}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			output := Base64UrlEncode([]byte(tc.inputString))
			if tc.expectedOutput != output {
				t.Errorf("\texpected output was incorrect: got - %s expected - %s", output, tc.expectedOutput)
			}
		})
	}
	t.Error("\ttest not implemented yet...")
}

func TestBase64UrlDecode(t *testing.T) {
	type testCase struct {
		name              string
		inputString       string
		expectedOutput    string
		expectedErrorCode string
	}
	testCases := []testCase{}
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
	t.Error("\ttest not implemented yet...")
}
