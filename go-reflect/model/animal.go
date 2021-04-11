package model

import "fmt"

type Animal interface {
	Move(int) bool
}

type Cat struct {
	Name  string   `custom:"name"`
	Age   int      `custom:"age"`
	Child []string `custom:"child"`
}

func (c *Cat) Move(distance int) bool {
	fmt.Println("Cat Move ", distance)
	return true
}

type Dog struct {
	Name  string   `custom:"name"`
	Age   int      `custom:"age"`
	Child []string `custom:"child"`
}

func (d *Dog) Move(distance int) bool {
	fmt.Println("Dog Move ", distance)
	return false
}
