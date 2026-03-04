package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"time"

	todopb "github.com/kenshin579/tutorials-go/grpc/todolist-streaming/gen/todo/v1"
	"golang.org/x/sync/errgroup"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	conn, err := grpc.NewClient("localhost:50052",
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		log.Fatalf("failed to connect: %v", err)
	}
	defer conn.Close()

	client := todopb.NewTodoStreamingServiceClient(conn)
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// 1. Client Streaming: 배치 생성
	fmt.Println("=== Client Streaming: BatchCreateTodos ===")
	demoClientStreaming(ctx, client)

	// 2. Server Streaming: 목록 조회
	fmt.Println("\n=== Server Streaming: ListTodos ===")
	demoServerStreaming(ctx, client)

	// 3. Bidirectional Streaming: 실시간 작업
	fmt.Println("\n=== Bidirectional Streaming: TodoUpdates ===")
	demoBidiStreaming(ctx, client)
}

func demoClientStreaming(ctx context.Context, client todopb.TodoStreamingServiceClient) {
	stream, err := client.BatchCreateTodos(ctx)
	if err != nil {
		log.Fatalf("BatchCreateTodos failed: %v", err)
	}

	titles := []string{"Learn gRPC", "Write tests", "Deploy service"}
	for _, title := range titles {
		if err := stream.Send(&todopb.CreateTodoRequest{Title: title}); err != nil {
			log.Fatalf("Send failed: %v", err)
		}
		log.Printf("Sent: %s", title)
	}

	resp, err := stream.CloseAndRecv()
	if err != nil {
		log.Fatalf("CloseAndRecv failed: %v", err)
	}
	log.Printf("Created %d todos: %v", resp.GetCreatedCount(), resp.GetIds())
}

func demoServerStreaming(ctx context.Context, client todopb.TodoStreamingServiceClient) {
	stream, err := client.ListTodos(ctx, &todopb.ListTodosRequest{})
	if err != nil {
		log.Fatalf("ListTodos failed: %v", err)
	}

	for {
		todo, err := stream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatalf("Recv failed: %v", err)
		}
		log.Printf("Received: id=%s title=%s completed=%v", todo.GetId(), todo.GetTitle(), todo.GetCompleted())
	}
}

func demoBidiStreaming(ctx context.Context, client todopb.TodoStreamingServiceClient) {
	stream, err := client.TodoUpdates(ctx)
	if err != nil {
		log.Fatalf("TodoUpdates failed: %v", err)
	}

	// errgroup으로 send/recv 고루틴 관리
	g, ctx := errgroup.WithContext(ctx)

	// recv 고루틴
	g.Go(func() error {
		for {
			event, err := stream.Recv()
			if err == io.EOF {
				return nil
			}
			if err != nil {
				return fmt.Errorf("recv error: %w", err)
			}
			log.Printf("Event: type=%s id=%s message=%s",
				event.GetEvent(), event.GetTodoId(), event.GetMessage())
		}
	})

	// send 고루틴
	g.Go(func() error {
		actions := []*todopb.TodoAction{
			{Action: todopb.ActionType_ACTION_TYPE_CREATE, Title: "Streaming Todo 1"},
			{Action: todopb.ActionType_ACTION_TYPE_CREATE, Title: "Streaming Todo 2"},
			{Action: todopb.ActionType_ACTION_TYPE_CREATE, Title: "Streaming Todo 3"},
		}

		for _, action := range actions {
			select {
			case <-ctx.Done():
				return ctx.Err()
			default:
				if err := stream.Send(action); err != nil {
					return fmt.Errorf("send error: %w", err)
				}
				log.Printf("Sent action: %s", action.GetAction())
			}
		}
		return stream.CloseSend()
	})

	if err := g.Wait(); err != nil {
		log.Fatalf("Bidi streaming error: %v", err)
	}
}
