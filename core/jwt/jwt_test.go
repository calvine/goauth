package jwt

import (
	"testing"
	"time"

	coreerrors "github.com/calvine/goauth/core/errors"
	"github.com/calvine/goauth/internal/testutils"
)

func TestDecodeAndValidateJWT(t *testing.T) {
	type testCase struct {
		name                      string
		encodedJWT                string
		validatorOptions          JWTValidatorOptions
		expectedErrorCode         string
		expectedJWTHeader         Header
		expectedJWTStandardClaims StandardClaims
	}
	testCases := []testCase{
		{
			name:       "GIVEN EXPECT ",
			encodedJWT: "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpc3MiOiJnb2F1dGgiLCJzdWIiOiIxMjM0NTY3ODkwIiwiYXVkIjpbImdvYXV0aCJdLCJuYmYiOjE1MTYyMzkwMjIsImlhdCI6MTUxNjIzOTAyMiwianRpIjoiNzg5NDU2KzEyMzAzMjEzNjU0OTg0OTY4NTQxKzQ0NTU1MiJ9.s-W_aZQE046I1SdfNaW4h5Yh1JgDXPcHQYWrSw0mfF0",
			validatorOptions: JWTValidatorOptions{
				AllowedAlgorithms: []string{
					Alg_HS256,
				},
				AllowAnyAudience: true,
				ExpectedIssuer:   "goauth",
				HMACOptions: HMACSigningOptions{
					Secret: "super secret key",
				},
			},
			expectedJWTHeader: Header{
				Algorithm: Alg_HS256,
				TokenType: Typ_JWT,
			},
			expectedJWTStandardClaims: StandardClaims{
				Issuer:    "goauth",
				Subject:   "1234567890",
				Audience:  "goauth",
				NotBefore: Time(time.Unix(1516239022, 0)),
				IssuedAt:  Time(time.Unix(1516239022, 0)),
				JWTID:     "789456+12303213654984968541+445552",
			},
		},
		{
			name:       "GIVEN EXPECT ",
			encodedJWT: "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpc3MiOiJnb2F1dGgjLCJzdWIiOiIxMjM0NTY3ODkwIiwiYXVkIjpbImdvYXV0aCJdLCJuYmYiOjE1MTYyMzkwMjIsImlhdCI6MTUxNjIzOTAyMiwianRpIjoiNzg5NDU2KzEyMzAzMjEzNjU0OTg0OTY4NTQxKzQ0NTU1MiJ9.s-W_aZQE046I1SdfNaW4h5Yh1JgDXPcHQYWrSw0mfF0",
			validatorOptions: JWTValidatorOptions{
				AllowedAlgorithms: []string{
					Alg_HS256,
				},
				AllowAnyAudience: true,
				ExpectedIssuer:   "goauth",
				HMACOptions: HMACSigningOptions{
					Secret: "super secret key",
				},
			},
			expectedErrorCode: coreerrors.ErrCodeJWTSignatureInvalid,
		},
		{
			name:       "GIVEN an expired token EXPECT error code jwt standard claims invalid ",
			encodedJWT: "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpc3MiOiJnb2F1dGgiLCJzdWIiOiIxMjM0NTY3ODkwIiwiYXVkIjpbImdvYXV0aCJdLCJleHAiOjE1MTYyMzkwMjIsIm5iZiI6MTUxNjIzOTAyMiwiaWF0IjoxNTE2MjM5MDIyLCJqdGkiOiI3ODk0NTYrMTIzMDMyMTM2NTQ5ODQ5Njg1NDErNDQ1NTUyIn0.L76z0VDzq2744_Da9T4YbfZ5rYQvC_bSevCq5rhpS3k",
			validatorOptions: JWTValidatorOptions{
				AllowedAlgorithms: []string{
					Alg_HS256,
				},
				AllowAnyAudience: true,
				ExpectedIssuer:   "goauth",
				HMACOptions: HMACSigningOptions{
					Secret: "super secret key",
				},
			},
			expectedErrorCode: coreerrors.ErrCodeJWTStandardClaimsInvalid,
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			validator, err := NewJWTValidator(tc.validatorOptions)
			if err != nil {
				t.Errorf("\tfailed to construct validator with given options: %s", err.Error())
				return
			}
			jwt, err := DecodeAndValidateJWT(tc.encodedJWT, validator)
			if err != nil {
				testutils.HandleTestError(t, err, tc.expectedErrorCode)
			} else if tc.expectedErrorCode != "" {
				t.Errorf("\texpected an error to occurr: %s", tc.expectedErrorCode)
			} else {
				if tc.expectedJWTHeader != jwt.Header {
					t.Errorf("\tjwt.Header value was not expected:\ng - %v\ne - %v", jwt.Header, tc.expectedJWTHeader)
				}
				if tc.expectedJWTStandardClaims != jwt.Claims {
					t.Errorf("\tjwt.Claims value was not expected:\ng - %v\ne - %v", jwt.Claims, tc.expectedJWTStandardClaims)
				}
			}

		})
	}
}

