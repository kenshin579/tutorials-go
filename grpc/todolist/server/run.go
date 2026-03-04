package server

import (
	"log"
	"net"

	todopb "github.com/kenshin579/tutorials-go/grpc/todolist/gen/todo/v1"
	"github.com/kenshin579/tutorials-go/grpc/todolist/interceptor"
	"google.golang.org/grpc"
)

func Run() {
	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	s := grpc.NewServer(
		grpc.ChainUnaryInterceptor(
			interceptor.UnaryLogging(),
		),
	)
	todopb.RegisterTodoServiceServer(s, NewTodoService())

	log.Println("gRPC server listening on :50051")
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
