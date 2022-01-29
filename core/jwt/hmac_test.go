package jwt

import (
	"testing"

	"github.com/calvine/goauth/internal/testutils"
)

func TestCalculateHMACSignature(t *testing.T) {
	type testCase struct {
		name                 string
		hmacOptions          HMACSigningOptions
		algorithm            JWTSigningAlgorithm
		encodedHeaderAndBody string
		expectedSignature    string
		expectedErrorCode    string
	}
	testCases := []testCase{
		{
			name: "GIVEN an encoded header and body, a secret key and the sha256 hashing algorithm EXPECT the proper signature to be produced",
			hmacOptions: HMACSigningOptions{
				Secret: "super secret key",
			},
			algorithm:            "HS256",
			encodedHeaderAndBody: "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpc3MiOiJnb2F1dGgiLCJzdWIiOiIxMjM0NTY3ODkwIiwiYXVkIjpbImdvYXV0aCJdLCJleHAiOjE1MTYyMzkwMjIsIm5iZiI6MTUxNjIzOTAyMiwiaWF0IjoxNTE2MjM5MDIyLCJqdGkiOiI3ODk0NTYrMTIzMDMyMTM2NTQ5ODQ5Njg1NDErNDQ1NTUyIn0",
			expectedSignature:    "L76z0VDzq2744_Da9T4YbfZ5rYQvC_bSevCq5rhpS3k",
		},
		{
			name: "GIVEN an encoded header and body, a secret key and the sha384 hashing algorithm EXPECT the proper signature to be produced",
			hmacOptions: HMACSigningOptions{
				Secret: "super secret key",
			},
			algorithm:            "HS384",
			encodedHeaderAndBody: "eyJhbGciOiJIUzM4NCIsInR5cCI6IkpXVCJ9.eyJpc3MiOiJnb2F1dGgiLCJzdWIiOiIxMjM0NTY3ODkwIiwiYXVkIjpbImdvYXV0aCJdLCJleHAiOjE1MTYyMzkwMjIsIm5iZiI6MTUxNjIzOTAyMiwiaWF0IjoxNTE2MjM5MDIyLCJqdGkiOiI3ODk0NTYrMTIzMDMyMTM2NTQ5ODQ5Njg1NDErNDQ1NTUyIn0",
			expectedSignature:    "6D6tQdMbij1bQ1-4QiNgp0E7KSDjuju9-S3E_ItKIuOUIfjHvW3eJXCTJG37Rp5B",
		},
		{
			name: "GIVEN an encoded header and body, a secret key and the sha512 hashing algorithm EXPECT the proper signature to be produced",
			hmacOptions: HMACSigningOptions{
				Secret: "super secret key",
			},
			algorithm:            "HS512",
			encodedHeaderAndBody: "eyJhbGciOiJIUzUxMiIsInR5cCI6IkpXVCJ9.eyJpc3MiOiJnb2F1dGgiLCJzdWIiOiIxMjM0NTY3ODkwIiwiYXVkIjpbImdvYXV0aCJdLCJleHAiOjE1MTYyMzkwMjIsIm5iZiI6MTUxNjIzOTAyMiwiaWF0IjoxNTE2MjM5MDIyLCJqdGkiOiI3ODk0NTYrMTIzMDMyMTM2NTQ5ODQ5Njg1NDErNDQ1NTUyIn0",
			expectedSignature:    "CBXefbK0WRMCiAYZIr1eujIJyCtomrEPTO4eEOorojKweDF5BLBua6ThlGDKsICkSEjUG5s6NxneSbGZ-cmUzQ",
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			signature, err := tc.hmacOptions.Sign(tc.algorithm, tc.encodedHeaderAndBody)
			if err != nil {
				testutils.HandleTestError(t, err, tc.expectedErrorCode)
			} else if tc.expectedErrorCode != "" {
				t.Errorf("\texpected an error to occurr: %s", tc.expectedErrorCode)
			} else {
				// signature := CalculateHMACSignature(tc.secret, tc.encodedHeaderAndBody, tc.hashFunc)
				if signature != tc.expectedSignature {
					t.Errorf("\tsignature is not expected value\n g - %s\ne - %s", signature, tc.expectedSignature)
				}
			}
		})
	}
}
