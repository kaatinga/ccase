package convert

import (
	"reflect"
	"testing"
)

func Test_splitCamelCase(t *testing.T) {
	tests := []struct {
		input string
		want  []string
	}{
		{"", []string{""}},
		{"oneWord", []string{"one", "word"}},
		{"OneWord", []string{"one", "word"}},
		{"oneWordTwo", []string{"one", "word", "two"}},
		{"OneWordTwo", []string{"one", "word", "two"}},
	}
	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			if got := splitCamelCase([]rune(tt.input)); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("splitCamelCase() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_checkDotGoExtension(t *testing.T) {
	tests := []struct {
		input string
		want  bool
	}{
		{"main.go", true},
		{"main.GO", true},
		{"main.go.go", true},
		{"main.go.GO", true},
		{"main.go.go.go", true},
		{"main.txt", false},
		{"main.go.txt", false},
		{"README.md", false},
		{"", false},
	}
	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			if got := isDotGoExtension([]rune(tt.input)); got != tt.want {
				t.Errorf("isDotGoExtension() = %v, want %v", got, tt.want)
			}
		})
	}
}
