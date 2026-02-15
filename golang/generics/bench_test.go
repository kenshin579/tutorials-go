package go_generics

import (
	"testing"

	"golang.org/x/exp/constraints"
)

// ============================================================
// interface 기반 합산
// ============================================================

type IntAdder int

func (a IntAdder) add(other interface{}) interface{} {
	return a + other.(IntAdder)
}

func sumWithInterface(items []interface{}) interface{} {
	var result IntAdder
	for _, item := range items {
		result = result.add(item).(IntAdder)
	}
	return result
}

// ============================================================
// generics 기반 합산
// ============================================================

func sumWithGenerics[T constraints.Integer | constraints.Float](items []T) T {
	var result T
	for _, item := range items {
		result += item
	}
	return result
}

// ============================================================
// 벤치마크: interface vs generics
// ============================================================

const benchSize = 10000

func BenchmarkSumInterface(b *testing.B) {
	items := make([]interface{}, benchSize)
	for i := 0; i < benchSize; i++ {
		items[i] = IntAdder(i)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		sumWithInterface(items)
	}
}

func BenchmarkSumGenerics(b *testing.B) {
	items := make([]int, benchSize)
	for i := 0; i < benchSize; i++ {
		items[i] = i
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		sumWithGenerics(items)
	}
}

// ============================================================
// 벤치마크: interface vs generics (contains/검색)
// ============================================================

func containsInterface(s []interface{}, target interface{}) bool {
	for _, v := range s {
		if v == target {
			return true
		}
	}
	return false
}

func containsGeneric[T comparable](s []T, target T) bool {
	for _, v := range s {
		if v == target {
			return true
		}
	}
	return false
}

func BenchmarkContainsInterface(b *testing.B) {
	items := make([]interface{}, benchSize)
	for i := 0; i < benchSize; i++ {
		items[i] = i
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		containsInterface(items, benchSize-1) // worst case: 마지막 요소
	}
}

func BenchmarkContainsGenerics(b *testing.B) {
	items := make([]int, benchSize)
	for i := 0; i < benchSize; i++ {
		items[i] = i
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		containsGeneric(items, benchSize-1) // worst case: 마지막 요소
	}
}

// ============================================================
// 벤치마크: 메모리 할당 비교
// ============================================================

func BenchmarkAllocInterface(b *testing.B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		items := make([]interface{}, 1000)
		for j := 0; j < 1000; j++ {
			items[j] = j
		}
	}
}

func BenchmarkAllocGenerics(b *testing.B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		items := make([]int, 1000)
		for j := 0; j < 1000; j++ {
			items[j] = j
		}
	}
}
