package go_getter_setter

import (
	"fmt"
	"log"
)

func ExamplePerson() {
	person := Person{}
	err := person.SetName("Frank")
	if err != nil {
		log.Fatal(err)
	}

	err = person.SetAge(5)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(person.Name())
	fmt.Println(person.Age())

	//Output:
	//Frank
	//5
}
