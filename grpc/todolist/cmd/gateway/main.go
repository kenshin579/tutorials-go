package main

import (
	"context"
	"log"
	"net/http"

	todopb "github.com/kenshin579/tutorials-go/grpc/todolist/gen/todo/v1"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	ctx := context.Background()

	mux := runtime.NewServeMux()
	opts := []grpc.DialOption{grpc.WithTransportCredentials(insecure.NewCredentials())}

	if err := todopb.RegisterTodoServiceHandlerFromEndpoint(ctx, mux, "localhost:50051", opts); err != nil {
		log.Fatalf("failed to register gateway: %v", err)
	}

	log.Println("gRPC-Gateway listening on :8080")
	if err := http.ListenAndServe(":8080", mux); err != nil {
		log.Fatalf("failed to serve gateway: %v", err)
	}
}
