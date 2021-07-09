package normalization

import (
	"testing"
	"time"
)

func TestReadBoolValue(t *testing.T) {
	testString := "yes"
	testInt8 := int8(1)
	testUint8 := uint8(1)
	testFloat32 := float32(1)
	testBool := true
	var testNilString *string
	var testNilInt8 *int8
	var testNilUint8 *uint8
	var testNilFloat32 *float32
	var testNilBool *bool
	type testArgs struct {
		DefaultToFalse bool
		Value          interface{}
	}
	testCases := []struct {
		ExpectedError  bool
		ExpectedOutput bool
		Input          testArgs
		Name           string
	}{
		{
			ExpectedError:  false,
			ExpectedOutput: true,
			Input:          testArgs{false, int8(1)},
			Name:           "testing int8 value of 1",
		},
		{
			ExpectedError:  false,
			ExpectedOutput: false,
			Input:          testArgs{false, int8(0)},
			Name:           "testing int8 value of 0",
		},
		{
			ExpectedError:  false,
			ExpectedOutput: true,
			Input:          testArgs{false, int16(1)},
			Name:           "testing int16 value of 1",
		},
		{
			ExpectedError:  false,
			ExpectedOutput: false,
			Input:          testArgs{false, int16(0)},
			Name:           "testing int16 value of 0",
		},
		{
			ExpectedError:  false,
			ExpectedOutput: true,
			Input:          testArgs{false, int(1)},
			Name:           "testing int value of 1",
		},
		{
			ExpectedError:  false,
			ExpectedOutput: false,
			Input:          testArgs{false, int(0)},
			Name:           "testing int value of 0",
		},
		{
			ExpectedError:  false,
			ExpectedOutput: true,
			Input:          testArgs{false, int32(1)},
			Name:           "testing int32 value of 1",
		},
		{
			ExpectedError:  false,
			ExpectedOutput: false,
			Input:          testArgs{false, int32(0)},
			Name:           "testing int32 value of 0",
		},
		{
			ExpectedError:  false,
			ExpectedOutput: true,
			Input:          testArgs{false, int64(1)},
			Name:           "testing int64 value of 1",
		},
		{
			ExpectedError:  false,
			ExpectedOutput: false,
			Input:          testArgs{false, int64(0)},
			Name:           "testing int64 value of 0",
		},
		{
			ExpectedError:  false,
			ExpectedOutput: true,
			Input:          testArgs{false, uint8(1)},
			Name:           "testing uint8 value of 1",
		},
		{
			ExpectedError:  false,
			ExpectedOutput: false,
			Input:          testArgs{false, uint8(0)},
			Name:           "testing uint8 value of 0",
		},
		{
			ExpectedError:  false,
			ExpectedOutput: true,
			Input:          testArgs{false, uint16(1)},
			Name:           "testing uint16 value of 1",
		},
		{
			ExpectedError:  false,
			ExpectedOutput: false,
			Input:          testArgs{false, uint16(0)},
			Name:           "testing uint16 value of 0",
		},
		{
			ExpectedError:  false,
			ExpectedOutput: true,
			Input:          testArgs{false, uint(1)},
			Name:           "testing uint value of 1",
		},
		{
			ExpectedError:  false,
			ExpectedOutput: false,
			Input:          testArgs{false, uint(0)},
			Name:           "testing uint value of 0",
		},
		{
			ExpectedError:  false,
			ExpectedOutput: true,
			Input:          testArgs{false, uint32(1)},
			Name:           "testing uint32 value of 1",
		},
		{
			ExpectedError:  false,
			ExpectedOutput: false,
			Input:          testArgs{false, uint32(0)},
			Name:           "testing uint32 value of 0",
		},
		{
			ExpectedError:  false,
			ExpectedOutput: true,
			Input:          testArgs{false, uint64(1)},
			Name:           "testing uint64 value of 1",
		},
		{
			ExpectedError:  false,
			ExpectedOutput: false,
			Input:          testArgs{false, uint64(0)},
			Name:           "testing uint64 value of 0",
		},
		{
			ExpectedError:  false,
			ExpectedOutput: true,
			Input:          testArgs{false, float32(1)},
			Name:           "testing float32 value of 1",
		},
		{
			ExpectedError:  false,
			ExpectedOutput: false,
			Input:          testArgs{false, float32(0)},
			Name:           "testing float32 value of 0",
		},
		{
			ExpectedError:  false,
			ExpectedOutput: true,
			Input:          testArgs{false, float64(1)},
			Name:           "testing float64 value of 1",
		},
		{
			ExpectedError:  false,
			ExpectedOutput: false,
			Input:          testArgs{false, float64(0)},
			Name:           "testing float64 value of 0",
		},
		{
			ExpectedError:  false,
			ExpectedOutput: true,
			Input:          testArgs{false, "yes"},
			Name:           "testing string value of yes",
		},
		{
			ExpectedError:  false,
			ExpectedOutput: true,
			Input:          testArgs{false, "y"},
			Name:           "testing string value of y",
		},
		{
			ExpectedError:  false,
			ExpectedOutput: true,
			Input:          testArgs{false, "t"},
			Name:           "testing string value of t",
		},
		{
			ExpectedError:  false,
			ExpectedOutput: true,
			Input:          testArgs{false, "true"},
			Name:           "testing string value of true",
		},
		{
			ExpectedError:  false,
			ExpectedOutput: true,
			Input:          testArgs{false, "1"},
			Name:           "testing string value of 1",
		},
		{
			ExpectedError:  false,
			ExpectedOutput: false,
			Input:          testArgs{false, "no"},
			Name:           "testing string value of n",
		},
		{
			ExpectedError:  false,
			ExpectedOutput: false,
			Input:          testArgs{false, "n"},
			Name:           "testing string value of n",
		},
		{
			ExpectedError:  false,
			ExpectedOutput: false,
			Input:          testArgs{false, "f"},
			Name:           "testing string value of f",
		},
		{
			ExpectedError:  false,
			ExpectedOutput: false,
			Input:          testArgs{false, "false"},
			Name:           "testing string value of false",
		},
		{
			ExpectedError:  false,
			ExpectedOutput: false,
			Input:          testArgs{false, "0"},
			Name:           "testing string value of 0",
		},
		{
			ExpectedError:  false,
			ExpectedOutput: true,
			Input:          testArgs{false, &testString},
			Name:           "testing *string value of yes",
		},
		{
			ExpectedError:  false,
			ExpectedOutput: true,
			Input:          testArgs{false, &testInt8},
			Name:           "testing *int8 value of 1",
		},
		{
			ExpectedError:  false,
			ExpectedOutput: true,
			Input:          testArgs{false, &testUint8},
			Name:           "testing *uint8 value of 1",
		},
		{
			ExpectedError:  false,
			ExpectedOutput: true,
			Input:          testArgs{false, &testFloat32},
			Name:           "testing *float32 value of 1",
		},
		{
			ExpectedError:  false,
			ExpectedOutput: true,
			Input:          testArgs{false, &testBool},
			Name:           "testing *bool value of true",
		},
		{
			ExpectedError:  true,
			ExpectedOutput: false,
			Input:          testArgs{false, testNilString},
			Name:           "testing nil *string value of yes",
		},
		{
			ExpectedError:  true,
			ExpectedOutput: false,
			Input:          testArgs{false, testNilInt8},
			Name:           "testing nil *int8 value of 1",
		},
		{
			ExpectedError:  true,
			ExpectedOutput: false,
			Input:          testArgs{false, testNilUint8},
			Name:           "testing nil *uint8 value of 1",
		},
		{
			ExpectedError:  true,
			ExpectedOutput: false,
			Input:          testArgs{false, testNilFloat32},
			Name:           "testing nil *float32 value of 1",
		},
		{
			ExpectedError:  true,
			ExpectedOutput: false,
			Input:          testArgs{false, testNilBool},
			Name:           "testing nil *bool value of true",
		},
		{
			ExpectedError:  true,
			ExpectedOutput: false,
			Input:          testArgs{false, time.Time{}},
			Name:           "testing nil time.Time with empty value",
		},
	}

	for _, test := range testCases {
		t.Run(test.Name, func(t *testing.T) {
			value := test.Input
			b, err := ReadBoolValue(value.Value, value.DefaultToFalse)
			if test.ExpectedError && err != nil {
				t.Logf("error returned but it was expectet: %s", err.Error())
			}
			if test.ExpectedError && err == nil {
				t.Error("expected an error to be thrown", test, test.Name)
			} else if !test.ExpectedError && err != nil {
				t.Error("failed to read value as bool", value, err, test.Name)
			} else if b != test.ExpectedOutput {
				t.Error("expected parsed bool value to be true", b, test.Name)
			}
		})
	}
}
