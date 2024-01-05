// Package convert provides a function to convert a string lower_snake_case.
package convert

import (
	"strings"
	"unicode"
)

type Case uint16

const unset Case = 0

const (
	lowerSnakeCase Case = 1 << iota
	upperSnakeCase
	lowerCamelCase
	upperCamelCase
	lowerKebabCase
	upperKebabCase
	inconsistentCase
	IsNotDotGo
	Ignore
)

func (c Case) IsMixedCase() bool {
	var found bool
	for i := Case(1); i < IsNotDotGo; i++ {
		if c&i == i {
			if found {
				return true
			}
			found = true
		}
	}

	return false
}

func (c Case) IsKebabCase() bool {
	return c&lowerKebabCase == lowerKebabCase || c&upperKebabCase == upperKebabCase
}

func (c Case) IsSnakeCase() bool {
	return c&lowerSnakeCase == lowerSnakeCase || c&upperSnakeCase == upperSnakeCase
}

func (c Case) String() string {
	if c&inconsistentCase == inconsistentCase {
		return "Inconsistent Case"
	}

	if c.IsMixedCase() {
		return "Mixed Case"
	}

	switch c {
	case lowerCamelCase:
		return "Lower Camel Case"
	case upperCamelCase:
		return "Upper Camel Case"
	case lowerSnakeCase:
		return "Lower Snake Case"
	case upperSnakeCase:
		return "Upper Snake Case"
	case lowerKebabCase:
		return "Lower Kebab Case"
	case upperKebabCase:
		return "Upper Kebab Case"
	default:
		return "Undetermined Case"
	}
}

const dotGoExtension = ".go"

func String(inputChars []rune) (Case, []string) {
	// File names that begin with “.” or “_” are ignored by the go tool
	if len(inputChars) == 0 || inputChars[0] == '_' || inputChars[0] == '.' {
		return Ignore, nil
	}

	if !isDotGoExtension(inputChars) {
		return IsNotDotGo, nil
	}

	inputChars = inputChars[:len(inputChars)-len(dotGoExtension)]

	return split(inputChars)
}

func (c Case) IsUpperCase() bool {
	for i := Case(2); i < inconsistentCase; i <<= 2 {
		if c&i == i {
			return true
		}
	}

	return false
}

func isDotGoExtension(inputChars []rune) bool {
	if len(inputChars) < len(dotGoExtension) {
		return false
	}

	for i := 1; i <= len(dotGoExtension); i++ {
		if inputChars[len(inputChars)-i]|0x20 != rune(dotGoExtension[len(dotGoExtension)-i]) {
			return false
		}
	}

	return true
}

func getCase(isUpper bool, preliminaryCase Case) Case {
	if preliminaryCase > upperKebabCase {
		return preliminaryCase
	}

	if isUpper {
		if !preliminaryCase.IsUpperCase() {
			return preliminaryCase << 1
		}
	} else {
		if preliminaryCase.IsUpperCase() {
			return preliminaryCase >> 1
		}
	}

	return preliminaryCase
}

func split(input []rune) (c Case, words []string) {
	var theFirstCharIsUpper bool
	if unicode.IsUpper(input[0]) {
		theFirstCharIsUpper = true
	}

	var hasUpper = theFirstCharIsUpper
	var currentWordFirstIndex int
	var abbreviationFound bool
	var lowerCaseFound bool
	for i, char := range input {
		switch {
		case char == ' ' || char == '-' || char == '_':
			lowerCaseFound = false
			word := strings.ToLower(string(input[currentWordFirstIndex:i]))
			if word != "" {
				words = append(words, word)
			}
			currentWordFirstIndex = i + 1
		case unicode.IsLower(char):
			lowerCaseFound = true
			if abbreviationFound {
				abbreviationFound = false
				word := strings.ToLower(string(input[currentWordFirstIndex : i-1]))
				if word != "" {
					words = append(words, word)
				}
				currentWordFirstIndex = i - 1
			}
		case unicode.IsUpper(char):
			hasUpper = true
			c = updateUpperCase(hasUpper, theFirstCharIsUpper, c)

			if lowerCaseFound {
				word := strings.ToLower(string(input[currentWordFirstIndex:i]))
				if word != "" {
					words = append(words, word)
				}
				currentWordFirstIndex = i
			} else if i > 0 {
				abbreviationFound = true
			}
			lowerCaseFound = false
		}

		c = updateCaseWithSeparator(char, hasUpper, theFirstCharIsUpper, c)
	}

	word := strings.ToLower(string(input[currentWordFirstIndex:]))
	if word != "" {
		words = append(words, strings.ToLower(string(input[currentWordFirstIndex:])))
	}

	if len(words) == 0 {
		return c, []string{""}
	}

	c = updateCaseIfUnset(c, theFirstCharIsUpper)

	return c, words
}

func updateCaseIfUnset(c Case, theFirstCharIsUpper bool) Case {
	if c == unset {
		c = getCase(theFirstCharIsUpper, upperCamelCase)
	}

	return c
}

func updateUpperCase(hasUpper, theFirstCharIsUpper bool, c Case) Case {
	if c.IsKebabCase() {
		c = setCaseCheckingInconsistency(hasUpper, theFirstCharIsUpper, c, upperKebabCase)
	}
	if c.IsSnakeCase() {
		c = setCaseCheckingInconsistency(hasUpper, theFirstCharIsUpper, c, upperSnakeCase)
	}

	return c
}

func updateCaseWithSeparator(separator rune, hasUpper, theFirstCharIsUpper bool, c Case) Case {
	switch separator {
	case '_':
		c = setCaseCheckingInconsistency(hasUpper, theFirstCharIsUpper, c, upperSnakeCase)
	case '-':
		c = setCaseCheckingInconsistency(hasUpper, theFirstCharIsUpper, c, upperKebabCase)
	case ' ':
		c |= inconsistentCase
	}

	return c
}

func setCaseCheckingInconsistency(hasUpper, theFirstCharIsUpper bool, c Case, targetCase Case) Case {
	if hasUpper && !theFirstCharIsUpper {
		c |= inconsistentCase
	} else {
		c |= getCase(theFirstCharIsUpper, targetCase)
	}
	return c
}
