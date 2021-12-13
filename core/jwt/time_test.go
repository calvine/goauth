package jwt

import (
	"encoding/json"
	"testing"
	"time"
)

func TestMarshalJSON(t *testing.T) {
	type testCase struct {
		name                  string
		timestampString       string
		expectedUnixTimestamp string
	}
	testCases := []testCase{
		{
			name:                  "GIVEN a valid Time EXPECT success",
			timestampString:       "2021-12-12T01:00:00Z",
			expectedUnixTimestamp: "1639270800", // Unix Timestamp Seconds Dec 12 2021 01:00:00 UTC
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			time, err := time.Parse(time.RFC3339, tc.timestampString)
			if err != nil {
				t.Errorf("failed to parse expectedTime string for test case: %s", err.Error())
			}
			tTime := Time(time)
			marshaledUnixTimestamp, err := json.Marshal(&tTime)
			if err != nil {
				t.Errorf("failed to unmarshal Time from timestampString: %s", err.Error())
			}
			if string(marshaledUnixTimestamp) != tc.expectedUnixTimestamp {
				t.Errorf("marshaledUnixTimestamp did not match expected value: got - %v expected - %v", string(marshaledUnixTimestamp), tc.expectedUnixTimestamp)
			}
		})
	}
}

func TestUnmarshalJSON(t *testing.T) {
	type testCase struct {
		name               string
		timestampString    string
		expectedTimeString string
	}
	testCases := []testCase{
		{
			name:               "GIVEN a valid unix timestamp in seconds EXPECT success",
			timestampString:    "1639270800", // Unix Timestamp Seconds Dec 12 2021 01:00:00 UTC
			expectedTimeString: "2021-12-12T01:00:00Z",
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			expectedTime, err := time.Parse(time.RFC3339, tc.expectedTimeString)
			if err != nil {
				t.Errorf("failed to parse expectedTime string for test case: %s", err.Error())
			}
			var unmarshaledTime Time
			err = json.Unmarshal([]byte(tc.timestampString), &unmarshaledTime)
			if err != nil {
				t.Errorf("failed to unmarshal Time from timestampString: %s", err.Error())
			}
			if !expectedTime.Equal(time.Time(unmarshaledTime)) {
				t.Errorf("unmarshaledTime did not match expected value: got - %v expected - %v", unmarshaledTime, expectedTime)
			}
		})
	}
}
