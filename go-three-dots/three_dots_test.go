package go_three_dots

import (
	"fmt"
	"reflect"
)

func sum(nums ...int) int {
	res := 0
	for _, n := range nums {
		res += n
	}
	return res
}

//1.함수에 가변인자로 정의하는 경우
func Example_가변인자_함수() {
	total1 := sum(1, 2, 3)
	total2 := sum(1, 2)
	fmt.Println(total1)
	fmt.Println(total2)

	//Output:
	//6
	//3
}

//2.가변인자 함수에 전달할 때 unpack해서 넘겨주는 경우
func Example_가변인자_함수에_전달하기() {
	numList := []int{2, 3, 5, 6}
	fmt.Println(sum(numList[0], numList[1], numList[2], numList[3]))
	//...표기법을 통해서 가변인자에 unpack해서 전달할 수 있다
	fmt.Println(sum(numList...))

	//Output:
	//16
	//16
}

//3.배열 리터럴에서 크기지정하는 경우
func Example_array_literal() {
	//배열 리터럴에서 ... 표기법은 리터럴의 요소 수와 동일한 길이를 지정한다
	strList := [...]string{"Frank", "Joe", "Angela"}
	fmt.Println(len(strList))

	typ := reflect.TypeOf(strList)
	fmt.Println("type", typ.Kind())

	//Output:
	//3
	//type array
}

//4.패키지 목록을 지정할 때 ...표기법은 패키지 목록을 wildcard로 사용된다
//$ go test ./...
