package string_unpacking

import (
	"strconv"
	"strings"
	"unicode/utf8"
)

func PackString(input string) string {
	if len(input) < 2 {
		return input
	}

	charGroup, _ := utf8.DecodeRuneInString(input)
	charCount := 1

	outputBuilder := strings.Builder{}

	for _, char := range input[1:] {
		if char == charGroup {
			charCount += 1
			continue
		}

		outputBuilder.WriteRune(charGroup)
		if charCount != 1 {
			outputBuilder.WriteString(strconv.Itoa(charCount))
		}

		charGroup = char
		charCount = 1
	}

	outputBuilder.WriteRune(charGroup)
	if charCount != 1 {
		outputBuilder.WriteString(strconv.Itoa(charCount))
	}

	return outputBuilder.String()
}
