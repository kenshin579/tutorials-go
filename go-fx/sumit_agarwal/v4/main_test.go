package main

import (
	"fmt"
	"net/rpc"
	"testing"
)

func TestRpc(t *testing.T) {
	//호출이 안됨
	client, _ := rpc.Dial("tcp", "localhost:8081")
	if err := client.Call("Handler.GetUsers", 1, nil); err != nil {
		fmt.Printf("Error:1 user.GetUsers() %+v", err)
	} else {
		fmt.Printf("user found")
	}
}
