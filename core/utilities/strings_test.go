package utilities

import "testing"

type testCase struct {
	input          string
	expectedOutput string
}

func TestUpperCaseFirstChar(t *testing.T) {
	testCases := []testCase{
		{
			expectedOutput: "",
			input:          "",
		},
		{
			expectedOutput: "ThingOrTwo",
			input:          "thingOrTwo",
		},
	}
	for _, test := range testCases {
		output := UpperCaseFirstChar(test.input)
		if output != test.expectedOutput {
			t.Errorf("output not expected: (expected: %s) (actual: %s)", test.expectedOutput, output)
		}
	}
}

func TestLowerCaseFirstChar(t *testing.T) {
	testCases := []testCase{
		{
			expectedOutput: "",
			input:          "",
		},
		{
			expectedOutput: "thingOrTwo",
			input:          "ThingOrTwo",
		},
	}
	for _, test := range testCases {
		output := LowerCaseFirstChar(test.input)
		if output != test.expectedOutput {
			t.Errorf("output not expected: (expected: %s) (actual: %s)", test.expectedOutput, output)
		}
	}
}
