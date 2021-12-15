package jwt

import (
	"crypto/sha256"
	"crypto/sha512"
	"hash"
	"testing"
	"time"

	coreerrors "github.com/calvine/goauth/core/errors"
	"github.com/calvine/goauth/internal/testutils"
)

func TestSplitEncodedJWT(t *testing.T) {
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
			parts, err := SplitEncodedJWT(tc.encodedJWT)
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
	testCases := []testCase{
		{
			name:          "GIVEN a valid encoded jwt header EXPECT success",
			encodedHeader: "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9",
			expectedHeader: Header{
				Algorithm: Alg_HS256,
				TokenType: Type_JWT,
			},
		},
		{
			name:          "GIVEN an encoded jwt header with invalid json EXPECT error code jwt malformed",
			encodedHeader: "eyJhbGdiOiJIUzI1NiIsInR5cCI6IkpXVCJ9",
			expectedHeader: Header{
				Algorithm: Alg_HS256,
				TokenType: Type_JWT,
			},
			expectedErrorCode: coreerrors.ErrCodeJWTMalformed,
		},
		{
			name:          "GIVEN an invalid base 64 encoded string EXPECT error code jwt malformed",
			encodedHeader: "eyJhbGdiOiJIUzI1NiIsQQQR5cCI6IkpXVCJ9",
			expectedHeader: Header{
				Algorithm: Alg_HS256,
				TokenType: Type_JWT,
			},
			expectedErrorCode: coreerrors.ErrCodeJWTMalformed,
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			decodedHeader, err := DecodeHeader(tc.encodedHeader)
			if err != nil {
				testutils.HandleTestError(t, err, tc.expectedErrorCode)
			} else if tc.expectedErrorCode != "" {
				t.Errorf("\texpected an error to occurr: %s", tc.expectedErrorCode)
			} else {
				if decodedHeader != tc.expectedHeader {
					t.Errorf("\tdecodedHeader value was not expected: got - %v expected - %v", decodedHeader, tc.expectedHeader)
				}
			}
		})
	}
}

func TestHeaderEncode(t *testing.T) {
	type testCase struct {
		name                  string
		header                Header
		expectedEncodedHeader string
		expectedErrorCode     string
	}
	testCases := []testCase{
		{
			name: "GIVEN a valid header EXPECT a properly encoded header",
			header: Header{
				Algorithm: Alg_HS256,
				TokenType: Type_JWT,
			},
			expectedEncodedHeader: "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9",
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			encodedHeader, err := tc.header.Encode()
			if err != nil {
				testutils.HandleTestError(t, err, tc.expectedErrorCode)
			} else if tc.expectedErrorCode != "" {
				t.Errorf("\texpected an error to occurr: %s", tc.expectedErrorCode)
			} else {
				if encodedHeader != tc.expectedEncodedHeader {
					t.Errorf("\tencodedHeader not expected: got - %s expected - %s", encodedHeader, tc.expectedEncodedHeader)
				}
			}
		})
	}
}

