package go_generics

import (
	"fmt"
)

// comparable은 == 또는 != 연산이 가능한 타입을 의미하는 내장 constraint이다.
// int, string, bool, struct(필드가 모두 comparable), 포인터, 배열 등이 해당된다.
// slice, map, func은 comparable이 아니다.
func contains[T comparable](s []T, target T) bool {
	for _, v := range s {
		if v == target {
			return true
		}
	}
	return false
}

func Example_comparable() {
	// int 슬라이스에서 검색
	fmt.Println(contains([]int{1, 2, 3, 4, 5}, 3))
	fmt.Println(contains([]int{1, 2, 3, 4, 5}, 6))

	// string 슬라이스에서 검색
	fmt.Println(contains([]string{"go", "java", "python"}, "go"))
	fmt.Println(contains([]string{"go", "java", "python"}, "rust"))

	//Output:
	//true
	//false
	//true
	//false
}

// comparable을 활용한 generic map key 검색
func mapKeys[K comparable, V any](m map[K]V) []K {
	keys := make([]K, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	return keys
}

func Example_comparableMapKeys() {
	m := map[string]int{
		"a": 1,
		"b": 2,
	}
	keys := mapKeys(m)
	// map 순서는 보장되지 않으므로 길이만 확인
	fmt.Println(len(keys))

	//Output:
	//2
}

// comparable constraint를 사용한 중복 제거
func unique[T comparable](s []T) []T {
	seen := make(map[T]bool)
	result := make([]T, 0)
	for _, v := range s {
		if !seen[v] {
			seen[v] = true
			result = append(result, v)
		}
	}
	return result
}

func Example_comparableUnique() {
	ints := unique([]int{1, 2, 2, 3, 3, 3})
	fmt.Println(ints)

	strs := unique([]string{"a", "b", "a", "c", "b"})
	fmt.Println(strs)

	//Output:
	//[1 2 3]
	//[a b c]
}
