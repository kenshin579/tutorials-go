package main

import (
	"sync"
	"testing"
)

// GoLand IDE에서 벤치마크 프로파일링을 실행하는 예제
// 직접 접근(slicerInBounds) vs 채널 통신(slicerInBoundsChannels) 방식의 성능을 비교한다.
//
// GoLand에서 벤치마크 좌측의 실행 아이콘 → "Profile"로 프로파일링을 실행할 수 있다.
// CLI에서 실행: go test -bench=. -cpuprofile=cpu.prof -memprofile=mem.prof
//
// 참고: https://blog.jetbrains.com/go/2019/04/03/profiling-go-applications-and-tests/
var pi []int

func printlner(i ...int) {
	pi = i
}

// mySliceType은 뮤텍스로 보호되는 슬라이스 타입이다.
// Get과 GetCh에서 checkBuffer를 호출하므로 의도적으로 CPU/메모리 부하가 발생한다.
type mySliceType struct {
	valuesGuard *sync.Mutex
	values      []int
}

// Get은 뮤텍스를 획득한 후 checkBuffer로 부하를 발생시키고 값을 반환한다.
func (s mySliceType) Get(idx int) int {
	s.valuesGuard.Lock()
	defer s.valuesGuard.Unlock()
	checkBuffer(s.values, idx)
	return s.values[idx]
}

// GetCh는 Get과 동일하지만 결과를 채널로 전달한다. 고루틴에서 호출하여 채널 통신 오버헤드를 측정한다.
func (s mySliceType) GetCh(ch chan int, idx int) {
	s.valuesGuard.Lock()
	defer s.valuesGuard.Unlock()
	checkBuffer(s.values, idx)
	ch <- s.values[idx]
}

func newMySliceType(values []int) mySliceType {
	return mySliceType{
		valuesGuard: &sync.Mutex{},
		values:      values,
	}
}

// fillBuffer는 의도적으로 큰 map을 생성하여 CPU와 메모리 부하를 발생시킨다.
// pprof에서 이 함수의 할당량이 높게 나타난다.
func fillBuffer(slice []int) map[int]int {
	result := map[int]int{}
	for i := 0; i < 100; i++ {
		for j := 0; j < len(slice); j++ {
			result[i*len(slice)+j] = slice[j]
		}
	}
	return result
}

// checkBuffer는 매 호출마다 map을 새로 할당하므로 GC 압박을 유발한다.
// 프로파일에서 alloc_space가 높게 나타나는 원인이 된다.
func checkBuffer(slice []int, idx int) {
	buffer := make(map[int]int, len(slice)*100)
	buffer = fillBuffer(slice)
	for i := range buffer {
		if i == idx {
			return
		}
	}
}

// slicerInBounds는 직접 호출 방식으로 슬라이스 값을 조회한다.
// 뮤텍스 경합만 발생하고 채널 오버헤드가 없어 상대적으로 빠르다.
func slicerInBounds(slice mySliceType) {
	for i := 0; i < 8; i++ {
		a0 := slice.Get(i*8 + 0)
		a1 := slice.Get(i*8 + 1)
		a2 := slice.Get(i*8 + 2)
		a3 := slice.Get(i*8 + 3)
		a4 := slice.Get(i*8 + 4)
		a5 := slice.Get(i*8 + 5)
		a6 := slice.Get(i*8 + 6)
		a7 := slice.Get(i*8 + 7)
		printlner(a0, a1, a2, a3, a4, a5, a6, a7)
	}
}

// slicerInBoundsChannels는 채널 통신 방식으로 슬라이스 값을 조회한다.
// 8개의 고루틴이 동시에 실행되지만, 고루틴 생성/채널 통신 오버헤드로 인해
// 직접 호출 방식보다 느릴 수 있다.
func slicerInBoundsChannels(slice mySliceType) {
	ch := make(chan int, 8)
	for i := 0; i < 8; i++ {
		go slice.GetCh(ch, i*8+0)
		go slice.GetCh(ch, i*8+1)
		go slice.GetCh(ch, i*8+2)
		go slice.GetCh(ch, i*8+3)
		go slice.GetCh(ch, i*8+4)
		go slice.GetCh(ch, i*8+5)
		go slice.GetCh(ch, i*8+6)
		go slice.GetCh(ch, i*8+7)
		a0 := <-ch
		a1 := <-ch
		a2 := <-ch
		a3 := <-ch
		a4 := <-ch
		a5 := <-ch
		a6 := <-ch
		a7 := <-ch
		printlner(a0, a1, a2, a3, a4, a5, a6, a7)
	}
}

// BenchmarkInBounds는 직접 호출 방식의 성능을 측정한다.
func BenchmarkInBounds(b *testing.B) {
	var mySlice []int
	for i := 0; i < 99; i++ {
		mySlice = append(mySlice, i)
	}
	ms := newMySliceType(mySlice)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		slicerInBounds(ms)
	}
}

// BenchmarkInBoundsChannels는 채널 통신 방식의 성능을 측정한다.
// BenchmarkInBounds와 비교하여 채널/고루틴 오버헤드를 확인할 수 있다.
func BenchmarkInBoundsChannels(b *testing.B) {
	var mySlice []int
	for i := 0; i < 99; i++ {
		mySlice = append(mySlice, i)
	}
	ms := newMySliceType(mySlice)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		slicerInBoundsChannels(ms)
	}
}
