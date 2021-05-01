package go_functional

import (
	"fmt"
	"strings"

	"github.com/kenshin579/tutorials-go/go-functional/util"
)

func setup() []string {
	return []string{"peach", "apple", "pear", "plum"}
}

func Example_Functional_Index() {
	strList := setup()

	fmt.Println(util.Index(strList, "pear"))
	//Output: 2
}

func Example_Functional_Include() {
	strList := setup()

	fmt.Println(util.Include(strList, "grape"))
	//Output: false
}

func Example_Functional_Any() {
	strList := setup()

	fmt.Println(util.Any(strList, func(v string) bool {
		return strings.HasPrefix(v, "p")
	}))
	//Output: true
}

func Example_Functional_All() {
	strList := setup()

	fmt.Println(util.All(strList, func(v string) bool {
		return strings.HasPrefix(v, "p")
	}))
	//Output: false
}

func Example_Functional_Filter() {
	strList := setup()

	fmt.Println(util.Filter(strList, func(v string) bool {
		return strings.Contains(v, "e")
	}))
	//Output: [peach apple pear]
}

func Example_Functional_Map() {
	strList := setup()

	fmt.Println(util.Map(strList, strings.ToUpper))
	//Output: [PEACH APPLE PEAR PLUM]
}
