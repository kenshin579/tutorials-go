package go_generics

import (
	"fmt"
	"strings"
)

// Filter - 조건을 만족하는 요소만 반환
func Filter[T any](s []T, predicate func(T) bool) []T {
	result := make([]T, 0)
	for _, v := range s {
		if predicate(v) {
			result = append(result, v)
		}
	}
	return result
}

func Example_filter() {
	// 짝수만 필터링
	evens := Filter([]int{1, 2, 3, 4, 5, 6}, func(n int) bool {
		return n%2 == 0
	})
	fmt.Println(evens)

	// 길이 3 이상인 문자열만 필터링
	long := Filter([]string{"go", "java", "py", "rust"}, func(s string) bool {
		return len(s) >= 3
	})
	fmt.Println(long)

	//Output:
	//[2 4 6]
	//[java rust]
}

// MapSlice - 슬라이스의 각 요소를 변환 (기존 Map 함수와 이름 충돌 방지)
func MapSlice[T, U any](s []T, transform func(T) U) []U {
	result := make([]U, len(s))
	for i, v := range s {
		result[i] = transform(v)
	}
	return result
}

func Example_mapSlice() {
	// int -> int (2배)
	doubled := MapSlice([]int{1, 2, 3}, func(n int) int {
		return n * 2
	})
	fmt.Println(doubled)

	// string -> string (대문자 변환)
	uppered := MapSlice([]string{"hello", "world"}, func(s string) string {
		return strings.ToUpper(s)
	})
	fmt.Println(uppered)

	// int -> string (타입 변환)
	strs := MapSlice([]int{1, 2, 3}, func(n int) string {
		return fmt.Sprintf("#%d", n)
	})
	fmt.Println(strs)

	//Output:
	//[2 4 6]
	//[HELLO WORLD]
	//[#1 #2 #3]
}

// Reduce - 슬라이스를 하나의 값으로 축약
func Reduce[T, U any](s []T, initial U, accumulator func(U, T) U) U {
	result := initial
	for _, v := range s {
		result = accumulator(result, v)
	}
	return result
}

func Example_reduce() {
	// 합계
	total := Reduce([]int{1, 2, 3, 4, 5}, 0, func(acc, n int) int {
		return acc + n
	})
	fmt.Println(total)

	// 문자열 합치기
	joined := Reduce([]string{"Go", "Generics", "Rock"}, "", func(acc, s string) string {
		if acc == "" {
			return s
		}
		return acc + " " + s
	})
	fmt.Println(joined)

	// 최댓값 구하기
	maxVal := Reduce([]int{3, 7, 2, 9, 4}, 0, func(acc, n int) int {
		if n > acc {
			return n
		}
		return acc
	})
	fmt.Println(maxVal)

	//Output:
	//15
	//Go Generics Rock
	//9
}

// Filter + Map + Reduce 조합 예제
func Example_filterMapReduce() {
	nums := []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}

	// 짝수만 필터링 → 2배로 변환 → 합산
	evens := Filter(nums, func(n int) bool { return n%2 == 0 })
	doubled := MapSlice(evens, func(n int) int { return n * 2 })
	total := Reduce(doubled, 0, func(acc, n int) int { return acc + n })

	fmt.Println(total) // (2+4+6+8+10) * 2 = 60

	//Output:
	//60
}
