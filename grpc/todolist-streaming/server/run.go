package server

import (
	"log"
	"net"

	todopb "github.com/kenshin579/tutorials-go/grpc/todolist-streaming/gen/todo/v1"
	"github.com/kenshin579/tutorials-go/grpc/todolist-streaming/interceptor"
	"google.golang.org/grpc"
)

func Run() {
	lis, err := net.Listen("tcp", ":50052")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	s := grpc.NewServer(
		grpc.ChainUnaryInterceptor(
			interceptor.UnaryLogging(),
		),
		grpc.ChainStreamInterceptor(
			interceptor.StreamLogging(),
		),
	)
	todopb.RegisterTodoStreamingServiceServer(s, NewTodoStreamingService())
	todopb.RegisterTodoUnaryServiceServer(s, NewTodoUnaryService())

	log.Println("gRPC streaming server listening on :50052")
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
