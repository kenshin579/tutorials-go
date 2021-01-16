package go_type_assertions

import (
	"fmt"
	"reflect"
)

//참고 : https://hoonyland.medium.com/%EB%B2%88%EC%97%AD-interfaces-in-go-d5ebece9a7ea

// person interface
type Person interface {
	getFullName() string
}

// salaried interface
type Salaried interface {
	getSalary() int
}

// Employee struct represents an employee in an origanization
type Employee struct {
	firstName string
	lastName  string
	salary    int
}

// using this method, Employee implements Person interface
func (e Employee) getFullName() string {
	return e.firstName + " " + e.lastName
}

// using this method, Employee implements Salaried interface
func (e Employee) getSalary() int {
	return e.salary
}

//타입 단언 : 다른 인터페이스로 변환할 때도 사용된다
func Example_TypeAssertion_Convert_Employee_Salary() {
	var johnP Person = Employee{"John", "Adams", 2000}

	// show john's salary
	fmt.Printf("full name : %v \n", reflect.ValueOf(johnP).Interface())

	// convert john to Salaried type
	johnS := johnP.(Salaried)

	fmt.Printf("salary : %v \n", johnS.getSalary())

	//Output:
	//full name : {John Adams 2000}
	//salary : 2000
}
