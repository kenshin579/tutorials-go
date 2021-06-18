package nointerface

import (
	"fmt"
)

//https://www.sohamkamani.com/golang/2019-03-29-polymorphism-without-interfaces/
type Animal struct {
	makeSound func(*Animal) string
	name      string
	legs      int
}

//여기서 makeNoise()는 makeNoiseFn의 wrapper 역할을 함
func (a *Animal) makeNoise() string {
	return a.makeSound(a)
}

func NewDog(name string) *Animal {
	return &Animal{
		makeSound: func(a *Animal) string {
			return a.name + " says 멍멍!"
		},
		legs: 4,
		name: name,
	}
}

func NewCat(name string) *Animal {
	return &Animal{
		makeSound: func(a *Animal) string {
			return a.name + " says 야옹!"
		},
		legs: 4,
		name: name,
	}
}

func Example_Polymorphism_Interface_사용하지_않는_방법() {
	var dog, cat *Animal

	dog = NewDog("초코")
	cat = NewCat("루시")

	fmt.Println(dog.makeNoise())
	fmt.Println(cat.makeNoise())

	//Output:
	//초코 says 멍멍!
	//루시 says 야옹!
}
