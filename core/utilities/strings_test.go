package utilities

import "testing"

type testCase struct {
	expectedOutput string
	input          string
	name           string
}

func TestUpperCaseFirstChar(t *testing.T) {
	testCases := []testCase{
		{
			expectedOutput: "",
			input:          "",
			name:           "empty string",
		},
		{
			expectedOutput: "ThingOrTwo",
			input:          "thingOrTwo",
			name:           "string test",
		},
	}
	for _, test := range testCases {
		output := UpperCaseFirstChar(test.input)
		if output != test.expectedOutput {
			t.Errorf("\t%s test failed: output not expected: (expected: %s) (actual: %s)", test.name, test.expectedOutput, output)
		}
	}
}

func TestLowerCaseFirstChar(t *testing.T) {
	testCases := []testCase{
		{
			expectedOutput: "",
			input:          "",
			name:           "enpty string",
		},
		{
			expectedOutput: "thingOrTwo",
			input:          "ThingOrTwo",
			name:           "string test",
		},
	}
	for _, test := range testCases {
		output := LowerCaseFirstChar(test.input)
		if output != test.expectedOutput {
			t.Errorf("\t%s test failed: output not expected: (expected: %s) (actual: %s)", test.name, test.expectedOutput, output)
		}
	}
}