func TestDecodeStandardClaims(t *testing.T) {
	type testCase struct {
		name                   string
		encodedStandardClaims  string
		expectedStandardClaims StandardClaims
		expectedErrorCode      string
	}
	testCases := []testCase{
		{
			name:                  "GIVEN valid base 64 encoded jwt standard claims EXPECT success",
			encodedStandardClaims: "eyJpc3MiOiJnb2F1dGgiLCJzdWIiOiIxMjM0NTY3ODkwIiwiYXVkIjpbImdvYXV0aCJdLCJleHAiOjE1MTYyMzkwMjIsIm5iZiI6MTUxNjIzOTAyMiwiaWF0IjoxNTE2MjM5MDIyLCJqdGkiOiI3ODk0NTYrMTIzMDMyMTM2NTQ5ODQ5Njg1NDErNDQ1NTUyIn0",
			expectedStandardClaims: StandardClaims{
				Issuer:         "goauth",
				Subject:        "1234567890",
				Audience:       "goauth",
				ExpirationTime: Time(time.Unix(1516239022, 0)),
				NotBefore:      Time(time.Unix(1516239022, 0)),
				IssuedAt:       Time(time.Unix(1516239022, 0)),
				JWTID:          "789456+12303213654984968541+445552",
			},
		},
		{
			name:                   "GIVEN base64 encoded jwt standard claims with invalid json EXPECT error code jwt malformed",
			encodedStandardClaims:  "eyJpc3MiOiJnb2F1dGgiLCJzdWIiOiIxMjM0NTY3ODkwIiwiYXVkIjpbImdvYXV0aCJdLCJleHAiOjE1MTYyMzkwMjIsIm5iZiI6MTUxNjIzOTAyMiwiaWF0IjoxNTE2MjM5MDIyLCJqdGkiOiI3ODk0NTYrMTIzMDMyMTM2NTQ5ODQ5Njg1NDErNDQ1NTUyQQQ",
			expectedStandardClaims: StandardClaims{},
			expectedErrorCode:      coreerrors.ErrCodeJWTMalformed,
		},
		{
			name:                   "GIVEN an invalid base 64 encoded string EXPECT error code jwt malformed",
			encodedStandardClaims:  "eyJpc3MiOiJZZZIF1~!@#dGgiLCJzdWIiOiIxMjM0NTY3ODkwIiwiYXVkIjpbImdvYXV0aCJdLCJleHAiOjE1MTYyMzkwMjIsIm5iZiI6MTUxNjIzOTAyMiwiaWF0IjoxNTE2MjM5MDIyLCJqdGkiOiI3ODk0NTYrMTIzMDMyMTM2NTQ5ODQ5Njg1NDErNDQ1NTUyQQQ",
			expectedStandardClaims: StandardClaims{},
			expectedErrorCode:      coreerrors.ErrCodeJWTMalformed,
		}}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			decodedStandardClaims, err := DecodeStandardClaims(tc.encodedStandardClaims)
			if err != nil {
				testutils.HandleTestError(t, err, tc.expectedErrorCode)
			} else if tc.expectedErrorCode != "" {
				t.Errorf("\texpected an error to occurr: %s", tc.expectedErrorCode)
			} else {
				if decodedStandardClaims != tc.expectedStandardClaims {
					t.Errorf("decodedStandardClaims value is not expected: got - %v expected - %v", decodedStandardClaims, tc.expectedStandardClaims)
				}
			}
		})
	}
}

func TestStandardClaimEncode(t *testing.T) {
	type testCase struct {
		name                          string
		standardClaims                StandardClaims
		expectedEncodedStandardClaims string
		expectedErrorCode             string
	}
	testCases := []testCase{
		{
			name: "GIVEN a set of standard claims EXPECT properly encoded standard claims",
			standardClaims: StandardClaims{
				Issuer:         "goauth",
				Subject:        "1234567890",
				Audience:       "goauth",
				ExpirationTime: Time(time.Unix(1516239022, 0)),
				NotBefore:      Time(time.Unix(1516239022, 0)),
				IssuedAt:       Time(time.Unix(1516239022, 0)),
				JWTID:          "789456+12303213654984968541+445552",
			},
			expectedEncodedStandardClaims: "eyJpc3MiOiJnb2F1dGgiLCJzdWIiOiIxMjM0NTY3ODkwIiwiYXVkIjpbImdvYXV0aCJdLCJleHAiOjE1MTYyMzkwMjIsIm5iZiI6MTUxNjIzOTAyMiwiaWF0IjoxNTE2MjM5MDIyLCJqdGkiOiI3ODk0NTYrMTIzMDMyMTM2NTQ5ODQ5Njg1NDErNDQ1NTUyIn0",
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			encodedStandardClaims, err := tc.standardClaims.Encode()
			if err != nil {
				testutils.HandleTestError(t, err, tc.expectedErrorCode)
			} else if tc.expectedErrorCode != "" {
				t.Errorf("\texpected an error to occurr: %s", tc.expectedErrorCode)
			} else {
				if encodedStandardClaims != tc.expectedEncodedStandardClaims {
					t.Errorf("\tencodedStandardClaims not expected:\ng - %s\ne - %s", encodedStandardClaims, tc.expectedEncodedStandardClaims)
				}
			}
		})
	}
}

