package model

type employee struct {
	Name      string
	Role      string
	MinSalary int
	MaxSalary int
}

type employeeBuilder struct {
	e employee
}

func EmployeeBuilder() *employeeBuilder {
	return &employeeBuilder{}
}

func (eb *employeeBuilder) Name(name string) *employeeBuilder {
	eb.e.Name = name
	return eb
}

func (eb *employeeBuilder) Role(role string) *employeeBuilder {
	if role == "manager" {
		eb.e.MinSalary = 20000
		eb.e.MaxSalary = 60000
	}
	eb.e.Role = role
	return eb
}

func (eb *employeeBuilder) Build() employee {
	return employee{
		Name:      eb.e.Name,
		Role:      eb.e.Role,
		MinSalary: eb.e.MinSalary,
		MaxSalary: eb.e.MaxSalary,
	}
}
