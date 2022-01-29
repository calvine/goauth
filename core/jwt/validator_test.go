package jwt

import (
	"testing"
	"time"

	coreerrors "github.com/calvine/goauth/core/errors"
	"github.com/calvine/goauth/internal/testutils"
)

func TestNewJWTValidator(t *testing.T) {
	type testCase struct {
		name                string
		jwtValidatorOptions JWTValidatorOptions
		expectedErrorCode   string
	}
	testCases := []testCase{
		{
			name: "GIVEN a valid set of jwt validator options EXPECT no errors to occur",
			jwtValidatorOptions: JWTValidatorOptions{
				AllowedAlgorithms: []JWTSigningAlgorithm{
					HS256,
				},
				HMACOptions: HMACSigningOptions{
					Secret: "test",
				},
			},
		},
		{
			name: "GIVEN jwt validator options with allow any issuer true and a value set for expected issuer EXPECT error code jwt validator no hmac secret provided",
			jwtValidatorOptions: JWTValidatorOptions{
				AllowedAlgorithms: []JWTSigningAlgorithm{
					HS256,
				},
				HMACOptions: HMACSigningOptions{
					Secret: "test",
				},
				AllowAnyIssuer: true,
				ExpectedIssuer: "goauth",
			},
			expectedErrorCode: coreerrors.ErrCodeJWTValidatorAllowAnyIssuerAndExpectedIssuerProvided,
		},
		{
			name: "GIVEN jwt validator options with HS algorithms allowed and no hmacSecret set EXPECT error code jwt validator no hmac secret provided",
			jwtValidatorOptions: JWTValidatorOptions{
				AllowedAlgorithms: []JWTSigningAlgorithm{
					HS256,
				},
			},
			expectedErrorCode: coreerrors.ErrCodeJWTValidatorNoHMACSecretProvided,
		},
		{
			name:                "GIVEN jwt validator options with no allowed algorithms EXPECT error code jwt validator no algorithm specified",
			jwtValidatorOptions: JWTValidatorOptions{},
			expectedErrorCode:   coreerrors.ErrCodeJWTValidatorNoAlgorithmSpecified,
		},
		{
			name: "GIVEN jwt validator options with aidience required true and no allowed audiences provided and allowAnyAudience false EXPECT error code jwt validator audience required but none provided",
			jwtValidatorOptions: JWTValidatorOptions{
				AllowedAlgorithms: []JWTSigningAlgorithm{
					HS256,
				},
				HMACOptions: HMACSigningOptions{
					Secret: "test",
				},
				AudienceRequired: true,
			},
			expectedErrorCode: coreerrors.ErrCodeJWTValidatorAudienceRequiredButNoneProvided,
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			validator, err := NewJWTValidator(tc.jwtValidatorOptions)
			if err != nil {
				testutils.HandleTestError(t, err, tc.expectedErrorCode)
			} else if tc.expectedErrorCode != "" {
				t.Errorf("\texpected an error to occurr: %s", tc.expectedErrorCode)
			} else {
				if validator.GetID() == "" {
					t.Error("\tvalidator id should never be an empty string")
				}
			}
		})
	}
}

