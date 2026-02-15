package go_generics

import (
	"fmt"
)

// identity는 입력값을 그대로 반환하는 generic 함수이다.
func identity[T any](v T) T {
	return v
}

// 명시적 타입 지정 vs 타입 추론
func Example_typeInference() {
	// 명시적 타입 지정 - 타입 파라미터를 직접 명시
	result1 := identity[int](42)
	fmt.Println(result1)

	result2 := identity[string]("hello")
	fmt.Println(result2)

	// 타입 추론 - 컴파일러가 인자로부터 타입을 자동 추론
	result3 := identity(42) // T = int 추론
	fmt.Println(result3)

	result4 := identity("hello") // T = string 추론
	fmt.Println(result4)

	//Output:
	//42
	//hello
	//42
	//hello
}

// pair는 두 개의 서로 다른 타입을 받는 generic 함수이다.
func pair[T, U any](a T, b U) string {
	return fmt.Sprintf("(%v, %v)", a, b)
}

// 여러 타입 파라미터에서의 타입 추론
func Example_typeInferenceMultipleParams() {
	// 명시적 타입 지정
	result1 := pair[int, string](1, "hello")
	fmt.Println(result1)

	// 타입 추론 - 두 인자로부터 T=int, U=string 추론
	result2 := pair(1, "hello")
	fmt.Println(result2)

	// 다양한 타입 조합도 추론 가능
	result3 := pair(3.14, true)
	fmt.Println(result3)

	//Output:
	//(1, hello)
	//(1, hello)
	//(3.14, true)
}

// toSlice는 가변 인자를 슬라이스로 변환하는 generic 함수이다.
func toSlice[T any](args ...T) []T {
	return args
}

// 타입 추론이 잘 동작하는 케이스
func Example_typeInferenceSuccess() {
	// 동일 타입 인자 - 추론 성공
	ints := toSlice(1, 2, 3)
	fmt.Println(ints)

	strings := toSlice("a", "b", "c")
	fmt.Println(strings)

	//Output:
	//[1 2 3]
	//[a b c]
}

// 타입 추론이 실패하여 명시적 지정이 필요한 케이스
func Example_typeInferenceExplicitRequired() {
	// 빈 인자 - 컴파일러가 타입을 추론할 수 없으므로 명시적 지정 필요
	emptyInts := toSlice[int]()
	fmt.Println(emptyInts)

	emptyStrings := toSlice[string]()
	fmt.Println(emptyStrings)

	// 리턴 타입만으로는 추론 불가 - 명시적 지정 필요
	// 아래 코드는 컴파일 에러: cannot infer T
	// result := toSlice()

	//Output:
	//[]
	//[]
}
