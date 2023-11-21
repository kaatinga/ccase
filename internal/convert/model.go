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
	MixedCase
	IsNotDotGo
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

var dotGoExtension = []rune(".go")

func String(input string) (Case, []string) {
	if input == "" {
		return Ignore, nil
	}

	inputChars := []rune(input)

	if !isDotGoExtension(inputChars) {
		return IsNotDotGo, nil
	}

	inputChars = inputChars[:len(inputChars)-len(dotGoExtension)]

	// File names that begin with “.” or “_” are ignored by the go tool//
	if inputChars[0] == '_' || inputChars[0] == '.' {
		return Ignore, nil
	}

	var upperCase bool
	if unicode.IsUpper(inputChars[0]) {
		upperCase = true
	}

	for _, char := range inputChars {
		if char == '_' {
			inputCase := getCase(upperCase, upperSnakeCase)
			words := splitSnakeCase(inputChars)
			var mixedWords []string
			for _, word := range words {
				camelWords := splitCamelCase([]rune(word))
				if len(camelWords) > 1 {
					inputCase = MixedCase
				}
				mixedWords = append(mixedWords, camelWords...)
			}

			if inputCase == MixedCase {
				return MixedCase, mixedWords
			}

			return inputCase, splitSnakeCase(inputChars)
		}
	}

	return getCase(upperCase, upperCamelCase), splitCamelCase(inputChars)
}

func isDotGoExtension(inputChars []rune) bool {
	if len(inputChars) < len(dotGoExtension) {
		return false
	}

	for i := 1; i < len(dotGoExtension); i++ {
		if unicode.ToLower(inputChars[len(inputChars)-i]) != dotGoExtension[len(dotGoExtension)-i] {
			return false
		}
	}

	return true
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
