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
