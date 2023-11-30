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
	for i := 0; i < 7; i++ {
		if c<<i == 1<<i {
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
	}

	if c&inconsistentCase == inconsistentCase {
		return "Inconsistent Case"
	}

	return "Mixed Case"
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
	if !isUpper && (preliminaryCase == upperCamelCase || preliminaryCase == upperSnakeCase) {
		return preliminaryCase >> 1
	}
	return preliminaryCase
}

func split(input []rune) (c Case, words []string) {
	var theFirstCharisUpper bool
	if unicode.IsUpper(input[0]) {
		theFirstCharisUpper = true
	}

	var hasUpper bool
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
			if c.IsKebabCase() {
				c |= upperKebabCase
			}
			if c.IsSnakeCase() {
				c |= upperSnakeCase
			}

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

		switch char {
		case '_':
			if hasUpper && !theFirstCharisUpper {
				c |= inconsistentCase
			} else {
				c |= getCase(theFirstCharisUpper, upperSnakeCase)
			}
		case '-':
			if hasUpper && !theFirstCharisUpper {
				c |= inconsistentCase
			} else {
				c |= getCase(theFirstCharisUpper, upperKebabCase)
			}
		case ' ':
			c |= inconsistentCase
		}
	}

	word := strings.ToLower(string(input[currentWordFirstIndex:]))
	if word != "" {
		words = append(words, strings.ToLower(string(input[currentWordFirstIndex:])))
	}

	if len(words) == 0 {
		return c, []string{""}
	}

	if c == unset {
		c = getCase(theFirstCharisUpper, upperCamelCase)
	}

	return c, words
}
