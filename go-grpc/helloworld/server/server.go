package main

import (
	"fmt"
	"log"
	"net"

	_ "github.com/jnewmano/grpc-json-proxy/codec"
	"github.com/kenshin579/tutorials-go/go-grpc/helloworld/chat"
	"google.golang.org/grpc"
)

func main() {

	lis, err := net.Listen("tcp", ":9000")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	s := grpc.NewServer()
	chat.RegisterChatServiceServer(s, &server{})

	fmt.Println("Starting gGRPC server")
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %s", err)
	}
}
