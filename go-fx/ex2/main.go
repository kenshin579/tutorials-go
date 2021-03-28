package main

import (
	"go.uber.org/fx"
)

type Student struct {
	Name string
	Age  int
}

func newStudent() *Student {
	return &Student{}
}

func main() {
	fx.New(
		fx.Provide(newStudent),
		fx.Invoke(doStuff),
	).Run()
}

func doStuff(student *Student) {
	student.Name = "frank"
	student.Age = 15
}