func TestValidateHeader(t *testing.T) {
	type testCase struct {
		name                string
		jwtValidatorOptions JWTValidatorOptions
		header              Header
		expectValid         bool
		expectedErrorCodes  []string
	}
	testCases := []testCase{
		// algorithm tests
		{
			name: "GIVEN a header with an allowed algorithm EXPECT no errors to occur",
			jwtValidatorOptions: JWTValidatorOptions{
				AllowedAlgorithms: []JWTSigningAlgorithm{
					HS256,
					HS384,
					HS512,
				},
				HMACOptions: HMACSigningOptions{
					Secret: "test secret",
				},
			},
			header: Header{
				Algorithm: HS256,
			},
			expectValid: true,
		},
		{
			name: "GIVEN a header with an algorithm that is not in the allowed algorithms list EXPECT error code jwt algorithm not allowed",
			jwtValidatorOptions: JWTValidatorOptions{
				AllowedAlgorithms: []JWTSigningAlgorithm{
					HS384,
					HS512,
				},
				HMACOptions: HMACSigningOptions{
					Secret: "test secret",
				},
			},
			header: Header{
				Algorithm: HS256,
			},
			expectedErrorCodes: []string{
				coreerrors.ErrCodeJWTAlgorithmNotAllowed,
			},
			expectValid: false,
		},
		// keyid tests
		{
			name: "GIVEN a header with a key id while key id is not required algorithm EXPECT no errors to occur",
			jwtValidatorOptions: JWTValidatorOptions{
				AllowedAlgorithms: []JWTSigningAlgorithm{
					HS256,
					HS384,
					HS512,
				},
				HMACOptions: HMACSigningOptions{
					Secret: "test secret",
				},
				KeyIDRequired: false,
			},
			header: Header{
				Algorithm: HS256,
				KeyID:     "1234",
			},
			expectValid: true,
		},
		{
			name: "GIVEN a header with a key id while key id is required EXPECT no errors to occur",
			jwtValidatorOptions: JWTValidatorOptions{
				AllowedAlgorithms: []JWTSigningAlgorithm{
					HS256,
					HS384,
					HS512,
				},
				HMACOptions: HMACSigningOptions{
					Secret: "test secret",
				},
				KeyIDRequired: true,
			},
			header: Header{
				Algorithm: HS256,
				KeyID:     "1234",
			},
			expectValid: true,
		},
		{
			name: "GIVEN a header with no key id while the key id is required EXPECT error code jwt validator missing key id",
			jwtValidatorOptions: JWTValidatorOptions{
				AllowedAlgorithms: []JWTSigningAlgorithm{
					HS256,
					HS384,
					HS512,
				},
				HMACOptions: HMACSigningOptions{
					Secret: "test secret",
				},
				KeyIDRequired: true,
			},
			header: Header{
				Algorithm: HS256,
			},
			expectedErrorCodes: []string{
				coreerrors.ErrCodeJWTValidatorKeyIDMissing,
			},
			expectValid: false,
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			validator, err := NewJWTValidator(tc.jwtValidatorOptions)
			if err != nil {
				t.Errorf("\texpected an error to occurr while building JWT validator: %s", err.GetErrorCode())
			}
			errs, valid := validator.ValidateHeader(tc.header)
			numErrors := len(errs)
			numExpectedErrors := len(tc.expectedErrorCodes)
			if numErrors != numExpectedErrors {
				t.Errorf("\texpected number of errors incorrect: got - %d expected - %d", numErrors, numExpectedErrors)
			}
			for _, e := range errs {
				errorCode := e.GetErrorCode()
				found := false
				for _, eec := range tc.expectedErrorCodes {
					if errorCode == eec {
						found = true
						break
					}
				}
				if !found {
					t.Errorf("\terror occurred that was not expected: %s", errorCode)
				}
			}
			if valid != tc.expectValid {
				t.Errorf("\texpected header valid state incorrect: got - %v expected - %v", valid, tc.expectValid)
			}
		})
	}
}

