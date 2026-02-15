package go_generics

import (
	"fmt"

	"golang.org/x/exp/constraints"
)

// Generic Min 함수 - constraints.Ordered를 활용
func MinOf[T constraints.Ordered](a, b T) T {
	if a < b {
		return a
	}
	return b
}

// Generic Max 함수
func MaxOf[T constraints.Ordered](a, b T) T {
	if a > b {
		return a
	}
	return b
}

func Example_minOf() {
	fmt.Println(MinOf(10, 20))
	fmt.Println(MinOf(3.14, 2.71))
	fmt.Println(MinOf("apple", "banana"))

	//Output:
	//10
	//2.71
	//apple
}

func Example_maxOf() {
	fmt.Println(MaxOf(10, 20))
	fmt.Println(MaxOf(3.14, 2.71))
	fmt.Println(MaxOf("apple", "banana"))

	//Output:
	//20
	//3.14
	//banana
}

// 슬라이스에서 최솟값을 찾는 함수
func MinSlice[T constraints.Ordered](s []T) (T, bool) {
	var zero T
	if len(s) == 0 {
		return zero, false
	}
	result := s[0]
	for _, v := range s[1:] {
		if v < result {
			result = v
		}
	}
	return result, true
}

// 슬라이스에서 최댓값을 찾는 함수
func MaxSlice[T constraints.Ordered](s []T) (T, bool) {
	var zero T
	if len(s) == 0 {
		return zero, false
	}
	result := s[0]
	for _, v := range s[1:] {
		if v > result {
			result = v
		}
	}
	return result, true
}

func Example_minMaxSlice() {
	minVal, _ := MinSlice([]int{5, 3, 8, 1, 9})
	fmt.Println(minVal)

	maxVal, _ := MaxSlice([]int{5, 3, 8, 1, 9})
	fmt.Println(maxVal)

	minStr, _ := MinSlice([]string{"banana", "apple", "cherry"})
	fmt.Println(minStr)

	//Output:
	//1
	//9
	//apple
}
