package detect_capital

import (
	"fmt"
	"unicode"
)

/*
https://leetcode.com/problems/detect-capital/
https://stackoverflow.com/questions/59293525/how-to-check-if-a-string-is-all-upper-or-lower-case-in-go
*/
func DetectCapitalUse(word string) bool {
	fmt.Println("word", word)
	count := 0
	for _, r := range word {
		fmt.Println("r", r)
		if unicode.IsUpper(r) {
			count++
		}
	}

	return len(word) == count || count == 0 ||
		count == 1 && unicode.IsUpper(rune(word[0]))
}
