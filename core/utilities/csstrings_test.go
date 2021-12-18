package utilities

import (
	"strings"
	"testing"
)

func TestNewCSString(t *testing.T) {
	type testCase struct {
		name             string
		inputValues      []string
		expectedCSString CSString
	}
	testCases := []testCase{
		{
			name:             "GIVEN a nil slice of input values EXPECT an empty CSString to be returned",
			inputValues:      nil,
			expectedCSString: "",
		},
		{
			name:             "GIVEN an empty slice of input values EXPECT an empty CSString to be returned",
			inputValues:      []string{},
			expectedCSString: "",
		},
		{
			name: "GIVEN a slice of input values EXPECT a CSString with the input values provided to be returned",
			inputValues: []string{
				"item1",
				"item2",
				"item3",
				"item4",
				"item5",
			},
			expectedCSString: "item1,item2,item3,item4,item5",
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			csString := NewCSString(tc.inputValues)
			if csString != tc.expectedCSString {
				t.Errorf("\tcsString value not expected: got - %s expected - %s", csString, tc.expectedCSString)
			}
		})
	}
}

func TestAdd(t *testing.T) {
	type testCase struct {
		name                  string
		initialCSString       CSString
		valuesToAdd           []string
		expectedNumberOfItems int
	}
	testCases := []testCase{
		{
			name:            "GIVEN three items to add to an empty CSString EXPECT three items to be added",
			initialCSString: "",
			valuesToAdd: []string{
				"item1",
				"item2",
				"item3",
			},
			expectedNumberOfItems: 3,
		},
		{
			name:            "GIVEN three items to add to an empty CSString EXPECT three items to be added",
			initialCSString: "initial_item",
			valuesToAdd: []string{
				"item1",
				"item2",
				"item3",
			},
			expectedNumberOfItems: 4,
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			var csString CSString = tc.initialCSString
			for _, v := range tc.valuesToAdd {
				csString.Add(v)
			}
			numberOfItems := len(strings.Split(string(csString), ","))
			if tc.expectedNumberOfItems != numberOfItems {
				t.Errorf("\tnumberOfItems is not expected value: got - %d expected - %d", numberOfItems, tc.expectedNumberOfItems)
			}
		})
	}
}

func TestGet(t *testing.T) {
	type testCase struct {
		name          string
		csString      CSString
		index         int
		expectedValue string
		expectedError error
	}
	testCases := []testCase{
		{
			name:          "GIVEN a valid index EXPECT the value at that index to be returned",
			csString:      "item1,items2,item3,item4,item5",
			index:         2,
			expectedValue: "item3",
		},
		{
			name:          "GIVEN a negative index EXPECT appropriate error to be returned",
			csString:      "item1,items2,item3,item4,item5",
			index:         -1,
			expectedValue: "",
			expectedError: ErrorCSStringIndexNegative,
		},
		{
			name:          "GIVEN an index out of bound for the CSString EXPECT appropriate error to be returned",
			csString:      "item1,items2,item3,item4,item5",
			index:         100,
			expectedValue: "",
			expectedError: ErrorCSStringIndexOutOfBounds,
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			value, err := tc.csString.Get(tc.index)
			if err != nil {
				if err != tc.expectedError {
					t.Errorf("\terror encountered that was not expected: %s", err)
				}
			} else if value != tc.expectedValue {
				t.Errorf("\tvalue not expected value: got - %s expected - %s", value, tc.expectedValue)
			}
		})
	}
}

func TestContains(t *testing.T) {
	type testCase struct {
		name          string
		csString      CSString
		itemToCheck   string
		expectedIndex int
		expectedFound bool
	}
	testCases := []testCase{
		{
			name:          "GIVEN an item in a CSString EXPECT the index to be returned and found to be true",
			csString:      "item1,item2,item3",
			itemToCheck:   "item2",
			expectedIndex: 1,
			expectedFound: true,
		},
		{
			name:          "GIVEN an item not in a CSString EXPECT the index to be -1 and found to be false",
			csString:      "item1,item2,item3",
			itemToCheck:   "other",
			expectedIndex: -1,
			expectedFound: false,
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			csString := tc.csString
			index, found := csString.Contains(tc.itemToCheck)
			if index != tc.expectedIndex {
				t.Errorf("\tindex was not expected value: got - %v expected -%v", index, tc.expectedIndex)
			}
			if found != tc.expectedFound {
				t.Errorf("\tfound was not expected value: got - %v expected -%v", found, tc.expectedFound)
			}
		})
	}
}

