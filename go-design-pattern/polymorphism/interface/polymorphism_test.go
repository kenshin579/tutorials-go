package _interface

import "fmt"

type Animal interface {
	makeNoise() string
}

type Dog struct {
	name string
	legs int
}

type Duck struct {
	name string
	legs int
}

func (d *Dog) makeNoise() string {
	return d.name + " says woof!"
}
func (d *Duck) makeNoise() string {
	return d.name + " says quack!"
}

func NewDog(name string) Animal {
	return &Dog{
		legs: 4,
		name: name,
	}
}

func NewDuck(name string) Animal {
	return &Duck{
		legs: 4,
		name: name,
	}
}

func Example_Polymorphism_Interface_사용하는_방법() {
	var dog, duck Animal

	dog = NewDog("fido")
	duck = NewDuck("donald")

	fmt.Println(dog.makeNoise())
	fmt.Println(duck.makeNoise())

	//Output:
	//fido says woof!
	//donald says quack!
}
