package string_unpacking

import (
	"errors"
	"strconv"
	"strings"
	"unicode"
)

var IncorrectStringError = errors.New("incorrect string")

func UnpackString(input string) (string, error) {
	var prevChar rune
	prevCharExists := false
	escaped := false

	outputBuilder := strings.Builder{}

	for _, char := range input {
		if escaped {
			if !unicode.IsDigit(char) && char != '\\' {
				return "", IncorrectStringError
			}

			prevChar = char
			prevCharExists = true
			escaped = false
			continue
		}

		if unicode.IsDigit(char) {
			if !prevCharExists {
				return "", IncorrectStringError
			}

			n, err := strconv.Atoi(string(char))
			if err != nil {
				panic(err)
			}

			for i := 0; i < n; i++ {
				outputBuilder.WriteRune(prevChar)
			}

			prevCharExists = false
			continue
		}

		if prevCharExists {
			outputBuilder.WriteRune(prevChar)
		}

		if char == '\\' {
			escaped = true
		} else {
			prevChar = char
			prevCharExists = true
		}
	}

	if prevCharExists {
		outputBuilder.WriteRune(prevChar)
	}

	return outputBuilder.String(), nil
}
