package models

import (
	"testing"
	"time"

	"github.com/calvine/goauth/core/nullable"
)

func TestJWTSigningMaterial_IsExpired(t *testing.T) {
	type testCase struct {
		name               string
		jwtSigningMaterial JWTSigningMaterial
		expectedResult     bool
	}
	testCases := []testCase{
		{
			name: "GIVEN a non expired jwt signing material EXPECT false",
			jwtSigningMaterial: JWTSigningMaterial{
				Expiration: nullable.NullableTime{
					HasValue: true,
					Value:    time.Now().Add(time.Minute),
				},
			},
			expectedResult: false,
		},
		{
			name: "GIVEN an expired jwt signing material EXPECT true",
			jwtSigningMaterial: JWTSigningMaterial{
				Expiration: nullable.NullableTime{
					HasValue: true,
					Value:    time.Now().Add(-time.Minute),
				},
			},
			expectedResult: true,
		},
		{
			name: "GIVEN a jwt signing material with no expiration value EXPECT false",
			jwtSigningMaterial: JWTSigningMaterial{
				Expiration: nullable.NullableTime{
					HasValue: false,
				},
			},
			expectedResult: false,
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result := tc.jwtSigningMaterial.IsExpired()
			if result != tc.expectedResult {
				t.Errorf("\tresult was not expected value: got - %v expected - %v", result, tc.expectedResult)
			}
		})
	}
}
