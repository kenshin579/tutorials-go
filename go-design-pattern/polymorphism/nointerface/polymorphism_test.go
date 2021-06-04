package nointerface

import (
	"fmt"
)

//https://www.sohamkamani.com/golang/2019-03-29-polymorphism-without-interfaces/
type Animal struct {
	makeNoiseFn func(*Animal) string
	name        string
	legs        int
}

//여기서 makeNoise()는 makeNoiseFn의 wrapper 역할을 함
func (a *Animal) makeNoise() string {
	return a.makeNoiseFn(a)
}

func NewDog(name string) *Animal {
	return &Animal{
		makeNoiseFn: func(a *Animal) string {
			return a.name + " says woof!"
		},
		legs: 4,
		name: name,
	}
}

func NewDuck(name string) *Animal {
	return &Animal{
		makeNoiseFn: func(a *Animal) string {
			return a.name + " says quack!"
		},
		legs: 4,
		name: name,
	}
}

func Example_Polymorphism_Interface_사용하지_않는_방법() {
	var dog, duck *Animal

	dog = NewDog("fido")
	duck = NewDuck("donald")

	fmt.Println(dog.makeNoise())
	fmt.Println(duck.makeNoise())

	//Output:
	//fido says woof!
	//donald says quack!
}
