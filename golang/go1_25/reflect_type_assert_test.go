package go1_25

import (
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_TypeAssert_성공(t *testing.T) {
	v := reflect.ValueOf(42)

	// Go 1.25: 제네릭 타입 단언 (메모리 할당 없음)
	n, ok := reflect.TypeAssert[int](v)
	assert.True(t, ok)
	assert.Equal(t, 42, n)
}

func Test_TypeAssert_실패(t *testing.T) {
	v := reflect.ValueOf(42)

	// 잘못된 타입으로 단언 → ok=false, 제로값 반환
	s, ok := reflect.TypeAssert[string](v)
	assert.False(t, ok)
	assert.Equal(t, "", s)
}

func Test_TypeAssert_인터페이스_비교(t *testing.T) {
	v := reflect.ValueOf("hello")

	// 기존 방식: Interface()를 통한 타입 단언
	val1, ok1 := v.Interface().(string)

	// Go 1.25 방식: TypeAssert 제네릭 함수
	val2, ok2 := reflect.TypeAssert[string](v)

	// 두 방식 모두 동일한 결과
	assert.Equal(t, ok1, ok2)
	assert.Equal(t, val1, val2)
	assert.Equal(t, "hello", val2)
}

func Test_TypeAssert_구조체(t *testing.T) {
	type Person struct {
		Name string
		Age  int
	}

	p := Person{Name: "Alice", Age: 30}
	v := reflect.ValueOf(p)

	result, ok := reflect.TypeAssert[Person](v)
	assert.True(t, ok)
	assert.Equal(t, "Alice", result.Name)
	assert.Equal(t, 30, result.Age)
}
