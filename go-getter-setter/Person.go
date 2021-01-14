package go_getter_setter

import "errors"

type Person struct {
	name string
	age  int
}

func (p *Person) Name() string {
	return p.name
}

func (p *Person) SetName(name string) error {
	if name == "" {
		return errors.New("invalid name")
	}
	p.name = name
	return nil
}

func (p *Person) Age() int {
	return p.age
}

func (p *Person) SetAge(age int) error {
	if age == 0 {
		return errors.New("invalid age")
	}
	p.age = age
	return nil
}
