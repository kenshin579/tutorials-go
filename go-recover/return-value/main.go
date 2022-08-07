package main

import (
	"errors"
	"fmt"
)

func main() {
	myFunc, err := MyFunc()
	fmt.Println("myFunc:", myFunc)
	fmt.Println("err:", err)
}

type Response struct {
	Message string
}

func MyFunc() (resp Response, err error) {
	defer func() {
		if r := recover(); r != nil {
			err = errors.New(fmt.Sprint(r))
			resp = Response{
				Message: "failure",
			}
		}
	}()
	panic("test")
	return Response{Message: "success"}, nil
}
