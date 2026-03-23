package go_strings

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestReverseString_다양한입력(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{"영문", "hello", "olleh"},
		{"한글", "안녕하세요", "요세하녕안"},
		{"빈문자열", "", ""},
		{"한글자", "a", "a"},
		{"회문", "racecar", "racecar"},
		{"이모지", "hello🌍", "🌍olleh"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.expected, ReverseString(tt.input))
		})
	}
}

func TestCamelToSnake_변환(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{"camelCase", "camelCase", "camel_case"},
		{"PascalCase", "PascalCase", "pascal_case"},
		{"단일소문자", "hello", "hello"},
		{"여러대문자", "HTTPServer", "h_t_t_p_server"},
		{"빈문자열", "", ""},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.expected, CamelToSnake(tt.input))
		})
	}
}

func TestSnakeToCamel_변환(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{"기본", "snake_case", "snakeCase"},
		{"세단어", "hello_world_go", "helloWorldGo"},
		{"단일단어", "hello", "hello"},
		{"빈문자열", "", ""},
		{"앞뒤언더스코어", "_hello_world_", "helloWorld"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.expected, SnakeToCamel(tt.input))
		})
	}
}

func TestTruncateWithEllipsis_잘라내기(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		maxLen   int
		expected string
	}{
		{"짧은문자열", "hi", 10, "hi"},
		{"정확히같은길이", "hello", 5, "hello"},
		{"잘라내기", "hello world", 8, "hello..."},
		{"매우짧은최대길이", "hello", 3, "hel"},
		{"한글잘라내기", "안녕하세요 세계입니다", 7, "안녕하세..."},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.expected, TruncateWithEllipsis(tt.input, tt.maxLen))
		})
	}
}

func TestCountWords_단어수세기(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected int
	}{
		{"일반문장", "hello world", 2},
		{"여러공백", "  hello   world  ", 2},
		{"빈문자열", "", 0},
		{"탭포함", "hello\tworld\ngo", 3},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.expected, CountWords(tt.input))
		})
	}
}

func TestIsPalindrome_회문검사(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected bool
	}{
		{"영문회문", "racecar", true},
		{"대소문자무시", "RaceCar", true},
		{"공백과특수문자무시", "A man, a plan, a canal: Panama", true},
		{"회문아님", "hello", false},
		{"빈문자열", "", true},
		{"숫자포함회문", "12321", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.expected, IsPalindrome(tt.input))
		})
	}
}
