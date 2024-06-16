package string_unpacking

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestPackString(t *testing.T) {
	testCases := []struct {
		input          string
		expectedOutput string
	}{
		{input: "aaaabccddddde", expectedOutput: "a4bc2d5e"},
		{input: "abcd", expectedOutput: "abcd"},
	}

	for _, testCase := range testCases {
		t.Run(testCase.input, func(t *testing.T) {
			output := PackString(testCase.input)
			require.Equal(t, testCase.expectedOutput, output)
		})
	}
}
