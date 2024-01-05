package convert

import (
	"fmt"
	"reflect"
	"testing"
)

func Test_split(t *testing.T) {
	tests := []struct {
		input string
		want  []string
		Case
	}{
		{"oneWord", []string{"one", "word"}, lowerCamelCase},
		{"one", []string{"one"}, lowerCamelCase},
		{"oneWord_test", []string{"one", "word", "test"}, inconsistentCase},
		{"oneWord_", []string{"one", "word"}, inconsistentCase},
		{"one_Word", []string{"one", "word"}, inconsistentCase},
		{"one_word", []string{"one", "word"}, lowerSnakeCase},
		{"OneWord", []string{"one", "word"}, upperCamelCase},
		{"oneWordTwo", []string{"one", "word", "two"}, lowerCamelCase},
		{"one Word Two", []string{"one", "word", "two"}, inconsistentCase},
		{"OneWordTwo", []string{"one", "word", "two"}, upperCamelCase},
		{"myTest", []string{"my", "test"}, lowerCamelCase},
		{"MyTest", []string{"my", "test"}, upperCamelCase},
		{"my123String", []string{"my123", "string"}, lowerCamelCase},
		{"GRPCHandler", []string{"grpc", "handler"}, upperCamelCase},
		{"GRPCHandlerTest", []string{"grpc", "handler", "test"}, upperCamelCase},
		{"HandlerGRPCTest", []string{"handler", "grpc", "test"}, upperCamelCase},
		{"Handler-GRPC-Test", []string{"handler", "grpc", "test"}, upperKebabCase},
		{"Handler_GRPC_Test", []string{"handler", "grpc", "test"}, upperSnakeCase},
		{"Handler_GRPC_Test_", []string{"handler", "grpc", "test"}, upperSnakeCase},
		{"Handler-GRPC-Test", []string{"handler", "grpc", "test"}, upperKebabCase},
		{"Handler_GRPC-Test", []string{"handler", "grpc", "test"}, upperKebabCase | upperSnakeCase},
	}
	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			c, got := split([]rune(tt.input))
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("split() = %v, want %v", got, tt.want)
			}
			if c.String() != tt.Case.String() {
				t.Logf("row case %b tt.Case %b", c, tt.Case)
				t.Errorf("split() = %v, want %v", c, tt.Case)
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
		{"-go", false},
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

func TestCase_IsUpperCase(t *testing.T) {
	tests := []struct {
		c    Case
		want bool
	}{
		{upperSnakeCase, true},
		{upperKebabCase, true},
		{upperCamelCase, true},
		{lowerSnakeCase, false},
		{lowerKebabCase, false},
		{lowerCamelCase, false},
		{lowerCamelCase | upperCamelCase, true},
	}
	for _, tt := range tests {
		t.Run(tt.c.String(), func(t *testing.T) {
			if got := tt.c.IsUpperCase(); got != tt.want {
				t.Errorf("IsUpperCase() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestCase_IsMixedCase(t *testing.T) {
	tests := []struct {
		c    Case
		want bool
	}{
		{upperSnakeCase, false},
		{upperKebabCase, false},
		{upperCamelCase, false},
		{lowerSnakeCase, false},
		{lowerKebabCase, false},
		{lowerCamelCase, false},
		{lowerCamelCase | upperCamelCase, true},
		{lowerCamelCase | upperCamelCase | lowerSnakeCase, true},
		{lowerCamelCase | upperCamelCase | lowerSnakeCase | upperSnakeCase, true},
		{inconsistentCase, false},
	}
	for _, tt := range tests {
		t.Run(tt.c.String(), func(t *testing.T) {
			if got := tt.c.IsMixedCase(); got != tt.want {
				t.Errorf("IsMixedCase() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_getCase(t *testing.T) {
	tests := []struct {
		isUpper         bool
		preliminaryCase Case
		want            Case
	}{
		{true, upperCamelCase, upperCamelCase},
		{true, lowerCamelCase, upperCamelCase},
		{true, upperSnakeCase, upperSnakeCase},
		{true, lowerSnakeCase, upperSnakeCase},
		{true, upperKebabCase, upperKebabCase},
		{true, lowerKebabCase, upperKebabCase},
		{false, upperCamelCase, lowerCamelCase},
		{false, lowerCamelCase, lowerCamelCase},
		{false, upperSnakeCase, lowerSnakeCase},
		{false, lowerSnakeCase, lowerSnakeCase},
		{false, upperKebabCase, lowerKebabCase},
		{false, lowerKebabCase, lowerKebabCase},
		{false, inconsistentCase, inconsistentCase},
	}
	for _, tt := range tests {
		t.Run(fmt.Sprintf("%s must be upper: %v", tt.preliminaryCase.String(), tt.isUpper), func(t *testing.T) {
			if got := getCase(tt.isUpper, tt.preliminaryCase); got != tt.want {
				t.Errorf("getCase() = %v, want %v", got, tt.want)
			}
		})
	}
}
