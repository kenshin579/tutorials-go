package go_testing

import "testing"

func TestAvg_Table_driven_test(t *testing.T) {
	for _, tt := range []struct {
		Nos    []int
		Result int
	}{
		{Nos: []int{2, 4}, Result: 3},
		{Nos: []int{1, 2, 5}, Result: 2},
		{Nos: []int{1}, Result: 1},
		{Nos: []int{}, Result: 0},
		{Nos: []int{2, -2}, Result: 0},
	} {
		if avg := Average(tt.Nos...); avg != tt.Result {
			t.Fatalf("expected average of %v to be %d, got %d\n", tt.Nos, tt.Result, avg)
		}
	}
}
