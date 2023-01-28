package go_strings

import (
	"fmt"
	"testing"
	"unicode/utf8"

	"github.com/stretchr/testify/assert"
)

// RuneLen 문자의 바이트 수를 구함
func Example_RuneLen() {
	var s string = "한"
	fmt.Println(len(s)) // 3: 한글은 3바이트로 저장하므로 3

	var r rune = '한'
	fmt.Println(utf8.RuneLen(r)) // 3: 한글은 3바이트로 저장하므로 3
	fmt.Println()

	//Output:
	//3
	//3
}

// RuneCountInString 문자열의 실제 길이를 구함
func Example_RuneCountInString() {
	var s1 string = "안녕하세요"
	fmt.Println(utf8.RuneCountInString(s1)) // 5: "안녕하세요"의 실제 길이는 5
	fmt.Println()

	//Output:
	//5
}

// DecodeRun 함수는 byte 슬라이스에서 첫 글자를 디코딩하여 리턴하고, 디코딩된 byte 수도 반환한다.
// 글자 하나만 디코딩하여 다음 글자를 디코딩하려면 b[3:]과 같이 부분 슬라이스로 만들어주면 된다

// DecodeLastRun 함수는 byte 슬라이스에서 마지막 글자를 디코딩하여 리턴하고, 디코딩된 바이트 수도 반환한다.
// 글자 하나만 디코딩하며 다른 글자를 디코딩하려면 b[:len(b)-3]과 같이 부분 슬라이스로 만들어주면 된다
func Example_DecodeRun() {
	b := []byte("안녕하세요")

	r, size := utf8.DecodeRune(b)
	fmt.Printf("%c %d\n", r, size) // 안 3: "안녕하세요"의 첫 글자를 디코딩하여 '안', 바이트 수 3

	r, size = utf8.DecodeRune(b[3:]) // '안'의 길이가 3이므로 인덱스 3부터 부분 슬라이스를 만들면 "녕하세요"가 됨
	fmt.Printf("%c %d\n", r, size)   // 녕 3: "녕하세요"를 첫 글자를 디코딩하여 '녕', 바이트 수 3

	r, size = utf8.DecodeLastRune(b)
	fmt.Printf("%c %d\n", r, size) // 요 3: "안녕하세요"의 마지막 글자를 디코딩하여 '요', 바이트 수 3

	// '요'의 길이가 3이므로 // 문자열 길이-3을 하여 부분 슬라이스를 만들면
	// "안녕하세"가 됨
	r, size = utf8.DecodeLastRune(b[:len(b)-3])

	fmt.Printf("%c %d\n", r, size) // 세 3: "안녕하세"의 마지막 글자를 디코딩하여 '세', 바이트 수 3
	fmt.Println()

	//Output:
	//안 3
	//녕 3
	//요 3
	//세 3
}

func Test_String_인덱스_접근(t *testing.T) {
	s1 := "Hello, world!"

	fmt.Printf("%c\n", s1[0]) // H: 인덱스 0이 첫 번째 글자
}

// DecodeRuneInString 함수는 UTF-8 문자열에서 첫 글자를 디코딩하여 리턴하고, 디코딩된 바이트 수도 반환한다
func Test_String_한글_인덱스_접근(t *testing.T) {
	s1 := "안녕하세요"

	fmt.Printf("%c\n", s1[0])         // 3바이트 중 1바이트만 출력하므로 한글이 정상적으로 출력되지 않음
	fmt.Printf("%c\n", s1[len(s1)-1]) // 3바이트 중 1바이트만 출력하므로 한글이 정상적으로 출력되지 않음

	r, _ := utf8.DecodeRuneInString(s1) // UTF-8 문자열의 첫 글자와 바이트 수를 리턴
	fmt.Printf("%c\n", r)               // 안: 문자열의 첫 글자

	r, _ = utf8.DecodeLastRuneInString(s1) // UTF-8 문자열의 첫 글자와 바이트 수를 리턴
	fmt.Printf("%c\n", r)                  // 요: 문자열의 마지막 글자
}

func Test_Valid(t *testing.T) {
	var b1 []byte = []byte("안녕하세요")
	assert.True(t, utf8.Valid(b1)) // true: "안녕하세요"는 UTF-8이 맞으므로 true
	var b2 []byte = []byte{0xff, 0xf1, 0xc1}
	assert.False(t, utf8.Valid(b2)) // false: 0xff 0xf1 0xc1은 UTF-8이 아니므로 false

	var r1 rune = '한'
	assert.True(t, utf8.ValidRune(r1)) // true: '한'은 UTF-8이 맞으므로 true
	var r2 rune = 0x11111111
	assert.False(t, utf8.ValidRune(r2)) // false: 0x11111111은 UTF-8이 아니므로 false

	var s1 string = "한글"
	assert.True(t, utf8.ValidString(s1)) // true: "한글"은 UTF-8이 맞으므로 true
	var s2 string = string([]byte{0xff, 0xf1, 0xc1})
	assert.False(t, utf8.ValidString(s2)) // false: 0xff 0xf1 0xc1은 UTF-8이 아니므로 false
}
