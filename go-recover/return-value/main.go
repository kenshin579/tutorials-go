package main

import (
	"errors"
	"fmt"
	"runtime/debug"
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
			debug.PrintStack()
			err = errors.New(fmt.Sprint(r))
			resp = Response{
				Message: "Failure",
			}
		}
	}()
	panic("test")
	return Response{Message: "Success"}, nil
}