func TestValidateClaims(t *testing.T) {
	type testCase struct {
		name                string
		jwtValidatorOptions JWTValidatorOptions
		body                StandardClaims
		expectValid         bool
		expectedErrorCodes  []string
	}
	testCases := []testCase{
		// issuer tests
		{
			name: "GIVEN claims with an expected issuer when no issuer is required EXPECT success",
			jwtValidatorOptions: JWTValidatorOptions{
				AllowedAlgorithms: []JWTSigningAlgorithm{
					HS256,
				},
				HMACOptions: HMACSigningOptions{
					Secret: "test",
				},
				ExpectedIssuer: "goauth",
				IssuerRequired: false,
			},
			body: StandardClaims{
				Issuer: "goauth",
			},
			expectValid: true,
		},
		{
			name: "GIVEN claims with an expected issuer when issuer is required EXPECT success",
			jwtValidatorOptions: JWTValidatorOptions{
				AllowedAlgorithms: []JWTSigningAlgorithm{
					HS256,
				},
				HMACOptions: HMACSigningOptions{
					Secret: "test",
				},
				ExpectedIssuer: "goauth",
				IssuerRequired: true,
			},
			body: StandardClaims{
				Issuer: "goauth",
			},
			expectValid: true,
		},
		{
			name: "GIVEN claims with an issuer and allow any issuer is true when issuer is required EXPECT success",
			jwtValidatorOptions: JWTValidatorOptions{
				AllowedAlgorithms: []JWTSigningAlgorithm{
					HS256,
				},
				HMACOptions: HMACSigningOptions{
					Secret: "test",
				},
				AllowAnyIssuer: true,
				IssuerRequired: true,
			},
			body: StandardClaims{
				Issuer: "goauth",
			},
			expectValid: true,
		},
		{
			name: "GIVEN claims with no issuer when issuer is not required EXPECT success",
			jwtValidatorOptions: JWTValidatorOptions{
				AllowedAlgorithms: []JWTSigningAlgorithm{
					HS256,
				},
				HMACOptions: HMACSigningOptions{
					Secret: "test",
				},
				IssuerRequired: false,
			},
			expectValid: true,
			body:        StandardClaims{},
		},
		{
			name: "GIVEN claims with an unexpected issuer when issuer is required EXPECT failure",
			jwtValidatorOptions: JWTValidatorOptions{
				AllowedAlgorithms: []JWTSigningAlgorithm{
					HS256,
				},
				HMACOptions: HMACSigningOptions{
					Secret: "test",
				},
				ExpectedIssuer: "goauth",
				IssuerRequired: false,
			},
			body: StandardClaims{
				Issuer: "other",
			},
			expectValid: false,
			expectedErrorCodes: []string{
				coreerrors.ErrCodeJWTIssuerInvalid,
			},
		},
		{
			name: "GIVEN claims with an unexpected issuer when issuer is not requiredEXPECT failure",
			jwtValidatorOptions: JWTValidatorOptions{
				AllowedAlgorithms: []JWTSigningAlgorithm{
					HS256,
				},
				HMACOptions: HMACSigningOptions{
					Secret: "test",
				},
				ExpectedIssuer: "goauth",
				IssuerRequired: true,
			},
			body: StandardClaims{
				Issuer: "other",
			},
			expectValid: false,
			expectedErrorCodes: []string{
				coreerrors.ErrCodeJWTIssuerInvalid,
			},
		},
		{
			name: "GIVEN claims with no issuer when issuer is required EXPECT failure",
			jwtValidatorOptions: JWTValidatorOptions{
				AllowedAlgorithms: []JWTSigningAlgorithm{
					HS256,
				},
				HMACOptions: HMACSigningOptions{
					Secret: "test",
				},
				IssuerRequired: true,
			},
			body:        StandardClaims{},
			expectValid: false,
			expectedErrorCodes: []string{
				coreerrors.ErrCodeJWTIssuerMissing,
			},
		},
		// expire tests
		{
			name: "GIVEN a jwt with a valid expire EXPECT success",
			jwtValidatorOptions: JWTValidatorOptions{
				AllowedAlgorithms: []JWTSigningAlgorithm{
					HS256,
				},
				HMACOptions: HMACSigningOptions{
					Secret: "test",
				},
			},
			body: StandardClaims{
				ExpirationTime: Time(time.Now().Add(time.Second * 10)),
			},
			expectValid: true,
		},
		{
			name: "GIVEN a jwt with a valid expire and exipre is required EXPECT success",
			jwtValidatorOptions: JWTValidatorOptions{
				AllowedAlgorithms: []JWTSigningAlgorithm{
					HS256,
				},
				ExpireRequired: true,
				HMACOptions: HMACSigningOptions{
					Secret: "test",
				},
			},
			body: StandardClaims{
				ExpirationTime: Time(time.Now().Add(time.Second * 10)),
			},
			expectValid: true,
		},
		{
			name: "GIVEN a jwt with no exipre and expire is not required EXPECT success",
			jwtValidatorOptions: JWTValidatorOptions{
				AllowedAlgorithms: []JWTSigningAlgorithm{
					HS256,
				},
				ExpireRequired: false,
				HMACOptions: HMACSigningOptions{
					Secret: "test",
				},
			},
			body:        StandardClaims{},
			expectValid: true,
		},
		{
			name:        "GIVEN a jwt that has expired EXPECT error code jwt expired",
			expectValid: false,
			jwtValidatorOptions: JWTValidatorOptions{
				AllowedAlgorithms: []JWTSigningAlgorithm{
					HS256,
				},
				HMACOptions: HMACSigningOptions{
					Secret: "test",
				},
			},
			body: StandardClaims{
				ExpirationTime: Time(time.Now().Add(time.Second * -1)),
			},
			expectedErrorCodes: []string{
				coreerrors.ErrCodeJWTExpired,
			},
		},
		{
			name:        "GIVEN a jwt that has expired and exp is required EXPECT error code jwt expired",
			expectValid: false,
			jwtValidatorOptions: JWTValidatorOptions{
				AllowedAlgorithms: []JWTSigningAlgorithm{
					HS256,
				},
				ExpireRequired: true,
				HMACOptions: HMACSigningOptions{
					Secret: "test",
				},
			},
			body: StandardClaims{
				ExpirationTime: Time(time.Now().Add(time.Second * -1)),
			},
			expectedErrorCodes: []string{
				coreerrors.ErrCodeJWTExpired,
			},
		},
		{
			name:        "GIVEN a jwt that has no expireation but expire is required EXPECT error code jwt expire missing",
			expectValid: false,
			jwtValidatorOptions: JWTValidatorOptions{
				AllowedAlgorithms: []JWTSigningAlgorithm{
					HS256,
				},
				ExpireRequired: true,
				HMACOptions: HMACSigningOptions{
					Secret: "test",
				},
			},
			body: StandardClaims{},
			expectedErrorCodes: []string{
				coreerrors.ErrCodeJWTExpireMissing,
			},
		},
		// audience tests
		{
			name: "GIVEN a jwt with audiences in allowed audiences and audiences are required EXPECT success",
			jwtValidatorOptions: JWTValidatorOptions{
				AllowedAlgorithms: []JWTSigningAlgorithm{
					HS256,
				},
				AllowedAudience: []string{
					"test1",
					"test2",
				},
				AudienceRequired: true,
				HMACOptions: HMACSigningOptions{
					Secret: "test",
				},
			},
			body: StandardClaims{
				Audience: "test1",
			},
			expectValid: true,
		},
		{
			name: "GIVEN a jwt with audiences in allowed audiences and audiences are not required EXPECT success",
			jwtValidatorOptions: JWTValidatorOptions{
				AllowedAlgorithms: []JWTSigningAlgorithm{
					HS256,
				},
				AllowedAudience: []string{
					"test1",
					"test2",
				},
				AudienceRequired: false,
				HMACOptions: HMACSigningOptions{
					Secret: "test",
				},
			},
			body: StandardClaims{
				Audience: "test2",
			},
			expectValid: true,
		},
		{
			name: "GIVEN a jwt with multiple audiences in allowed audiences and audiences are not required EXPECT success",
			jwtValidatorOptions: JWTValidatorOptions{
				AllowedAlgorithms: []JWTSigningAlgorithm{
					HS256,
				},
				AllowedAudience: []string{
					"test1",
					"test2",
				},
				AudienceRequired: false,
				HMACOptions: HMACSigningOptions{
					Secret: "test",
				},
			},
			body: StandardClaims{
				Audience: "test1,test2",
			},
			expectValid: true,
		},
		{
			name: "GIVEN a jwt with audiences and allow any audience is true and audiences are required EXPECT success",
			jwtValidatorOptions: JWTValidatorOptions{
				AllowedAlgorithms: []JWTSigningAlgorithm{
					HS256,
				},
				AllowAnyAudience: true,
				AudienceRequired: true,
				HMACOptions: HMACSigningOptions{
					Secret: "test",
				},
			},
			body: StandardClaims{
				Audience: "other1",
			},
			expectValid: true,
		},
		{
			name: "GIVEN a jwt with audiences and allow any audience is true and audiences are not required EXPECT success",
			jwtValidatorOptions: JWTValidatorOptions{
				AllowedAlgorithms: []JWTSigningAlgorithm{
					HS256,
				},
				AllowAnyAudience: true,
				AudienceRequired: false,
				HMACOptions: HMACSigningOptions{
					Secret: "test",
				},
			},
			body: StandardClaims{
				Audience: "other2",
			},
			expectValid: true,
		},
		{
			name: "GIVEN a jwt with multiple audiences and one invalid audience in allowed audiences and audiences are not required EXPECT success",
			jwtValidatorOptions: JWTValidatorOptions{
				AllowedAlgorithms: []JWTSigningAlgorithm{
					HS256,
				},
				AllowedAudience: []string{
					"test1",
					"test2",
				},
				AudienceRequired: false,
				HMACOptions: HMACSigningOptions{
					Secret: "test",
				},
			},
			body: StandardClaims{
				Audience: "test1,other2",
			},
			expectValid: false,
			expectedErrorCodes: []string{
				coreerrors.ErrCodeJWTValidatorAudienceInvalid,
			},
		},
		{
			name: "GIVEN a jwt no audience and audiences are required EXPECT error code jwt audience missing",
			jwtValidatorOptions: JWTValidatorOptions{
				AllowedAlgorithms: []JWTSigningAlgorithm{
					HS256,
				},
				AllowedAudience: []string{
					"test1",
					"test2",
				},
				AudienceRequired: true,
				HMACOptions: HMACSigningOptions{
					Secret: "test",
				},
			},
			body:        StandardClaims{},
			expectValid: false,
			expectedErrorCodes: []string{
				coreerrors.ErrCodeJWTValidatorAudienceMissing,
			},
		},
		{
			name: "GIVEN a jwt with an audiences not in allowed audiences and audiences are required EXPECT error code jwt audience invalid",
			jwtValidatorOptions: JWTValidatorOptions{
				AllowedAlgorithms: []JWTSigningAlgorithm{
					HS256,
				},
				AllowedAudience: []string{
					"test1",
					"test2",
				},
				AudienceRequired: true,
				HMACOptions: HMACSigningOptions{
					Secret: "test",
				},
			},
			body: StandardClaims{
				Audience: "other1",
			},
			expectValid: false,
			expectedErrorCodes: []string{
				coreerrors.ErrCodeJWTValidatorAudienceInvalid,
			},
		},
		{
			name: "GIVEN a jwt with an audiences not in allowed audiences and audiences are not required EXPECT error code jwt audience invalid",
			jwtValidatorOptions: JWTValidatorOptions{
				AllowedAlgorithms: []JWTSigningAlgorithm{
					HS256,
				},
				AllowedAudience: []string{
					"test1",
					"test2",
				},
				AudienceRequired: false,
				HMACOptions: HMACSigningOptions{
					Secret: "test",
				},
			},
			body: StandardClaims{
				Audience: "other2",
			},
			expectValid: false,
			expectedErrorCodes: []string{
				coreerrors.ErrCodeJWTValidatorAudienceInvalid,
			},
		},
		// issued at tests
		{
			name: "GIVEN a jwt with a valid issued at EXPECT success",
			jwtValidatorOptions: JWTValidatorOptions{
				AllowedAlgorithms: []JWTSigningAlgorithm{
					HS256,
				},
				HMACOptions: HMACSigningOptions{
					Secret: "test",
				},
			},
			body: StandardClaims{
				IssuedAt: Time(time.Now().Add(time.Second * -10)),
			},
			expectValid: true,
		},
		{
			name: "GIVEN a jwt with a valid issued at and issued at is required EXPECT success",
			jwtValidatorOptions: JWTValidatorOptions{
				AllowedAlgorithms: []JWTSigningAlgorithm{
					HS256,
				},
				IssuedAtRequired: true,
				HMACOptions: HMACSigningOptions{
					Secret: "test",
				},
			},
			body: StandardClaims{
				IssuedAt: Time(time.Now().Add(time.Second * -10)),
			},
			expectValid: true,
		},
		{
			name: "GIVEN a jwt with no issued at and issued at is not required EXPECT success",
			jwtValidatorOptions: JWTValidatorOptions{
				AllowedAlgorithms: []JWTSigningAlgorithm{
					HS256,
				},
				IssuedAtRequired: false,
				HMACOptions: HMACSigningOptions{
					Secret: "test",
				},
			},
			body:        StandardClaims{},
			expectValid: true,
		},
		{
			name:        "GIVEN a jwt that has an invalid issued at EXPECT error code jwt issued at invalid",
			expectValid: false,
			jwtValidatorOptions: JWTValidatorOptions{
				AllowedAlgorithms: []JWTSigningAlgorithm{
					HS256,
				},
				HMACOptions: HMACSigningOptions{
					Secret: "test",
				},
			},
			body: StandardClaims{
				IssuedAt: Time(time.Now().Add(time.Second * 10)), // in the future
			},
			expectedErrorCodes: []string{
				coreerrors.ErrCodeJWTIssuedAtInvalid,
			},
		},
		{
			name:        "GIVEN a jwt that has an invalid issued at and issued at is required EXPECT error code jwt issued at invalid",
			expectValid: false,
			jwtValidatorOptions: JWTValidatorOptions{
				AllowedAlgorithms: []JWTSigningAlgorithm{
					HS256,
				},
				IssuedAtRequired: true,
				HMACOptions: HMACSigningOptions{
					Secret: "test",
				},
			},
			body: StandardClaims{
				IssuedAt: Time(time.Now().Add(time.Second * 10)),
			},
			expectedErrorCodes: []string{
				coreerrors.ErrCodeJWTIssuedAtInvalid,
			},
		},
		{
			name:        "GIVEN a jwt that has no issued at but expire is required EXPECT error code jwt issued at missing",
			expectValid: false,
			jwtValidatorOptions: JWTValidatorOptions{
				AllowedAlgorithms: []JWTSigningAlgorithm{
					HS256,
				},
				IssuedAtRequired: true,
				HMACOptions: HMACSigningOptions{
					Secret: "test",
				},
			},
			body: StandardClaims{},
			expectedErrorCodes: []string{
				coreerrors.ErrCodeJWTIssuedAtMissing,
			},
		},
		// not before tests
		{
			name: "GIVEN a jwt with a valid not before EXPECT success",
			jwtValidatorOptions: JWTValidatorOptions{
				AllowedAlgorithms: []JWTSigningAlgorithm{
					HS256,
				},
				HMACOptions: HMACSigningOptions{
					Secret: "test",
				},
			},
			body: StandardClaims{
				NotBefore: Time(time.Now().Add(time.Second * -10)),
			},
			expectValid: true,
		},
		{
			name: "GIVEN a jwt with a valid not before and not before is required EXPECT success",
			jwtValidatorOptions: JWTValidatorOptions{
				AllowedAlgorithms: []JWTSigningAlgorithm{
					HS256,
				},
				NotBeforeRequired: true,
				HMACOptions: HMACSigningOptions{
					Secret: "test",
				},
			},
			body: StandardClaims{
				NotBefore: Time(time.Now().Add(time.Second * -10)),
			},
			expectValid: true,
		},
		{
			name: "GIVEN a jwt with no not before and not before is not required EXPECT success",
			jwtValidatorOptions: JWTValidatorOptions{
				AllowedAlgorithms: []JWTSigningAlgorithm{
					HS256,
				},
				NotBeforeRequired: false,
				HMACOptions: HMACSigningOptions{
					Secret: "test",
				},
			},
			body:        StandardClaims{},
			expectValid: true,
		},
		{
			name:        "GIVEN a jwt that has an invalid not before EXPECT error code jwt not before invalid",
			expectValid: false,
			jwtValidatorOptions: JWTValidatorOptions{
				AllowedAlgorithms: []JWTSigningAlgorithm{
					HS256,
				},
				HMACOptions: HMACSigningOptions{
					Secret: "test",
				},
			},
			body: StandardClaims{
				NotBefore: Time(time.Now().Add(time.Second * 10)), // in the future
			},
			expectedErrorCodes: []string{
				coreerrors.ErrCodeJWTNotBeforeInFuture,
			},
		},
		{
			name:        "GIVEN a jwt that has an invalid not before and not before is required EXPECT error code jwt not before invalid",
			expectValid: false,
			jwtValidatorOptions: JWTValidatorOptions{
				AllowedAlgorithms: []JWTSigningAlgorithm{
					HS256,
				},
				NotBeforeRequired: true,
				HMACOptions: HMACSigningOptions{
					Secret: "test",
				},
			},
			body: StandardClaims{
				NotBefore: Time(time.Now().Add(time.Second * 10)),
			},
			expectedErrorCodes: []string{
				coreerrors.ErrCodeJWTNotBeforeInFuture,
			},
		},
		{
			name:        "GIVEN a jwt that has no not before but expire is required EXPECT error code jwt not before missing",
			expectValid: false,
			jwtValidatorOptions: JWTValidatorOptions{
				AllowedAlgorithms: []JWTSigningAlgorithm{
					HS256,
				},
				NotBeforeRequired: true,
				HMACOptions: HMACSigningOptions{
					Secret: "test",
				},
			},
			body: StandardClaims{},
			expectedErrorCodes: []string{
				coreerrors.ErrCodeJWTNotBeforeMissing,
			},
		},
		// subject tests
		{
			name: "GIVEN jwt with subject when the subject is required EXPECT success",
			jwtValidatorOptions: JWTValidatorOptions{
				AllowedAlgorithms: []JWTSigningAlgorithm{
					HS256,
				},
				HMACOptions: HMACSigningOptions{
					Secret: "test",
				},
				SubjectRequired: true,
			},
			body: StandardClaims{
				Subject: "user id",
			},
			expectValid: true,
		},
		{
			name: "GIVEN jwt with subject when the subject is not required EXPECT success",
			jwtValidatorOptions: JWTValidatorOptions{
				AllowedAlgorithms: []JWTSigningAlgorithm{
					HS256,
				},
				HMACOptions: HMACSigningOptions{
					Secret: "test",
				},
				SubjectRequired: false,
			},
			body: StandardClaims{
				Subject: "user id",
			},
			expectValid: true,
		},
		{
			name: "GIVEN jwt with no subject when the subject is not required EXPECT success",
			jwtValidatorOptions: JWTValidatorOptions{
				AllowedAlgorithms: []JWTSigningAlgorithm{
					HS256,
				},
				HMACOptions: HMACSigningOptions{
					Secret: "test",
				},
				SubjectRequired: false,
			},
			body:        StandardClaims{},
			expectValid: true,
		},
		{
			name: "GIVEN jwt with no id when the jwt id is required EXPECT error code jwt id missing",
			jwtValidatorOptions: JWTValidatorOptions{
				AllowedAlgorithms: []JWTSigningAlgorithm{
					HS256,
				},
				HMACOptions: HMACSigningOptions{
					Secret: "test",
				},
				SubjectRequired: true,
			},
			body:        StandardClaims{},
			expectValid: false,
			expectedErrorCodes: []string{
				coreerrors.ErrCodeJWTSubjectMissing,
			},
		},
		// jwt id tests
		{
			name: "GIVEN jwt with an id when the jwt id is required EXPECT success",
			jwtValidatorOptions: JWTValidatorOptions{
				AllowedAlgorithms: []JWTSigningAlgorithm{
					HS256,
				},
				HMACOptions: HMACSigningOptions{
					Secret: "test",
				},
				JTIRequired: true,
			},
			body: StandardClaims{
				JWTID: "unique id",
			},
			expectValid: true,
		},
		{
			name: "GIVEN jwt with an id when the jwt id is not required EXPECT success",
			jwtValidatorOptions: JWTValidatorOptions{
				AllowedAlgorithms: []JWTSigningAlgorithm{
					HS256,
				},
				HMACOptions: HMACSigningOptions{
					Secret: "test",
				},
				JTIRequired: false,
			},
			body: StandardClaims{
				JWTID: "unique id",
			},
			expectValid: true,
		},
		{
			name: "GIVEN jwt with an id when the jwt id is required EXPECT success",
			jwtValidatorOptions: JWTValidatorOptions{
				AllowedAlgorithms: []JWTSigningAlgorithm{
					HS256,
				},
				HMACOptions: HMACSigningOptions{
					Secret: "test",
				},
				JTIRequired: true,
			},
			body: StandardClaims{
				JWTID: "unique id",
			},
			expectValid: true,
		},
		{
			name: "GIVEN jwt with no id when the jwt id is not required EXPECT success",
			jwtValidatorOptions: JWTValidatorOptions{
				AllowedAlgorithms: []JWTSigningAlgorithm{
					HS256,
				},
				HMACOptions: HMACSigningOptions{
					Secret: "test",
				},
				JTIRequired: false,
			},
			body:        StandardClaims{},
			expectValid: true,
		},
		{
			name: "GIVEN jwt with no id when the jwt id is required EXPECT error code jwt id missing",
			jwtValidatorOptions: JWTValidatorOptions{
				AllowedAlgorithms: []JWTSigningAlgorithm{
					HS256,
				},
				HMACOptions: HMACSigningOptions{
					Secret: "test",
				},
				JTIRequired: true,
			},
			body:        StandardClaims{},
			expectValid: false,
			expectedErrorCodes: []string{
				coreerrors.ErrCodeJWTIDMissing,
			},
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			validator, err := NewJWTValidator(tc.jwtValidatorOptions)
			if err != nil {
				t.Errorf("\tunexpected an error to occurr while building JWT validator: %s", err.GetErrorCode())
			}
			errs, valid := validator.ValidateClaims(tc.body)
			numErrors := len(errs)
			numExpectedErrors := len(tc.expectedErrorCodes)
			if numErrors != numExpectedErrors {
				t.Errorf("\texpected number of errors incorrect: got - %d expected - %d", numErrors, numExpectedErrors)
			}
			for _, e := range errs {
				errorCode := e.GetErrorCode()
				found := false
				for _, eec := range tc.expectedErrorCodes {
					if errorCode == eec {
						found = true
						break
					}
				}
				if !found {
					t.Errorf("\terror occurred that was not expected: %s", errorCode)
				}
			}
			if valid != tc.expectValid {
				t.Errorf("\texpected claims valid state incorrect: got - %v expected - %v", valid, tc.expectValid)
			}
		})
	}
}

