package go_strings

import (
	"strings"
	"unicode"
)

// ReverseString reverses the given string, handling Unicode correctly.
func ReverseString(s string) string {
	runes := []rune(s)
	for i, j := 0, len(runes)-1; i < j; i, j = i+1, j-1 {
		runes[i], runes[j] = runes[j], runes[i]
	}
	return string(runes)
}

// CamelToSnake converts a camelCase or PascalCase string to snake_case.
func CamelToSnake(s string) string {
	var result strings.Builder
	for i, r := range s {
		if unicode.IsUpper(r) {
			if i > 0 {
				result.WriteRune('_')
			}
			result.WriteRune(unicode.ToLower(r))
		} else {
			result.WriteRune(r)
		}
	}
	return result.String()
}

// SnakeToCamel converts a snake_case string to camelCase.
func SnakeToCamel(s string) string {
	parts := strings.Split(s, "_")
	var result strings.Builder
	for i, part := range parts {
		if part == "" {
			continue
		}
		if i == 0 {
			result.WriteString(strings.ToLower(part))
		} else {
			result.WriteString(strings.ToUpper(part[:1]) + strings.ToLower(part[1:]))
		}
	}
	return result.String()
}

// TruncateWithEllipsis truncates a string to maxLen and appends "..." if truncated.
func TruncateWithEllipsis(s string, maxLen int) string {
	runes := []rune(s)
	if len(runes) <= maxLen {
		return s
	}
	if maxLen <= 3 {
		return string(runes[:maxLen])
	}
	return string(runes[:maxLen-3]) + "..."
}

// CountWords counts the number of words in a string.
func CountWords(s string) int {
	return len(strings.Fields(s))
}

// IsPalindrome checks if a string is a palindrome, ignoring case and non-alphanumeric characters.
func IsPalindrome(s string) bool {
	var cleaned []rune
	for _, r := range strings.ToLower(s) {
		if unicode.IsLetter(r) || unicode.IsDigit(r) {
			cleaned = append(cleaned, r)
		}
	}
	for i, j := 0, len(cleaned)-1; i < j; i, j = i+1, j-1 {
		if cleaned[i] != cleaned[j] {
			return false
		}
	}
	return true
}
