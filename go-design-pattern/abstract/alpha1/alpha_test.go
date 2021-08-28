package alpha1

import (
	"fmt"
	"testing"
)

//Abstract Interface
type iAlpha interface {
	work()
	common()
}

//Abstract Concrete Type
type alpha struct {
	name string
}

func (a *alpha) common() {
	fmt.Println("common called")
}

//Implementing Type
type beta struct {
	alpha
}

func (b *beta) work() {
	fmt.Println("work called")
	fmt.Printf("name is %s\n", b.name) //추상 클래스에 있는 필드임
	b.common()                         //추상 메서드임
}

func Test(t *testing.T) {
	a := alpha{
		name: "test",
	}
	b := &beta{
		alpha: a,
	}
	b.work()
}
