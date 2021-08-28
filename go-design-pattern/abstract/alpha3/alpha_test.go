package alpha2

import (
	"fmt"
	"testing"
)

//Abstract Interface
type iAlpha interface {
	work()
	common(iAlpha)
}

//Abstract Concrete Type
type alpha struct {
	name string
}

func (a *alpha) common(i iAlpha) {
	fmt.Println("common called")
	i.work()
}

//Implementing Type
type beta struct {
	alpha
}

func (b *beta) work() {
	fmt.Println("work called")
	fmt.Printf("Name is %s\n", b.name)
}

func Test(t *testing.T) {
	a := alpha{
		name: "test",
	}
	b := &beta{
		alpha: a,
	}
	b.common(b)
}
