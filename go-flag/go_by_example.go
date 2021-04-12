package main

import (
	"flag"
	"fmt"
)

//https://mingrammer.com/gobyexample/command-line-flags/
func main() {
	//옵션 이름과 기본값을 설정한다
	wordPtr := flag.String("word", "foo", "a string")
	numbPtr := flag.Int("numb", 42, "an int")
	boolPtr := flag.Bool("fork", false, "a bool")

	var svar string
	//StringVar() 저장할 변수를 인자로 받는다
	flag.StringVar(&svar, "svar", "bar", "a string var")

	//커맨드 라인 파싱을 실행함
	flag.Parse()

	fmt.Println("word:", *wordPtr)
	fmt.Println("numb:", *numbPtr)
	fmt.Println("fork:", *boolPtr)
	fmt.Println("svar:", svar)
	fmt.Println("tail:", flag.Args())
}
