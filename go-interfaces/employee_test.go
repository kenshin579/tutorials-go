package go_interfaces

import "fmt"

type Employee interface {
	GetSalary() int
	//Age() int
	//Random() (string, error)
}

type Engineer struct {
	Name string
}

func (e Engineer) GetSalary() int {
	return 1000
}

func ExampleEmployee() {
	// This will throw an error
	var programmers []Employee
	elliot := Engineer{Name: "Elliot"}
	// Engineer does not implement the Employee interface
	// you'll need to implement Age() and Random()
	programmers = append(programmers, elliot)
	fmt.Println(programmers)

	//Output:
	//[{Elliot}]
}
