package models

import (
	"testing"
	"time"

	"github.com/calvine/goauth/core/constants/contact"
	"github.com/calvine/goauth/core/nullable"
)

func TestContactIsConfirmed(t *testing.T) {
	type testCase struct {
		name                string
		confirmedDate       nullable.NullableTime
		expectedIsConfirmed bool
	}
	testCases := []testCase{
		{
			name: "GIVEN valid confirmedDate EXPECT true",
			confirmedDate: nullable.NullableTime{
				HasValue: true,
				Value:    time.Now().Add(time.Second * -1),
			},
			expectedIsConfirmed: true,
		},
		{
			name: "GIVEN future confirmedDate EXPECT false",
			confirmedDate: nullable.NullableTime{
				HasValue: true,
				Value:    time.Now().Add(time.Second * 5),
			},
			expectedIsConfirmed: false,
		},
		{
			name:                "GIVEN null confirmedDate EXPECT false",
			expectedIsConfirmed: false,
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			contact := Contact{
				ConfirmedDate: tc.confirmedDate,
			}
			isConfirmed := contact.IsConfirmed()
			if isConfirmed != tc.expectedIsConfirmed {
				t.Errorf("\tconfirmed status is not what was expected: got - %v expected - %v", isConfirmed, tc.expectedIsConfirmed)
			}
		})
	}
}

func TestNormalizeContactPrincipal(t *testing.T) {
	type testCase struct {
		name                               string
		contactType                        contact.Type
		contactPrincipal                   string
		expectedNormalizedContactPrincipal string
	}
	testCases := []testCase{
		{
			name:                               "GIVEN a email contact EXPECT the contact principal to be converted to lower case",
			contactType:                        contact.Email,
			contactPrincipal:                   "My_Email_123@email.org",
			expectedNormalizedContactPrincipal: "my_email_123@email.org",
		},
		{
			name:                               "GIVEN a mobile contact EXPECT the contact principal to have any dashes removed",
			contactType:                        contact.Mobile,
			contactPrincipal:                   "+1-478-867-5309",
			expectedNormalizedContactPrincipal: "+14788675309",
		},
		{
			name:                               "GIVEN an invalid contact type EXPECT the contact principal will be converted to lower case",
			contactType:                        "INVALID TYPE",
			contactPrincipal:                   "JustChecking!",
			expectedNormalizedContactPrincipal: "justchecking!",
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			normalizedContactPrincipal := NormalizeContactPrincipal(tc.contactType, tc.contactPrincipal)
			if normalizedContactPrincipal != tc.expectedNormalizedContactPrincipal {
				t.Errorf("\tnormalized contact principal is not expected: got - %s expect - %s", normalizedContactPrincipal, tc.expectedNormalizedContactPrincipal)
			}
		})
	}
}
