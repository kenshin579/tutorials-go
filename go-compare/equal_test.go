package go_compare

import (
	"testing"

	"golang.org/x/exp/slices"
)

var equalIntTests = []struct {
	s1, s2 []int
	want   bool
}{
	{
		[]int{1, 2, 3},
		[]int{1, 2, 3},
		true,
	},
	{
		[]int{1, 2, 3},
		[]int{1, 2, 3, 4},
		false,
	},
}

func TestSlice_Equal(t *testing.T) {
	for _, test := range equalIntTests {
		if result := slices.Equal(test.s1, test.s2); result != test.want {
			t.Errorf("Equal(%v, %v) = %t, want %t", test.s1, test.s2, result, test.want)
		}
	}
}

// compare 함수를 넘겨줄 수 있음
func TestSlices_EqualFunc(t *testing.T) {
	for _, test := range equalIntTests {
		if result := slices.EqualFunc(test.s1, test.s2, func(i, j int) bool {
			return i == j
		}); result != test.want {
			t.Errorf("Equal(%v, %v) = %t, want %t", test.s1, test.s2, result, test.want)
		}
	}
}
