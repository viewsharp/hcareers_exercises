package string_unpacking

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestUnpackString_Base(t *testing.T) {
	testCases := []struct {
		input          string
		expectedOutput string
		expectedError  error
	}{
		{input: "a4bc2d5e", expectedOutput: "aaaabccddddde"},
		{input: "abcd", expectedOutput: "abcd"},
		{input: "3abc", expectedError: IncorrectStringError},
		{input: "45", expectedError: IncorrectStringError},
		{input: "aaa10b", expectedError: IncorrectStringError},
		{input: "aaa0b", expectedOutput: "aab"},
		{input: "", expectedOutput: ""},
		{input: "d\n5abc", expectedOutput: "d\n\n\n\n\nabc"},
	}

	for _, testCase := range testCases {
		t.Run(testCase.input, func(t *testing.T) {
			output, err := UnpackString(testCase.input)
			require.Equal(t, testCase.expectedOutput, output)
			require.Equal(t, testCase.expectedError, err)
		})
	}
}

func TestUnpackString_Escaping(t *testing.T) {
	testCases := []struct {
		input          string
		expectedOutput string
		expectedError  error
	}{
		{input: `qwe\4\5`, expectedOutput: "qwe45"},
		{input: `qwe\45`, expectedOutput: "qwe44444"},
		{input: `qwe\\5`, expectedOutput: `qwe\\\\\`},
		{input: `qw\ne`, expectedError: IncorrectStringError},
	}

	for _, testCase := range testCases {
		t.Run(testCase.input, func(t *testing.T) {
			output, err := UnpackString(testCase.input)
			require.Equal(t, testCase.expectedOutput, output)
			require.Equal(t, testCase.expectedError, err)
		})
	}
}
