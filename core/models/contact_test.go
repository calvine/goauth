package models

import (
	"testing"
	"time"

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
				t.Fail()
			}
		})
	}
}