func TestContainsCaseInsensitive(t *testing.T) {
	type testCase struct {
		name          string
		csString      CSString
		itemToCheck   string
		expectedIndex int
		expectedFound bool
	}
	testCases := []testCase{
		{
			name:          "GIVEN an item in a CSString EXPECT the index to be returned and found to be true",
			csString:      "item1,item2,item3",
			itemToCheck:   "item2",
			expectedIndex: 1,
			expectedFound: true,
		},
		{
			name:          "GIVEN an item in a CSString but of a differing case EXPECT the index to be returned and found to be true",
			csString:      "item1,item2,item3",
			itemToCheck:   "Item2",
			expectedIndex: 1,
			expectedFound: true,
		},
		{
			name:          "GIVEN an item not in a CSString EXPECT the index to be -1 and found to be false",
			csString:      "item1,item2,item3",
			itemToCheck:   "other",
			expectedIndex: -1,
			expectedFound: false,
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			index, found := tc.csString.ContainsCaseInsensitive(tc.itemToCheck)
			if index != tc.expectedIndex {
				t.Errorf("\tindex was not expected value: got - %v expected -%v", index, tc.expectedIndex)
			}
			if found != tc.expectedFound {
				t.Errorf("\tfound was not expected value: got - %v expected -%v", found, tc.expectedFound)
			}
		})
	}
}

func TestToSlice(t *testing.T) {
	type testCase struct {
		name          string
		csString      CSString
		expectedItems []string
		shouldBeNil   bool
	}
	testCases := []testCase{
		{
			name:          "GIVEN an empty CSString EXPECT a nil slice to be returned",
			csString:      "",
			expectedItems: []string{},
			shouldBeNil:   true,
		},
		{
			name:     "GIVEN a CSString with three items EXPECT a slice of string with the same three items in the same order will be returned",
			csString: "item1,item2,item3",
			expectedItems: []string{
				"item1",
				"item2",
				"item3",
			},
			shouldBeNil: false,
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ss := tc.csString.ToSlice()
			numSS := len(ss)
			numExpectedItems := len(tc.expectedItems)
			if numExpectedItems != numSS {
				t.Errorf("\tnumber of items in slice not expected value: got - %d expected - %d", numSS, numExpectedItems)
				return
			}
			if tc.shouldBeNil && ss != nil {
				t.Errorf("\tslice returned should be nil but its not: %v", ss)
			} else {
				for i := 0; i < len(ss); i++ {
					if ss[i] != tc.expectedItems[i] {
						t.Errorf("\titem at index %d in slice not expected value: got - %s expected - %s", i, ss[i], tc.expectedItems[i])
					}
				}
			}
		})
	}
}

func TestItemCount(t *testing.T) {
	type testCase struct {
		name              string
		csString          CSString
		expectedItemCount int
	}
	testCases := []testCase{
		{
			name:              "GIVEN an empty CSString EXPECT zero to be returned",
			csString:          "",
			expectedItemCount: 0,
		},
		{
			name:              "GIVEN a CSString with one value EXPECT one to be returned",
			csString:          "item1",
			expectedItemCount: 1,
		},
		{
			name:              "GIVEN a CSString with two value EXPECT two to be returned",
			csString:          "item1,item2",
			expectedItemCount: 2,
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			count := tc.csString.ItemCount()
			if count != tc.expectedItemCount {
				t.Errorf("\tcount not expected value: got - %d expected - %d", count, tc.expectedItemCount)
			}
		})
	}
}

func TestMarshalJSON(t *testing.T) {
	type testCase struct {
		name          string
		csString      CSString
		expectedJSON  string
		expectedError bool
	}
	testCases := []testCase{
		{
			name:         "GIVEN a valid CSString EXPECT an array of strings in JSON format",
			csString:     "item1,item2,item3",
			expectedJSON: `["item1","item2","item3"]`,
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			data, err := tc.csString.MarshalJSON()
			if err != nil {
				if !tc.expectedError {
					t.Errorf("\t an unexpected error occurred: %s", err.Error())
				}
				return
			}
			dataString := string(data)
			if dataString != tc.expectedJSON {
				t.Errorf("\tjson data not expected: got - %s  expected: - %s", dataString, tc.expectedJSON)
			}
		})
	}
}

func TestUnmarshalJSON(t *testing.T) {
	type testCase struct {
		name             string
		inputData        string
		expectedCSString CSString
		expectedError    bool
	}
	testCases := []testCase{
		{
			name:             "GIVEN a JSON array of strings EXPECT the CSString to be correct",
			inputData:        `["item1","item2","item3"]`,
			expectedCSString: "item1,item2,item3",
		},
		{
			name:          "GIVEN a JSON array of strings and other types EXPECT an error to be returned",
			inputData:     `["item1","item2",4, false]`,
			expectedError: true,
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			var csString CSString
			err := csString.UnmarshalJSON([]byte(tc.inputData))
			if err != nil {
				if !tc.expectedError {
					t.Errorf("\t an unexpected error occurred: %s", err.Error())
				}
				return
			}
			if csString != tc.expectedCSString {
				t.Errorf("\tcsString value not expected: got - %s expected - %s", csString, tc.expectedCSString)
			}
		})
	}
}