func TestEncodeAndSign(t *testing.T) {
	type testCase struct {
		name               string
		jwt                JWT
		signer             JWTSigner
		expectedEncodedJWT string
		expectedErrorCode  string
	}
	testCases := []testCase{
		{
			name: "GIVEN a valid jwt with the HS256 algorithm EXPECT a properly encoded jwt with an HS256 signature",
			jwt: JWT{
				Header: Header{
					Algorithm: Alg_HS256,
					TokenType: Typ_JWT,
				},
				Claims: StandardClaims{
					Issuer:         "goauth",
					Subject:        "1234567890",
					Audience:       "goauth,othersite",
					ExpirationTime: Time(time.Unix(1516239022, 0)),
					NotBefore:      Time(time.Unix(1516239022, 0)),
					IssuedAt:       Time(time.Unix(1516239022, 0)),
					JWTID:          "12345678900987654321+_+_--",
				},
			},
			signer: HMACSigningOptions{
				Secret: "super secret key",
			},
			expectedEncodedJWT: "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpc3MiOiJnb2F1dGgiLCJzdWIiOiIxMjM0NTY3ODkwIiwiYXVkIjpbImdvYXV0aCIsIm90aGVyc2l0ZSJdLCJleHAiOjE1MTYyMzkwMjIsIm5iZiI6MTUxNjIzOTAyMiwiaWF0IjoxNTE2MjM5MDIyLCJqdGkiOiIxMjM0NTY3ODkwMDk4NzY1NDMyMStfK18tLSJ9.bfdGZTx9ZPXwPg5GFJdyNGwghFwcyNfDBuDnlaTFh84",
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			encodedJWT, err := tc.jwt.EncodeAndSign(tc.signer)
			if err != nil {
				testutils.HandleTestError(t, err, tc.expectedErrorCode)
			} else if tc.expectedErrorCode != "" {
				t.Errorf("\texpected an error to occurr: %s", tc.expectedErrorCode)
			} else {
				if encodedJWT != tc.expectedEncodedJWT {
					t.Errorf("\tencodedJWT value was not expected:\ng - %v\ne - %v", encodedJWT, tc.expectedEncodedJWT)
				}
			}
		})
	}
}

func TestEncodeSignedJWT(t *testing.T) {
	type testCase struct {
		name               string
		jwt                JWT
		expectedEncodedJWT string
		expectedErrorCode  string
	}
	testCases := []testCase{
		{
			name: "GIVEN a valid jwt that is already signed EXPECT a properly encoded jwt",
			jwt: JWT{
				Header: Header{
					Algorithm: Alg_HS256,
					TokenType: Typ_JWT,
				},
				Claims: StandardClaims{
					Issuer:         "goauth",
					Subject:        "1234567890",
					Audience:       "goauth,othersite",
					ExpirationTime: Time(time.Unix(1516239022, 0)),
					NotBefore:      Time(time.Unix(1516239022, 0)),
					IssuedAt:       Time(time.Unix(1516239022, 0)),
					JWTID:          "12345678900987654321+_+_--",
				},
				Signature: "bfdGZTx9ZPXwPg5GFJdyNGwghFwcyNfDBuDnlaTFh84",
			},
			expectedEncodedJWT: "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpc3MiOiJnb2F1dGgiLCJzdWIiOiIxMjM0NTY3ODkwIiwiYXVkIjpbImdvYXV0aCIsIm90aGVyc2l0ZSJdLCJleHAiOjE1MTYyMzkwMjIsIm5iZiI6MTUxNjIzOTAyMiwiaWF0IjoxNTE2MjM5MDIyLCJqdGkiOiIxMjM0NTY3ODkwMDk4NzY1NDMyMStfK18tLSJ9.bfdGZTx9ZPXwPg5GFJdyNGwghFwcyNfDBuDnlaTFh84",
		},
		{
			name: "GIVEN a jwt witha missing signature EXPECT error code jwt signature missing",
			jwt: JWT{
				Header: Header{
					Algorithm: Alg_HS256,
					TokenType: Typ_JWT,
				},
				Claims: StandardClaims{
					Issuer:         "goauth",
					Subject:        "1234567890",
					Audience:       "goauth,othersite",
					ExpirationTime: Time(time.Unix(1516239022, 0)),
					NotBefore:      Time(time.Unix(1516239022, 0)),
					IssuedAt:       Time(time.Unix(1516239022, 0)),
					JWTID:          "12345678900987654321+_+_--",
				},
			},
			expectedErrorCode: coreerrors.ErrCodeJWTSignatureMissing,
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			encodedJWT, err := tc.jwt.EncodeSignedJWT()
			if err != nil {
				testutils.HandleTestError(t, err, tc.expectedErrorCode)
			} else if tc.expectedErrorCode != "" {
				t.Errorf("\texpected an error to occurr: %s", tc.expectedErrorCode)
			} else {
				if encodedJWT != tc.expectedEncodedJWT {
					t.Errorf("\tencodedJWT value was not expected:\ng - %v\ne - %v", encodedJWT, tc.expectedEncodedJWT)
				}
			}
		})
	}
}

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
				TokenType: Typ_JWT,
			},
		},
		{
			name:          "GIVEN an encoded jwt header with invalid json EXPECT error code jwt malformed",
			encodedHeader: "eyJhbGdiOiJIUzI1NiIsInR5cCI6IkpXVCJ9",
			expectedHeader: Header{
				Algorithm: Alg_HS256,
				TokenType: Typ_JWT,
			},
			expectedErrorCode: coreerrors.ErrCodeJWTMalformed,
		},
		{
			name:          "GIVEN an invalid base 64 encoded string EXPECT error code jwt malformed",
			encodedHeader: "eyJhbGdiOiJIUzI1NiIsQQQR5cCI6IkpXVCJ9",
			expectedHeader: Header{
				Algorithm: Alg_HS256,
				TokenType: Typ_JWT,
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
				TokenType: Typ_JWT,
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
