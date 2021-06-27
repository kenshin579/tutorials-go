package _interface

import (
	"fmt"
	"reflect"
)

const (
	DogType AnimalType = "dog"
	CatType AnimalType = "cat"
)

type AnimalType string

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

func NewAnimal(animalType AnimalType, name string) Animal {
	if animalType == DogType {
		return NewDog(name)
	} else if animalType == CatType {
		return NewCat(name)
	} else {
		fmt.Errorf("unknown type: %s", animalType)
	}
	return nil
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

//https://stackoverflow.com/questions/20170275/how-to-find-the-type-of-an-object-in-go
func Example_Polymorphism_Type() {

	animal := NewAnimal(DogType, "dogname")
	fmt.Println(animal.makeSound())

	fmt.Printf("type : %T\n", animal)
	fmt.Printf("type : %s\n", reflect.TypeOf(animal))

	//Output:
	//dogname says 멍멍!
	//type : *_interface.Dog
	//type : *_interface.Dog
}
