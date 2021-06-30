package utilities

import "testing"

type hashTestCase struct {
	ExpectedOutput, Input, Name string
}

func TestSHA256(t *testing.T) {
	tests := []hashTestCase{
		{
			ExpectedOutput: "9f86d081884c7d659a2feaa0c55ad015a3bf4f1b2b0b822cd15d6c15b0f00a08",
			Input:          "test",
			Name:           "Test SHA256 Hash 1",
		}, {

			ExpectedOutput: "552ba4939a4136a54a3216425d80abe2e6b7360aea8e819e158d633150e216e0",
			Input:          "Lorem ipsum dolor sit amet, consectetur adipiscing elit. Maecenas sit amet ipsum turpis. Integer ultrices nibh at urna ullamcorper, at finibus velit facilisis. Fusce sed neque in nulla volutpat placerat in ut sem. Morbi pulvinar turpis varius sollicitudin bibendum. Donec scelerisque mi vel justo imperdiet ultricies. Etiam pretium arcu nisl, eget consequat urna commodo ut. Vestibulum eget urna arcu.",
			Name:           "Test SHA256 Hash 2",
		}, {

			ExpectedOutput: "33eb0576bd8ecb5317c08dfaa4c3c2853ac740b23c248ef65959c4fe12eca4cf",
			Input:          "This is a test sentence.",
			Name:           "Test SHA256 Hash 3",
		}, {

			ExpectedOutput: "a773cf70d0eca27342cd8ea51d762a7642fc664d63b14a142296650258587d96",
			Input:          "98730149387094892t873194urhfghgui9re783y1hreuofjgueb9e8fvywei2d9fvu8yb",
			Name:           "Test SHA256 Hash 4",
		}, {

			ExpectedOutput: "839e7749d46ee728149bfa5807c9b023b4313b1c78448eac5d45de9b3f458494",
			Input:          "!@#$%^&*()+}{\":?><,./;'[]=-0987654321`~",
			Name:           "Test SHA256 Hash 5",
		},
	}

	for _, test := range tests {
		output, err := SHA256(test.Input)
		if err != nil {
			t.Error("failed to hash input", test.Input, err)
		}
		if output != test.ExpectedOutput {
			t.Error("hash did not match expected value", test.ExpectedOutput, output)
		}

	}
}

func TestSHA512(t *testing.T) {
	tests := []hashTestCase{
		{
			ExpectedOutput: "ee26b0dd4af7e749aa1a8ee3c10ae9923f618980772e473f8819a5d4940e0db27ac185f8a0e1d5f84f88bc887fd67b143732c304cc5fa9ad8e6f57f50028a8ff",
			Input:          "test",
			Name:           "Test SHA512 Hash 1",
		}, {

			ExpectedOutput: "77cc083c93d0d08abe852b75b984df91f42039190fc8ecdd2138a227df2a577d41b259f3660fa97ee0b5f858510ea5fe8d9e4a1311c334fffa77d7455e59c31a",
			Input:          "Lorem ipsum dolor sit amet, consectetur adipiscing elit. Maecenas sit amet ipsum turpis. Integer ultrices nibh at urna ullamcorper, at finibus velit facilisis. Fusce sed neque in nulla volutpat placerat in ut sem. Morbi pulvinar turpis varius sollicitudin bibendum. Donec scelerisque mi vel justo imperdiet ultricies. Etiam pretium arcu nisl, eget consequat urna commodo ut. Vestibulum eget urna arcu.",
			Name:           "Test SHA512 Hash 2",
		}, {

			ExpectedOutput: "04a9a680a8776339f12133a55d677c9846131f4ae9f19e452ba890db70460aabe369a0003ac3e498bec13d384e6209386c96a3825af81c13654fd15ced3a123a",
			Input:          "This is a test sentence.",
			Name:           "Test SHA512 Hash 3",
		}, {

			ExpectedOutput: "9880f6b1b5e69f6d9586cba756980fdf3531a46a3dc39c9455a989ab7876f5b0ca25ebbe746dc7a091a13b4acd078655aeb907cacb51765465bfdaf8ac690b30",
			Input:          "98730149387094892t873194urhfghgui9re783y1hreuofjgueb9e8fvywei2d9fvu8yb",
			Name:           "Test SHA512 Hash 4",
		}, {

			ExpectedOutput: "84ef172117cb1ec3a3b5c598c6e5ab53f5eb17d495f397ee83124bd62fc7fbc1e48246ca9446b5efad551390cce589921628bdda66380b9520ea6a83eb0c7c8b",
			Input:          "!@#$%^&*()+}{\":?><,./;'[]=-0987654321`~",
			Name:           "Test SHA512 Hash 5",
		},
	}

	for _, test := range tests {
		output, err := SHA512(test.Input)
		if err != nil {
			t.Error("failed to hash input", test.Input, err)
		}
		if output != test.ExpectedOutput {
			t.Error("hash did not match expected value", test.ExpectedOutput, output)
		}

	}
}
