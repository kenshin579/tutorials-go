package model

type Employee struct {
	Name      string
	Role      string
	MinSalary int
	MaxSalary int
}

type EmployeeBuilder struct {
	e Employee
}

func (b *EmployeeBuilder) Build() Employee {
	return b.e
}

func (b *EmployeeBuilder) Name(name string) *EmployeeBuilder {
	b.e.Name = name
	return b
}

func (b *EmployeeBuilder) Role(role string) *EmployeeBuilder {
	if role == "manager" {
		b.e.MinSalary = 20000
		b.e.MaxSalary = 60000
	}
	b.e.Role = role
	return b
}
