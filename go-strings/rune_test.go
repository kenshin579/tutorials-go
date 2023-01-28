package go_strings

import (
	"fmt"
	"unicode/utf8"
)

// rune: 유니코드(UTF-8)을 표현하는 타입
func Example_Rune() {
	//get rune from string
	var r1 rune
	str := "홍길동"

	r1 = []rune(str)[0]

	//get rune from character
	var r2 rune
	r2 = '홍'

	// get rune from character (unicode point)
	var r3 rune
	r3 = '\ud64d'

	//get len of run type (한글)
	len1 := utf8.RuneLen('홍')

	//get len of run type (영어)
	len2 := utf8.RuneLen('a')

	fmt.Println(r1, r2, r3, len1, len2)

	//Output:
	//54861 54861 54861 3 1
}

// byte: ASCII 코드로 표현하는 타입
func Example_byte() {
	//get byte from string
	var b1 byte
	str := "apple"
	b1 = []byte(str)[0]

	//get byte from character
	var b2 byte
	b2 = 'a'

	//get byte from character (ascii)
	var b3 byte
	b3 = 97

	//get byte from character (hex)
	var b4 byte
	b4 = 0x61

	//get len of byte type
	bLen := len("a")

	fmt.Println(b1, b2, b3, b4, bLen)

	//Output:
	//97 97 97 97 1
}

// todo: Output 결과가 왜 true가 안되는지 잘 모르겠음
func Example_한글() {
	s1 := "happy 행복"
	fmt.Println(s1)

	r1 := []rune(s1)
	fmt.Println(r1)

	for i := 0; i < len(r1); i++ {
		fmt.Printf("%d, %s\n", r1[i], string(r1[i]))
	}

	//Output:
	//happy 행복
	//[104 97 112 112 121 32 54665 48373]
	//104, h
	//97, a
	//112, p
	//112, p
	//121, y
	//32,
	//54665, 행
	//48373, 복
}
