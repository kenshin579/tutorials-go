package go_generics

import (
	"fmt"
	"sort"
)

// MapValues - map의 value를 슬라이스로 추출
func MapValues[K comparable, V any](m map[K]V) []V {
	values := make([]V, 0, len(m))
	for _, v := range m {
		values = append(values, v)
	}
	return values
}

func Example_mapValues() {
	m := map[string]int{"a": 1, "b": 2, "c": 3}
	values := MapValues(m)
	sort.Ints(values) // map 순서는 비결정적이므로 정렬
	fmt.Println(values)

	//Output:
	//[1 2 3]
}

// MapMerge - 두 map을 합치기 (두 번째 map이 우선)
func MapMerge[K comparable, V any](m1, m2 map[K]V) map[K]V {
	result := make(map[K]V, len(m1)+len(m2))
	for k, v := range m1 {
		result[k] = v
	}
	for k, v := range m2 {
		result[k] = v
	}
	return result
}

func Example_mapMerge() {
	m1 := map[string]int{"a": 1, "b": 2}
	m2 := map[string]int{"b": 20, "c": 3}
	merged := MapMerge(m1, m2)

	// 정렬된 순서로 출력
	keys := make([]string, 0, len(merged))
	for k := range merged {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	for i, k := range keys {
		if i > 0 {
			fmt.Print(" ")
		}
		fmt.Printf("%s:%d", k, merged[k])
	}
	fmt.Println()

	// Output:
	// a:1 b:20 c:3
}

// MapFilter - 조건을 만족하는 key-value 쌍만 반환
func MapFilter[K comparable, V any](m map[K]V, predicate func(K, V) bool) map[K]V {
	result := make(map[K]V)
	for k, v := range m {
		if predicate(k, v) {
			result[k] = v
		}
	}
	return result
}

func Example_mapFilter() {
	scores := map[string]int{"alice": 85, "bob": 42, "carol": 91, "dave": 67}

	// 점수 70 이상만 필터링
	passed := MapFilter(scores, func(k string, v int) bool {
		return v >= 70
	})

	keys := make([]string, 0, len(passed))
	for k := range passed {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	for i, k := range keys {
		if i > 0 {
			fmt.Print(" ")
		}
		fmt.Printf("%s:%d", k, passed[k])
	}
	fmt.Println()

	// Output:
	// alice:85 carol:91
}
