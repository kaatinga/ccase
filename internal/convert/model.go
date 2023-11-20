package convert

import (
	"strings"
	"unicode"
)

const (
	lowerSnakeCase Case = 1 << iota
	upperSnakeCase
	lowerCamelCase
	upperCamelCase
	Ignore
)

type Case byte

func (c Case) String() string {
	switch c {
	case lowerCamelCase:
		return "Lower Camel Case"
	case upperCamelCase:
		return "Upper Camel Case"
	case lowerSnakeCase:
		return "Lower Snake Case"
	case upperSnakeCase:
		return "Upper Snake Case"
	}

	return ""
}

func String(input string) (Case, []string) {
	if input == "" {
		return Ignore, nil
	}

	inputChars := []rune(input)

	// File names that begin with “.” or “_” are ignored by the go tool//
	if inputChars[0] == '_' || inputChars[0] == '.' {
		return Ignore, nil
	}

	var upperCase bool
	if unicode.IsUpper(inputChars[0]) {
		upperCase = true
	}

	for _, char := range input {
		if char == '_' {
			return getCase(upperCase, upperSnakeCase), splitSnakeCase(inputChars)
		}
	}

	return getCase(upperCase, upperCamelCase), splitCamelCase(inputChars)
}

func getCase(isUpper bool, preliminaryCase Case) Case {
	switch preliminaryCase {
	case lowerSnakeCase, upperSnakeCase:
		if isUpper {
			return upperSnakeCase
		} else {
			return lowerSnakeCase
		}
	case lowerCamelCase, upperCamelCase:
		if isUpper {
			return upperCamelCase
		} else {
			return lowerCamelCase
		}
	default:
		return Ignore
	}
}

func splitCamelCase(input []rune) []string {
	var currentWordFirstIndex int
	var lowerCaseFound bool
	var words []string
	for i, char := range input {
		if unicode.IsLower(char) {
			lowerCaseFound = true
		} else {
			if lowerCaseFound {
				words = append(words, strings.ToLower(string(input[currentWordFirstIndex:i])))
				currentWordFirstIndex = i
			}
		}
	}

	return append(words, strings.ToLower(string(input[currentWordFirstIndex:])))
}

func splitSnakeCase(input []rune) []string {
	output := strings.Split(string(input), "_")
	for i := range output {
		output[i] = strings.ToLower(output[i])
	}
	return output
}
