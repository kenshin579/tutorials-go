package main

import (
	"context"
	"log"
	"time"

	todopb "github.com/kenshin579/tutorials-go/grpc/todolist/gen/todo/v1"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	conn, err := grpc.NewClient("localhost:50051",
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		log.Fatalf("failed to connect: %v", err)
	}
	defer conn.Close()

	client := todopb.NewTodoServiceClient(conn)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Create
	createResp, err := client.CreateTodo(ctx, &todopb.CreateTodoRequest{Title: "Learn gRPC"})
	if err != nil {
		log.Fatalf("CreateTodo failed: %v", err)
	}
	log.Printf("Created: %v", createResp.GetTodo())

	todoID := createResp.GetTodo().GetId()

	// List
	listResp, err := client.ListTodos(ctx, &todopb.ListTodosRequest{})
	if err != nil {
		log.Fatalf("ListTodos failed: %v", err)
	}
	log.Printf("Todos: %v", listResp.GetTodos())

	// Update
	updateResp, err := client.UpdateTodo(ctx, &todopb.UpdateTodoRequest{
		Id:        todoID,
		Title:     "Learn gRPC (done)",
		Completed: true,
	})
	if err != nil {
		log.Fatalf("UpdateTodo failed: %v", err)
	}
	log.Printf("Updated: %v", updateResp.GetTodo())

	// Get
	getResp, err := client.GetTodo(ctx, &todopb.GetTodoRequest{Id: todoID})
	if err != nil {
		log.Fatalf("GetTodo failed: %v", err)
	}
	log.Printf("Got: %v", getResp.GetTodo())

	// Delete
	_, err = client.DeleteTodo(ctx, &todopb.DeleteTodoRequest{Id: todoID})
	if err != nil {
		log.Fatalf("DeleteTodo failed: %v", err)
	}
	log.Printf("Deleted: %s", todoID)
}
