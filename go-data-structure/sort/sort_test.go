package sort

import (
	"fmt"
	"sort"
	"strings"

	"github.com/kenshin579/tutorials-go/go-data-structure/sort/model"
	"golang.org/x/exp/slices"
)

func Example_sortInts_Primitive_Type() {
	list := []int{4, 2, 3, 1}
	sort.Ints(list)
	fmt.Println(list)

	//Output:
	// [1 2 3 4]
}

func Example_sortStrings_Primitive_Type() {
	list := []string{"d", "f", "a", "y", "e"}
	sort.Strings(list)
	fmt.Println(list)

	//Output:
	// [a d e f y]
}

func Example_sortFloat64s_Primitive_Type() {
	list := []float64{3, 2, 8, 5, 4}
	sort.Float64s(list)
	fmt.Println(list)

	//Output:
	// [2 3 4 5 8]
}

func Example_sortSliceStable_Struct_With_Custom_Comparator() {
	employee := []struct {
		Name string
		Age  int
	}{
		{"Alice", 23},
		{"David", 2},
		{"Eve", 2},
		{"Bob", 25},
	}

	// Sort by age, keeping original order or equal elements.
	sort.SliceStable(employee, func(i, j int) bool {
		return employee[i].Age < employee[j].Age
	})
	fmt.Println(employee)

	//Output:
	//[{David 2} {Eve 2} {Alice 23} {Bob 25}]
}

/*
SortFunc is generic 함수로 func 에서 index로 받지 않고 실제 객체를 인자로 받음
- SortFunc은 1.18부터 추가됨
*/
func Example_slicesSortFunc() {
	employees := []model.Employee{
		{"Alice", 23},
		{"Eve", 2},
		{"Bob", 25},
	}

	slices.SortFunc(employees, func(x, y model.Employee) bool {
		return x.Age < y.Age
	})
	fmt.Println(employees)

	//Output:
	//[{Eve 2} {Alice 23} {Bob 25}]

}

func Example_sortSort_Any_Collection_By_Implementing_Sort_Interface() {
	employees := []model.Employee{
		{"Alice", 23},
		{"Eve", 2},
		{"Bob", 25},
	}
	sort.Sort(model.ByAge(employees))
	fmt.Println(employees)
	//Output:
	//[{Eve 2} {Alice 23} {Bob 25}]
}

func Example_sortInts_ReverseIntSlice() {
	numbers := []int{4, 3, 2, 1, 0, 4, 7, 5}

	// sort ints ascending
	sort.Ints(numbers)
	fmt.Println(numbers)

	// sort ints descending
	sort.Sort(sort.Reverse(sort.IntSlice(numbers)))
	fmt.Println(numbers)

	//Output:
	//[0 1 2 3 4 4 5 7]
	//[7 5 4 4 3 2 1 0]
}

func Example_SliceStable_Reverse_Sort_Collection_Based_On_Index() {
	employees := []model.Employee{
		{"Alice", 23},
		{"Eve", 2},
		{"Bob", 25},
	}

	sort.SliceStable(employees, func(i, j int) bool {
		return i > j
	})
	fmt.Println(employees)

	//Output:
	//[{Bob 25} {Eve 2} {Alice 23}]
}

func Example_sortStrings_Sort_Map_By_Key_or_Value() {
	m := map[string]int{
		"Alice": 2,
		"Cecil": 1,
		"Bob":   3,
	}

	sortKeys := make([]string, 0, len(m))

	for k := range m {
		fmt.Println("k", k)
		sortKeys = append(sortKeys, k)
	}
	sort.Strings(sortKeys) //keys로 정렬을 함

	//정렬한 keys 값으로 데이터를 출력함
	for _, k := range sortKeys {
		fmt.Println(k, m[k])
	}

	//Output:
	//Alice 2
	//Bob 3
	//Cecil 1
}

func Example_sortSliceStable_Map_By_Key_or_Value2() {
	m := map[string]int{
		"1100::1111": 2,
		"1000::2222": 1,
		"0900::2222": 3,
		"2200::2222": 3,
	}

	sortKeys := make([]string, 0, len(m))

	for k := range m {
		//fmt.Println("k", k)
		sortKeys = append(sortKeys, k)
	}
	sort.SliceStable(sortKeys, func(i, j int) bool {
		return strings.Split(sortKeys[i], "::")[0] < strings.Split(sortKeys[j], "::")[0]
	})

	//fmt.Println(sortKeys)

	//정렬한 keys 값으로 데이터를 출력함
	for _, k := range sortKeys {
		fmt.Println(k, m[k])
	}

	//Output:
	//0900::2222 3
	//1000::2222 1
	//1100::1111 2
	//2200::2222 3
}
