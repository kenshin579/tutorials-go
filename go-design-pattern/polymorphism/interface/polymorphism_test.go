package _interface

import "fmt"

type Animal interface {
	makeSound() string
}

type Dog struct {
	name string
	legs int
}

type Cat struct {
	name string
	legs int
}

func (d *Dog) makeSound() string {
	return d.name + " says 멍멍!"
}
func (c *Cat) makeSound() string {
	return c.name + " says 야옹!"
}

func NewDog(name string) Animal {
	return &Dog{
		legs: 4,
		name: name,
	}
}

func NewCat(name string) Animal {
	return &Cat{
		legs: 4,
		name: name,
	}
}

func Example_Polymorphism_Interface_사용하는_방법() {
	var dog, cat Animal

	dog = NewDog("초코")
	cat = NewCat("루시")

	fmt.Println(dog.makeSound())
	fmt.Println(cat.makeSound())

	//Output:
	//초코 says 멍멍!
	//루시 says 야옹!
}
