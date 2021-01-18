package go_methods

import "fmt"

type Employee struct {
	name   string
	age    int
	salary int
}

func (e Employee) Name() string {
	return e.name
}

func (e *Employee) Salary() int {
	return e.salary
}

func Example_Method_Value_Receiver() {
	employee := Employee{"Frank", 20, 1000}
	fmt.Println(employee.Name())

	//Output:
	//Frank
}

func Example_Method_Pointer_Receiver() {
	employee := Employee{"Frank", 20, 1000}
	fmt.Println(employee.Name())

	//Output:
	//Frank
}
