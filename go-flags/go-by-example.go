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
	flag.StringVar(&svar, "svar", "bar", "a string var")

	flag.Parse() //커맨드 라인 파싱을 실행함

	fmt.Println("word:", *wordPtr)
	fmt.Println("numb:", *numbPtr)
	fmt.Println("fork:", *boolPtr)
	fmt.Println("svar:", svar)
	fmt.Println("tail:", flag.Args())

}
