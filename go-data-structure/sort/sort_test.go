package sort

import (
	"fmt"
	"sort"

	"github.com/kenshin579/tutorials-go/go-data-structure/sort/model"
)

func Example_Sort_Int_Primitive_Type() {
	list := []int{4, 2, 3, 1}
	sort.Ints(list)
	fmt.Println(list)

	//Output:
	// [1 2 3 4]
}

func Example_Sort_String_Primitive_Type() {
	list := []string{"d", "f", "a", "y", "e"}
	sort.Strings(list)
	fmt.Println(list)

	//Output:
	// [a d e f y]
}

func Example_Sort_Floats_Primitive_Type() {
	list := []float64{3, 2, 8, 5, 4}
	sort.Float64s(list)
	fmt.Println(list)

	//Output:
	// [2 3 4 5 8]
}

func Example_Sort_Struct_With_Custom_Comparator() {
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

func Example_Sort_Any_Collection_By_Implementing_Sort_Interface() {
	family := []model.Employee{
		{"Alice", 23},
		{"Eve", 2},
		{"Bob", 25},
	}
	sort.Sort(model.ByAge(family))
	fmt.Println(family)
	//Output:
	//[{Eve 2} {Alice 23} {Bob 25}]
}

func Example_Sort_Map_By_Key_or_Value() {
	m := map[string]int{
		"Alice": 2,
		"Cecil": 1,
		"Bob":   3,
	}

	sortKeys := make([]string, 0, len(m))

	for k := range m {
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
