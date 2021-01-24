package go_testing

import "testing"

//Function Mocking
func TestPrintSize(t *testing.T) {
	var got string
	oldShow := show
	show = func(v ...interface{}) {
		if len(v) != 1 {
			t.Fatalf("expected show to be called with 1 param, got %d", len(v))
		}
		var ok bool
		got, ok = v[0].(string)
		if !ok {
			t.Fatal("expected show to be called with a string")
		}
	}

	for _, tt := range []struct {
		N   int
		Out string
	}{
		{2, "SMALL"},
		{3, "SMALL"},
		{9, "SMALL"},
		{10, "LARGE"},
		{11, "LARGE"},
		{100, "LARGE"},
	} {
		got = ""
		printSize(tt.N)
		if got != tt.Out {
			t.Fatalf("on %d, expected '%s', got '%s'\n", tt.N, tt.Out, got)
		}
	}

	// careful though, we must not forget to restore it to its original value
	// before finishing the test, or it might interfere with other tests in our
	// suite, giving us unexpected and hard to trace behavior.
	show = oldShow
}
