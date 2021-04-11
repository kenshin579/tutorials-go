package model

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_Builder_Employee(t *testing.T) {
	emp1 := EmployeeBuilder().
		Name("Michael Scott").
		Role("manager").
		Build()
	fmt.Println(emp1)
	assert.Equal(t, "Michael Scott", emp1.Name)
	assert.Equal(t, "manager", emp1.Role)

	emp2 := EmployeeBuilder().
		Name("Michael Scott").
		Role("manager").
		Build()
	fmt.Println(emp2)
	assert.Equal(t, "Michael Scott", emp2.Name)
	assert.Equal(t, "manager", emp2.Role)
}
