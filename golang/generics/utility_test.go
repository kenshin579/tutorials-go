package go_generics

import (
	"fmt"
	"sort"
	"strings"
)

// ============================================================
// 1. Result 타입 - 에러 처리를 위한 Generic 래퍼
// ============================================================

// Result - 성공 값 또는 에러를 담는 타입
type Result[T any] struct {
	value T
	err   error
}

func OK[T any](value T) Result[T] {
	return Result[T]{value: value}
}

func Fail[T any](err error) Result[T] {
	return Result[T]{err: err}
}

func (r Result[T]) IsOK() bool {
	return r.err == nil
}

func (r Result[T]) Unwrap() (T, error) {
	return r.value, r.err
}

// MapResult - Result 값을 변환
func MapResult[T, U any](r Result[T], f func(T) U) Result[U] {
	if r.err != nil {
		return Fail[U](r.err)
	}
	return OK(f(r.value))
}

func Example_resultType() {
	// 성공 케이스
	r1 := OK(42)
	fmt.Println(r1.IsOK())
	val, _ := r1.Unwrap()
	fmt.Println(val)

	// 변환
	r2 := MapResult(r1, func(n int) string {
		return fmt.Sprintf("value=%d", n)
	})
	str, _ := r2.Unwrap()
	fmt.Println(str)

	// 실패 케이스
	r3 := Fail[int](fmt.Errorf("something failed"))
	fmt.Println(r3.IsOK())
	_, err := r3.Unwrap()
	fmt.Println(err)

	// Output:
	// true
	// 42
	// value=42
	// false
	// something failed
}

// ============================================================
// 2. Pair 타입 - 두 값을 묶는 Generic 타입
// ============================================================

// Pair - 두 값을 묶는 제네릭 타입
type Pair[K, V any] struct {
	Key   K
	Value V
}

func NewPair[K, V any](k K, v V) Pair[K, V] {
	return Pair[K, V]{Key: k, Value: v}
}

// Zip - 두 슬라이스를 Pair 슬라이스로 합침
func Zip[K, V any](keys []K, values []V) []Pair[K, V] {
	n := len(keys)
	if len(values) < n {
		n = len(values)
	}
	result := make([]Pair[K, V], n)
	for i := 0; i < n; i++ {
		result[i] = Pair[K, V]{Key: keys[i], Value: values[i]}
	}
	return result
}

func Example_zipPairs() {
	names := []string{"Alice", "Bob", "Charlie"}
	scores := []int{95, 87, 92}

	pairs := Zip(names, scores)
	for _, p := range pairs {
		fmt.Printf("%s:%d\n", p.Key, p.Value)
	}

	// Output:
	// Alice:95
	// Bob:87
	// Charlie:92
}

// ============================================================
// 3. GroupBy - 슬라이스를 그룹핑하는 함수
// ============================================================

// GroupBy - 키 함수로 슬라이스를 그룹핑
func GroupBy[T any, K comparable](items []T, keyFn func(T) K) map[K][]T {
	result := make(map[K][]T)
	for _, item := range items {
		key := keyFn(item)
		result[key] = append(result[key], item)
	}
	return result
}

func Example_groupBy() {
	words := []string{"apple", "banana", "avocado", "blueberry", "cherry"}

	groups := GroupBy(words, func(s string) string {
		return string(s[0]) // 첫 글자로 그룹핑
	})

	// 키를 정렬하여 출력 순서 보장
	keys := make([]string, 0, len(groups))
	for k := range groups {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	for _, k := range keys {
		fmt.Printf("%s: %s\n", k, strings.Join(groups[k], ", "))
	}

	// Output:
	// a: apple, avocado
	// b: banana, blueberry
	// c: cherry
}

// ============================================================
// 4. ChunkSlice - 슬라이스를 청크 단위로 분할
// ============================================================

// ChunkSlice - 슬라이스를 지정 크기의 청크로 분할
func ChunkSlice[T any](items []T, size int) [][]T {
	if size <= 0 {
		return nil
	}
	var chunks [][]T
	for i := 0; i < len(items); i += size {
		end := i + size
		if end > len(items) {
			end = len(items)
		}
		chunks = append(chunks, items[i:end])
	}
	return chunks
}

func Example_chunkSlice() {
	nums := []int{1, 2, 3, 4, 5, 6, 7}
	chunks := ChunkSlice(nums, 3)
	for _, chunk := range chunks {
		fmt.Println(chunk)
	}

	// Output:
	// [1 2 3]
	// [4 5 6]
	// [7]
}
