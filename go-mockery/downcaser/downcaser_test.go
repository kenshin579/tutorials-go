package downcaser

import "testing"

//go:generate mockery --inpackage --testonly --name=downcaser
//https://blog.gopheracademy.com/advent-2015/reducing-boilerplate-with-go-generate/
func TestMock(t *testing.T) {
	m := &mockDowncaser{}
	m.On("Downcase", "FOO").Return("foo", nil)
	m.Downcase("FOO")
	m.AssertNumberOfCalls(t, "Downcase", 1)
}
