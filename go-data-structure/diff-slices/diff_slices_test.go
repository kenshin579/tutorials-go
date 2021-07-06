package diff_slices

import (
	"fmt"
)

//https://stackoverflow.com/questions/19374219/how-to-find-the-difference-between-two-slices-of-strings
// difference returns the elements in `a` that aren't in `b`.
func difference(a, b []string) []string {
	bMap := make(map[string]struct{}, len(b))
	for _, x := range b {
		bMap[x] = struct{}{}
	}
	var diff []string
	for _, x := range a {
		if _, found := bMap[x]; !found {
			diff = append(diff, x)
		}
	}
	return diff
}

func ExampleDiff_Two_Slices() {
	fmt.Println(difference(
		[]string{"111", "222", "333"},
		[]string{"111", "222"},
	))

	fmt.Println(difference(
		[]string{"111", "222"},
		[]string{"111", "222", "333"},
	))

	fmt.Println(difference(
		[]string{"111", "222"},
		[]string{},
	))

	//Output:
	//[333]
	//[]
	//[111 222]
}
