package testutilities

import "testing"

type testStruct1 struct {
	P1 string
	P2 int
	P3 *string
	P4 []float64
	p5 string
}

// fun fact the the field less embedded struct is not exported, it is non accessable,
// and if the embedded struct is exported it is. This seems obvious,
// but I thought it was interesting nonetheless
// type testStruct2 struct {
// 	P6 uint
// 	testStruct1
// }

type testStruct3 struct {
	P7  complex128
	Ts1 testStruct1
}

func TestThingsValuesAreEqual(t *testing.T) {
	testString1 := "teststring"
	testString2 := "teststring"
	testCases := []struct {
		name           string
		input1         interface{}
		input2         interface{}
		expectedOutput bool
	}{
		{
			name:           "two nils are equal",
			input1:         nil,
			input2:         nil,
			expectedOutput: true,
		},
		{
			name:           "two strings are equal",
			input1:         "test string",
			input2:         "test string",
			expectedOutput: true,
		},
		{
			name:           "two strings are not equal",
			input1:         "test string123",
			input2:         "test string",
			expectedOutput: false,
		},
		{
			name:           "two int are equal",
			input1:         1,
			input2:         1,
			expectedOutput: true,
		},
		{
			name:           "two int are not equal",
			input1:         1,
			input2:         0,
			expectedOutput: false,
		},
		{
			name:           "two float64 are equal",
			input1:         1.234,
			input2:         1.234,
			expectedOutput: true,
		},
		{
			name:           "two float64 are not equal",
			input1:         0.123,
			input2:         1.123,
			expectedOutput: false,
		},
		{
			name: "two structs without embedded structs are equal",
			input1: testStruct1{
				P1: "test",
				P2: 123,
				P3: &testString1,
				P4: []float64{1.234, 2.345, 3.456, 4.567, 5.678, 6.789, 7.890},
				p5: "1", // this field is not exported and should be ignored by the equality comparison...
			},
			input2: testStruct1{
				P1: "test",
				P2: 123,
				P3: &testString2,
				P4: []float64{1.234, 2.345, 3.456, 4.567, 5.678, 6.789, 7.890},
				p5: "2", // this field is not exported and should be ignored by the equality comparison...
			},
			expectedOutput: true,
		},
		{
			name: "two structs without embedded structs are not equal",
			input1: testStruct1{
				P1: "test",
				P2: 123,
				P3: &testString1,
				P4: []float64{1.234, 2.345, 3.456, 4.567, 5.678, 6.789, 7.890},
			},
			input2: testStruct1{
				P1: "test",
				P2: 123,
				P3: &testString2,
				P4: []float64{1.233, 2.345, 3.456, 4.567, 5.678, 6.789, 7.890},
			},
			expectedOutput: false,
		},
		{
			name: "two structs with embedded struct with field for struct are equal",
			input1: testStruct3{
				P7: 123,
				Ts1: testStruct1{
					P1: "test",
					P2: 123,
					P3: &testString1,
					P4: []float64{1.234, 2.345, 3.456, 4.567, 5.678, 6.789, 7.890},
					p5: "1", // this field is not exported and should be ignored by the equality comparison...
				},
			},
			input2: testStruct3{
				P7: 123,
				Ts1: testStruct1{
					P1: "test",
					P2: 123,
					P3: &testString2,
					P4: []float64{1.234, 2.345, 3.456, 4.567, 5.678, 6.789, 7.890},
					p5: "2", // this field is not exported and should be ignored by the equality comparison...
				},
			},
			expectedOutput: true,
		},
		{
			name: "two structs with embedded struct with field for struct are not equal",
			input1: testStruct3{
				P7: 123,
				Ts1: testStruct1{
					P1: "test",
					P2: 123,
					P3: &testString1,
					P4: []float64{1.234, 2.345, 3.456, 4.567, 5.678, 6.789, 7.890},
				},
			},
			input2: testStruct3{
				P7: 123,
				Ts1: testStruct1{
					P1: "test",
					P2: 123,
					P3: &testString2,
					P4: []float64{1.233, 2.345, 3.456, 4.567, 5.678, 6.789, 7.890},
				},
			},
			expectedOutput: false,
		},
		{
			name: "two maps are equal",
			input1: map[string]string{
				"k1": "v1",
				"k2": "v2",
				"k3": "v3",
			},
			input2: map[string]string{
				"k1": "v1",
				"k2": "v2",
				"k3": "v3",
			},
			expectedOutput: true,
		},
		{
			name: "two maps are not equal",
			input1: map[string]string{
				"k1": "v1",
				"k2": "v2",
				"k3": "v3",
			},
			input2: map[string]string{
				"k1": "v1",
				"k2": "v2",
				"k3": "v4",
			},
			expectedOutput: false,
		},
		{
			name: "two maps with non matching keys",
			input1: map[string]string{
				"k1": "v1",
				"k2": "v2",
				"k3": "v3",
			},
			input2: map[string]string{
				"k1": "v1",
				"k2": "v2",
				"k3": "v3",
				"k4": "v4",
			},
			expectedOutput: false,
		},
	}
	for _, tt := range testCases {
		t.Run(tt.name, func(t *testing.T) {
			result := Equals(tt.input1, tt.input2)
			if result.AreEqual != tt.expectedOutput {
				t.Errorf("equality comparison not expected output: got: %t - expected: %t", result.AreEqual, tt.expectedOutput)
			}
		})
	}
}
