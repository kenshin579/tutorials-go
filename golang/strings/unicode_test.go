package go_strings

import (
	"testing"
	"unicode"

	"github.com/stretchr/testify/assert"
)

func Test_Unicode_Is(t *testing.T) {
	var r1 = '한'
	assert.True(t, unicode.Is(unicode.Hangul, r1))
	assert.False(t, unicode.Is(unicode.Latin, r1))

	var r2 = '漢'
	assert.True(t, unicode.Is(unicode.Han, r2))
	assert.False(t, unicode.Is(unicode.Hangul, r2))

	var r3 = 'a'
	assert.True(t, unicode.Is(unicode.Latin, r3))
	assert.False(t, unicode.Is(unicode.Hangul, r3))

}

// In(): 문자가 여러 범위 테이블 중에 포함되는지 확인
func Test_Unicode_In(t *testing.T) {
	var r1 = '한'

	assert.True(t, unicode.In(r1, unicode.Latin, unicode.Han, unicode.Hangul))
}

func Test_Unicode_Funcs(t *testing.T) {
	assert.True(t, unicode.IsGraphic('1'))   // true: 1은 화면에 표시되는 숫자이므로 true
	assert.True(t, unicode.IsGraphic('a'))   // true: a는 화면에 표시되는 문자이므로 true
	assert.True(t, unicode.IsGraphic('한'))   // true: '한'은 화면에 표시되는 문자이므로 true
	assert.True(t, unicode.IsGraphic('漢'))   // true: '漢'은 화면에 표시되는 문자이므로 true
	assert.False(t, unicode.IsGraphic('\n')) // false: \n 화면에 표시되는 문자가 아니므로 false

	assert.True(t, unicode.IsLetter('a'))  // true: a는 문자이므로 true
	assert.False(t, unicode.IsLetter('1')) // false: 1은 문자가 아니므로 false

	assert.True(t, unicode.IsDigit('1'))     // true: 1은 숫자이므로 true
	assert.True(t, unicode.IsControl('\n'))  // true: \n은 제어 문자이므로 true
	assert.True(t, unicode.IsMark('\u17c9')) // true: \u17c9는 마크이므로 true

	assert.True(t, unicode.IsPrint('1')) // true: 1은 Go 언어에서 표시할 수 있으므로 true
	assert.True(t, unicode.IsPunct('.')) // true: .은 문장 부호이므로 true

	assert.True(t, unicode.IsSpace(' '))  // true: ' '는 공백이므로 true
	assert.True(t, unicode.IsSymbol('♥')) // true: ♥는 심볼이므로 true

	assert.True(t, unicode.IsUpper('A')) // true: A는 대문자이므로 true
	assert.True(t, unicode.IsLower('a')) // true: a는 소문자이므로 true
}
