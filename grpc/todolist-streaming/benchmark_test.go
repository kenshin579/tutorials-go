package todolist_streaming_test

import (
	"context"
	"fmt"
	"io"
	"net"
	"testing"
	"time"

	todopb "github.com/kenshin579/tutorials-go/grpc/todolist-streaming/gen/todo/v1"
	"github.com/kenshin579/tutorials-go/grpc/todolist-streaming/server"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/test/bufconn"
	"google.golang.org/protobuf/types/known/timestamppb"
)

const benchItemCount = 100

// setupBenchServer 벤치마크마다 독립 서버를 생성하여 데이터 누적 문제 방지
func setupBenchServer(b *testing.B) (*bufconn.Listener, *server.TodoStreamingService, *server.TodoUnaryService) {
	b.Helper()
	l := bufconn.Listen(bufSize)

	streamingSvc := server.NewTodoStreamingService()
	unarySvc := server.NewTodoUnaryService()

	s := grpc.NewServer()
	todopb.RegisterTodoStreamingServiceServer(s, streamingSvc)
	todopb.RegisterTodoUnaryServiceServer(s, unarySvc)

	go func() {
		s.Serve(l)
	}()
	b.Cleanup(func() { s.Stop() })

	return l, streamingSvc, unarySvc
}

func newBenchClients(b *testing.B, l *bufconn.Listener) (todopb.TodoStreamingServiceClient, todopb.TodoUnaryServiceClient) {
	b.Helper()
	dialer := func(ctx context.Context, _ string) (net.Conn, error) {
		return l.DialContext(ctx)
	}
	conn, err := grpc.NewClient("passthrough:///bufnet",
		grpc.WithContextDialer(dialer),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		b.Fatalf("failed to dial: %v", err)
	}
	b.Cleanup(func() { conn.Close() })
	return todopb.NewTodoStreamingServiceClient(conn), todopb.NewTodoUnaryServiceClient(conn)
}

func generateTestTodos(count int) []*todopb.Todo {
	todos := make([]*todopb.Todo, count)
	for i := 0; i < count; i++ {
		todos[i] = &todopb.Todo{
			Id:        fmt.Sprintf("todo-%d", i),
			Title:     fmt.Sprintf("Benchmark Todo %d", i),
			Completed: false,
			CreatedAt: timestamppb.New(time.Now()),
		}
	}
	return todos
}

// --- Create 벤치마크: Unary 반복 호출 vs Client Streaming ---

func BenchmarkUnaryCreateTodos(b *testing.B) {
	l, _, _ := setupBenchServer(b)
	_, unaryClient := newBenchClients(b, l)
	ctx := context.Background()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		for j := 0; j < benchItemCount; j++ {
			_, err := unaryClient.CreateTodo(ctx, &todopb.CreateTodoRequest{
				Title: fmt.Sprintf("bench-%d-%d", i, j),
			})
			if err != nil {
				b.Fatalf("CreateTodo failed: %v", err)
			}
		}
	}
}

func BenchmarkStreamingCreateTodos(b *testing.B) {
	l, _, _ := setupBenchServer(b)
	streamingClient, _ := newBenchClients(b, l)
	ctx := context.Background()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		stream, err := streamingClient.BatchCreateTodos(ctx)
		if err != nil {
			b.Fatalf("BatchCreateTodos failed: %v", err)
		}

		for j := 0; j < benchItemCount; j++ {
			if err := stream.Send(&todopb.CreateTodoRequest{
				Title: fmt.Sprintf("bench-%d-%d", i, j),
			}); err != nil {
				b.Fatalf("Send failed: %v", err)
			}
		}

		if _, err := stream.CloseAndRecv(); err != nil {
			b.Fatalf("CloseAndRecv failed: %v", err)
		}
	}
}

// --- Get 벤치마크: Unary 전체 조회 vs Server Streaming ---

func BenchmarkUnaryGetTodos(b *testing.B) {
	l, streamingSvc, unarySvc := setupBenchServer(b)
	_, unaryClient := newBenchClients(b, l)
	ctx := context.Background()

	// 사전 데이터 로드
	testData := generateTestTodos(benchItemCount)
	streamingSvc.AddTestData(testData)
	unarySvc.AddTestData(testData)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := unaryClient.GetTodos(ctx, &todopb.ListTodosRequest{})
		if err != nil {
			b.Fatalf("GetTodos failed: %v", err)
		}
	}
}

func BenchmarkStreamingGetTodos(b *testing.B) {
	l, streamingSvc, _ := setupBenchServer(b)
	streamingClient, _ := newBenchClients(b, l)
	ctx := context.Background()

	// 사전 데이터 로드
	testData := generateTestTodos(benchItemCount)
	streamingSvc.AddTestData(testData)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		stream, err := streamingClient.ListTodos(ctx, &todopb.ListTodosRequest{})
		if err != nil {
			b.Fatalf("ListTodos failed: %v", err)
		}

		for {
			_, err := stream.Recv()
			if err == io.EOF {
				break
			}
			if err != nil {
				b.Fatalf("Recv failed: %v", err)
			}
		}
	}
}