func TestValidateSignature(t *testing.T) {
	type testCase struct {
		name                 string
		jwtValidatorOptions  JWTValidatorOptions
		alg                  JWTSigningAlgorithm
		encodedHeaderAndBody string
		signature            string
		expectValid          bool
		expectedErrorCode    string
	}
	testCases := []testCase{
		{
			name: "GIVEN a valid HS256 jwt signature for the given encodedHeadAndBody EXPECT the signature to be valid",
			jwtValidatorOptions: JWTValidatorOptions{
				AllowedAlgorithms: []JWTSigningAlgorithm{
					HS256,
				},
				HMACOptions: HMACSigningOptions{
					Secret: "test",
				},
			},
			alg:                  HS256,
			encodedHeaderAndBody: "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxMjM0NTY3ODkwIiwibmFtZSI6IkpvaG4gRG9lIiwiaWF0IjoxNTE2MjM5MDIyfQ",
			signature:            "5mhBHqs5_DTLdINd9p5m7ZJ6XD0Xc55kIaCRY5r6HRA",
			expectValid:          true,
		},
		{
			name: "GIVEN a valid HS384 jwt signature for the given encodedHeadAndBody EXPECT the signature to be valid",
			jwtValidatorOptions: JWTValidatorOptions{
				AllowedAlgorithms: []JWTSigningAlgorithm{
					HS384,
				},
				HMACOptions: HMACSigningOptions{
					Secret: "test",
				},
			},
			alg:                  HS384,
			encodedHeaderAndBody: "eyJhbGciOiJIUzM4NCIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxMjM0NTY3ODkwIiwibmFtZSI6IkpvaG4gRG9lIiwiaWF0IjoxNTE2MjM5MDIyfQ",
			signature:            "KOZqnJ-wEzC-JvqqIHGKBIGgbYHH2Fej71TpBctnIguBkf3EdSYiwuRMSz35uY8E",
			expectValid:          true,
		},
		{
			name: "GIVEN a valid HS512 jwt signature for the given encodedHeadAndBody EXPECT the signature to be valid",
			jwtValidatorOptions: JWTValidatorOptions{
				AllowedAlgorithms: []JWTSigningAlgorithm{
					HS512,
				},
				HMACOptions: HMACSigningOptions{
					Secret: "test",
				},
			},
			alg:                  HS512,
			encodedHeaderAndBody: "eyJhbGciOiJIUzUxMiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxMjM0NTY3ODkwIiwibmFtZSI6IkpvaG4gRG9lIiwiaWF0IjoxNTE2MjM5MDIyfQ",
			signature:            "VXfjNdZn9mDxRYhiaCi8rYYtcuNe3KCfK3LvggWSaHwjZsag9ugMOuDPOeeBD3oNhK-cOkTvRLy_ERbgnEyxYA",
			expectValid:          true,
		},
		{
			name: "GIVEN an invalid HS256 jwt signature for the given encodedHeadAndBody EXPECT the signature to be valid",
			jwtValidatorOptions: JWTValidatorOptions{
				AllowedAlgorithms: []JWTSigningAlgorithm{
					HS256,
				},
				HMACOptions: HMACSigningOptions{
					Secret: "test",
				},
			},
			alg:                  HS256,
			encodedHeaderAndBody: "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxMjM0NTY3ODkwIiwibmFtZSI6IkpvaG4gRG9lIiwiaWF0IjoxNTE2MjM5MDIyfQ",
			signature:            "5mhhHqs5_DTLdINd9p5m7ZJ6XD0Xc55kIaCRY5r6HRA",
			expectValid:          false,
		},
		{
			name: "GIVEN an invalid HS384 jwt signature for the given encodedHeadAndBody EXPECT the signature to be valid",
			jwtValidatorOptions: JWTValidatorOptions{
				AllowedAlgorithms: []JWTSigningAlgorithm{
					HS384,
				},
				HMACOptions: HMACSigningOptions{
					Secret: "test",
				},
			},
			alg:                  HS384,
			encodedHeaderAndBody: "eyJhbGciOiJIUzM4NCIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxMjM0NTY3ODkwIiwibmFtZSI6IkpvaG4gRG9lIiwiaWF0IjoxNTE2MjM5MDIyfQ",
			signature:            "KOZqqJ-wEzC-JvqqIHGKBIGgbYHH2Fej71TpBctnIguBkf3EdSYiwuRMSz35uY8E",
			expectValid:          false,
		},
		{
			name: "GIVEN an invalid HS512 jwt signature for the given encodedHeadAndBody EXPECT the signature to be valid",
			jwtValidatorOptions: JWTValidatorOptions{
				AllowedAlgorithms: []JWTSigningAlgorithm{
					HS512,
				},
				HMACOptions: HMACSigningOptions{
					Secret: "test",
				},
			},
			alg:                  HS512,
			encodedHeaderAndBody: "eyJhbGciOiJIUzUxMiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxMjM0NTY3ODkwIiwibmFtZSI6IkpvaG4gRG9lIiwiaWF0IjoxNTE2MjM5MDIyfQ",
			signature:            "VXfjNddn9mDxRYhiaCi8rYYtcuNe3KCfK3LvggWSaHwjZsag9ugMOuDPOeeBD3oNhK-cOkTvRLy_ERbgnEyxYA",
			expectValid:          false,
		},
		{
			name: "GIVEN a non supported jwt signature algorithm for the given encodedHeadAndBody EXPECT error code jwt algorithm not implemented",
			jwtValidatorOptions: JWTValidatorOptions{
				AllowedAlgorithms: []JWTSigningAlgorithm{
					NONE,
				},
			},
			alg:                  NONE,
			encodedHeaderAndBody: "eyJhbGciOiJIUzUxMiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxMjM0NTY3ODkwIiwibmFtZSI6IkpvaG4gRG9lIiwiaWF0IjoxNTE2MjM5MDIyfQ",
			signature:            "VXfjNddn9mDxRYhiaCi8rYYtcuNe3KCfK3LvggWSaHwjZsag9ugMOuDPOeeBD3oNhK-cOkTvRLy_ERbgnEyxYA",
			expectValid:          false,
			expectedErrorCode:    coreerrors.ErrCodeJWTAlgorithmNotImplemented,
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			validator, err := NewJWTValidator(tc.jwtValidatorOptions)
			if err != nil {
				t.Errorf("\texpected an error to occurr while building JWT validator: %s", err.GetErrorCode())
			}
			valid, err := validator.ValidateSignature(tc.alg, tc.encodedHeaderAndBody, tc.signature)
			if err != nil {
				testutils.HandleTestError(t, err, tc.expectedErrorCode)
			} else if tc.expectedErrorCode != "" {
				t.Errorf("\texpected an error to occurr: %s", tc.expectedErrorCode)
			}
			if valid != tc.expectValid {
				t.Errorf("\texpected signature valid state incorrect: got - %v expected - %v", valid, tc.expectValid)
			}
		})
	}
}
