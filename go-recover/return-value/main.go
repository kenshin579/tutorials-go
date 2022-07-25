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
			err = convertRecoverToError(r)
			resp = Response{
				Message: "failure",
			}
		}
	}()
	panic("test")
	return Response{Message: "success"}, nil
}

func convertRecoverToError(r interface{}) error {
	switch x := r.(type) {
	case string:
		return errors.New(x)
	case error:
		return x
	default:
		return errors.New(fmt.Sprint(x))
	}
}