func TestCalculateHMACSignature(t *testing.T) {
	type testCase struct {
		name                 string
		hashFunc             func() hash.Hash
		secret               string
		encodedHeaderAndBody string
		expectedSignature    string
	}
	testCases := []testCase{
		{
			name:                 "GIVEN an encoded header and body, a secret key and the sha256 hashing algorithm EXPECT the proper signature to be produced",
			hashFunc:             sha256.New,
			secret:               "super secret key",
			encodedHeaderAndBody: "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpc3MiOiJnb2F1dGgiLCJzdWIiOiIxMjM0NTY3ODkwIiwiYXVkIjpbImdvYXV0aCJdLCJleHAiOjE1MTYyMzkwMjIsIm5iZiI6MTUxNjIzOTAyMiwiaWF0IjoxNTE2MjM5MDIyLCJqdGkiOiI3ODk0NTYrMTIzMDMyMTM2NTQ5ODQ5Njg1NDErNDQ1NTUyIn0",
			expectedSignature:    "L76z0VDzq2744_Da9T4YbfZ5rYQvC_bSevCq5rhpS3k",
		},
		{
			name:                 "GIVEN an encoded header and body, a secret key and the sha384 hashing algorithm EXPECT the proper signature to be produced",
			hashFunc:             sha512.New384,
			secret:               "super secret key",
			encodedHeaderAndBody: "eyJhbGciOiJIUzM4NCIsInR5cCI6IkpXVCJ9.eyJpc3MiOiJnb2F1dGgiLCJzdWIiOiIxMjM0NTY3ODkwIiwiYXVkIjpbImdvYXV0aCJdLCJleHAiOjE1MTYyMzkwMjIsIm5iZiI6MTUxNjIzOTAyMiwiaWF0IjoxNTE2MjM5MDIyLCJqdGkiOiI3ODk0NTYrMTIzMDMyMTM2NTQ5ODQ5Njg1NDErNDQ1NTUyIn0",
			expectedSignature:    "6D6tQdMbij1bQ1-4QiNgp0E7KSDjuju9-S3E_ItKIuOUIfjHvW3eJXCTJG37Rp5B",
		},
		{
			name:                 "GIVEN an encoded header and body, a secret key and the sha512 hashing algorithm EXPECT the proper signature to be produced",
			hashFunc:             sha512.New,
			secret:               "super secret key",
			encodedHeaderAndBody: "eyJhbGciOiJIUzUxMiIsInR5cCI6IkpXVCJ9.eyJpc3MiOiJnb2F1dGgiLCJzdWIiOiIxMjM0NTY3ODkwIiwiYXVkIjpbImdvYXV0aCJdLCJleHAiOjE1MTYyMzkwMjIsIm5iZiI6MTUxNjIzOTAyMiwiaWF0IjoxNTE2MjM5MDIyLCJqdGkiOiI3ODk0NTYrMTIzMDMyMTM2NTQ5ODQ5Njg1NDErNDQ1NTUyIn0",
			expectedSignature:    "CBXefbK0WRMCiAYZIr1eujIJyCtomrEPTO4eEOorojKweDF5BLBua6ThlGDKsICkSEjUG5s6NxneSbGZ-cmUzQ",
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			signature := CalculateHMACSignature(tc.secret, tc.encodedHeaderAndBody, tc.hashFunc)
			if signature != tc.expectedSignature {
				t.Errorf("\tsignature is not expected value\n g - %s\ne - %s", signature, tc.expectedSignature)
			}
		})
	}
}

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
