package go_generics

import (
	"cmp"
	"fmt"
	"maps"
	"slices"
)

// slices 패키지 활용 예제 (Go 1.21+)
func Example_slicesSort() {
	nums := []int{5, 3, 8, 1, 9, 2}
	slices.Sort(nums)
	fmt.Println(nums)

	strs := []string{"banana", "apple", "cherry"}
	slices.Sort(strs)
	fmt.Println(strs)

	//Output:
	//[1 2 3 5 8 9]
	//[apple banana cherry]
}

func Example_slicesContains() {
	nums := []int{1, 2, 3, 4, 5}
	fmt.Println(slices.Contains(nums, 3))
	fmt.Println(slices.Contains(nums, 6))

	//Output:
	//true
	//false
}

func Example_slicesMinMax() {
	nums := []int{5, 3, 8, 1, 9}
	fmt.Println(slices.Min(nums))
	fmt.Println(slices.Max(nums))

	//Output:
	//1
	//9
}

func Example_slicesSortFunc() {
	// 커스텀 정렬: 절대값 기준
	nums := []int{-5, 3, -1, 8, -2}
	slices.SortFunc(nums, func(a, b int) int {
		absA, absB := a, b
		if absA < 0 {
			absA = -absA
		}
		if absB < 0 {
			absB = -absB
		}
		return absA - absB
	})
	fmt.Println(nums)

	//Output:
	//[-1 -2 3 -5 8]
}

// maps 패키지 활용 예제 (Go 1.21+)
func Example_mapsClone() {
	original := map[string]int{"a": 1, "b": 2, "c": 3}
	cloned := maps.Clone(original)
	cloned["d"] = 4 // clone 수정해도 원본 영향 없음

	fmt.Println(len(original))
	fmt.Println(len(cloned))

	//Output:
	//3
	//4
}

func Example_mapsEqual() {
	m1 := map[string]int{"a": 1, "b": 2}
	m2 := map[string]int{"a": 1, "b": 2}
	m3 := map[string]int{"a": 1, "b": 3}

	fmt.Println(maps.Equal(m1, m2))
	fmt.Println(maps.Equal(m1, m3))

	//Output:
	//true
	//false
}

// cmp 패키지 활용 예제 (Go 1.21+)
func Example_cmpCompare() {
	fmt.Println(cmp.Compare(1, 2))  // -1 (1 < 2)
	fmt.Println(cmp.Compare(2, 2))  // 0  (2 == 2)
	fmt.Println(cmp.Compare(3, 2))  // 1  (3 > 2)

	//Output:
	//-1
	//0
	//1
}

func Example_cmpOr() {
	// cmp.Or는 첫 번째 non-zero 값을 반환한다
	// 설정값 우선순위 패턴에 유용: 사용자 설정 > 환경 변수 > 기본값
	userConfig := ""
	envConfig := ""
	defaultConfig := "localhost:8080"

	addr := cmp.Or(userConfig, envConfig, defaultConfig)
	fmt.Println(addr)

	// 환경 변수가 있는 경우
	envConfig2 := "0.0.0.0:9090"
	addr2 := cmp.Or(userConfig, envConfig2, defaultConfig)
	fmt.Println(addr2)

	// Output:
	// localhost:8080
	// 0.0.0.0:9090
}
