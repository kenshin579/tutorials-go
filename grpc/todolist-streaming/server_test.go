package todolist_streaming_test

import (
	"context"
	"fmt"
	"io"
	"log"
	"net"
	"testing"

	todopb "github.com/kenshin579/tutorials-go/grpc/todolist-streaming/gen/todo/v1"
	"github.com/kenshin579/tutorials-go/grpc/todolist-streaming/server"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"golang.org/x/sync/errgroup"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/test/bufconn"
)

const bufSize = 1024 * 1024

var lis *bufconn.Listener

func init() {
	lis = bufconn.Listen(bufSize)
	s := grpc.NewServer()
	todopb.RegisterTodoStreamingServiceServer(s, server.NewTodoStreamingService())
	go func() {
		if err := s.Serve(lis); err != nil {
			log.Fatalf("server exited with error: %v", err)
		}
	}()
}

func bufDialer(ctx context.Context, _ string) (net.Conn, error) {
	return lis.DialContext(ctx)
}

func newTestClient(t *testing.T) todopb.TodoStreamingServiceClient {
	t.Helper()
	conn, err := grpc.NewClient("passthrough:///bufnet",
		grpc.WithContextDialer(bufDialer),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	require.NoError(t, err)
	t.Cleanup(func() { conn.Close() })
	return todopb.NewTodoStreamingServiceClient(conn)
}

// --- Server Streaming 테스트 ---

func TestServerStreaming_ListTodos(t *testing.T) {
	client := newTestClient(t)
	ctx := context.Background()

	// 먼저 Client Streaming으로 데이터 생성
	batchStream, err := client.BatchCreateTodos(ctx)
	require.NoError(t, err)

	titles := []string{"Todo A", "Todo B", "Todo C"}
	for _, title := range titles {
		require.NoError(t, batchStream.Send(&todopb.CreateTodoRequest{Title: title}))
	}
	batchResp, err := batchStream.CloseAndRecv()
	require.NoError(t, err)
	assert.Equal(t, int32(3), batchResp.GetCreatedCount())

	// Server Streaming으로 조회
	listStream, err := client.ListTodos(ctx, &todopb.ListTodosRequest{})
	require.NoError(t, err)

	var received []*todopb.Todo
	for {
		todo, err := listStream.Recv()
		if err == io.EOF {
			break
		}
		require.NoError(t, err)
		received = append(received, todo)
	}

	assert.GreaterOrEqual(t, len(received), 3)
}

func TestServerStreaming_ListTodos_CompletedOnly(t *testing.T) {
	client := newTestClient(t)
	ctx := context.Background()

	// completed_only 필터로 조회 (완료된 것이 없으므로 빈 결과)
	stream, err := client.ListTodos(ctx, &todopb.ListTodosRequest{CompletedOnly: true})
	require.NoError(t, err)

	var count int
	for {
		todo, err := stream.Recv()
		if err == io.EOF {
			break
		}
		require.NoError(t, err)
		assert.True(t, todo.GetCompleted())
		count++
	}
	// 아직 완료된 Todo가 없을 수 있음
	t.Logf("completed todos: %d", count)
}

// --- Client Streaming 테스트 ---

func TestClientStreaming_BatchCreateTodos(t *testing.T) {
	client := newTestClient(t)
	ctx := context.Background()

	stream, err := client.BatchCreateTodos(ctx)
	require.NoError(t, err)

	titles := []string{"Batch 1", "Batch 2", "Batch 3", "Batch 4", "Batch 5"}
	for _, title := range titles {
		require.NoError(t, stream.Send(&todopb.CreateTodoRequest{Title: title}))
	}

	resp, err := stream.CloseAndRecv()
	require.NoError(t, err)
	assert.Equal(t, int32(5), resp.GetCreatedCount())
	assert.Len(t, resp.GetIds(), 5)
}

func TestClientStreaming_BatchCreateTodos_EmptyTitle(t *testing.T) {
	client := newTestClient(t)
	ctx := context.Background()

	stream, err := client.BatchCreateTodos(ctx)
	require.NoError(t, err)

	// 빈 title 전송
	require.NoError(t, stream.Send(&todopb.CreateTodoRequest{Title: ""}))

	_, err = stream.CloseAndRecv()
	require.Error(t, err)
}

// --- Bidirectional Streaming 테스트 ---

func TestBidiStreaming_TodoUpdates(t *testing.T) {
	client := newTestClient(t)
	ctx := context.Background()

	stream, err := client.TodoUpdates(ctx)
	require.NoError(t, err)

	g, ctx := errgroup.WithContext(ctx)
	var events []*todopb.TodoEvent

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
			events = append(events, event)
		}
	})

	// send 고루틴
	g.Go(func() error {
		actions := []*todopb.TodoAction{
			{Action: todopb.ActionType_ACTION_TYPE_CREATE, Title: "Bidi Todo 1"},
			{Action: todopb.ActionType_ACTION_TYPE_CREATE, Title: "Bidi Todo 2"},
		}
		for _, action := range actions {
			if err := stream.Send(action); err != nil {
				return fmt.Errorf("send error: %w", err)
			}
		}
		return stream.CloseSend()
	})

	require.NoError(t, g.Wait())
	assert.Len(t, events, 2)
	assert.Equal(t, todopb.EventType_EVENT_TYPE_CREATED, events[0].GetEvent())
	assert.Equal(t, todopb.EventType_EVENT_TYPE_CREATED, events[1].GetEvent())
}

func TestBidiStreaming_TodoUpdates_CreateAndComplete(t *testing.T) {
	client := newTestClient(t)
	ctx := context.Background()

	stream, err := client.TodoUpdates(ctx)
	require.NoError(t, err)

	// CREATE → 응답 받기 → COMPLETE 순서로 진행 (순차적)
	require.NoError(t, stream.Send(&todopb.TodoAction{
		Action: todopb.ActionType_ACTION_TYPE_CREATE,
		Title:  "Complete Me",
	}))

	event, err := stream.Recv()
	require.NoError(t, err)
	assert.Equal(t, todopb.EventType_EVENT_TYPE_CREATED, event.GetEvent())
	todoID := event.GetTodoId()

	// COMPLETE
	require.NoError(t, stream.Send(&todopb.TodoAction{
		Action: todopb.ActionType_ACTION_TYPE_COMPLETE,
		TodoId: todoID,
	}))

	event, err = stream.Recv()
	require.NoError(t, err)
	assert.Equal(t, todopb.EventType_EVENT_TYPE_COMPLETED, event.GetEvent())

	// DELETE
	require.NoError(t, stream.Send(&todopb.TodoAction{
		Action: todopb.ActionType_ACTION_TYPE_DELETE,
		TodoId: todoID,
	}))

	event, err = stream.Recv()
	require.NoError(t, err)
	assert.Equal(t, todopb.EventType_EVENT_TYPE_DELETED, event.GetEvent())

	require.NoError(t, stream.CloseSend())
}

func TestBidiStreaming_TodoUpdates_NotFound(t *testing.T) {
	client := newTestClient(t)
	ctx := context.Background()

	stream, err := client.TodoUpdates(ctx)
	require.NoError(t, err)

	// 존재하지 않는 Todo 완료 시도
	require.NoError(t, stream.Send(&todopb.TodoAction{
		Action: todopb.ActionType_ACTION_TYPE_COMPLETE,
		TodoId: "nonexistent-id",
	}))

	event, err := stream.Recv()
	require.NoError(t, err)
	assert.Equal(t, todopb.EventType_EVENT_TYPE_ERROR, event.GetEvent())
	assert.Contains(t, event.GetMessage(), "not found")

	require.NoError(t, stream.CloseSend())
}

// --- Directional Channel 패턴 테스트 ---

func TestBidiStreaming_DirectionalChannel(t *testing.T) {
	client := newTestClient(t)
	ctx := context.Background()

	stream, err := client.TodoUpdates(ctx)
	require.NoError(t, err)

	// directional channel 패턴
	eventCh := make(chan *todopb.TodoEvent, 10)
	errCh := make(chan error, 1)

	// recv 고루틴 (<-chan 방향으로만 전달)
	go func() {
		defer close(eventCh)
		for {
			event, err := stream.Recv()
			if err == io.EOF {
				return
			}
			if err != nil {
				errCh <- err
				return
			}
			eventCh <- event
		}
	}()

	// send (chan<- 방향)
	go func() {
		actions := []*todopb.TodoAction{
			{Action: todopb.ActionType_ACTION_TYPE_CREATE, Title: "Chan Todo 1"},
			{Action: todopb.ActionType_ACTION_TYPE_CREATE, Title: "Chan Todo 2"},
		}
		for _, action := range actions {
			if err := stream.Send(action); err != nil {
				errCh <- err
				return
			}
		}
		stream.CloseSend()
	}()

	// eventCh에서 이벤트 수신
	var received []*todopb.TodoEvent
	for event := range eventCh {
		received = append(received, event)
	}

	select {
	case err := <-errCh:
		require.NoError(t, err)
	default:
	}

	assert.Len(t, received, 2)
	for _, event := range received {
		assert.Equal(t, todopb.EventType_EVENT_TYPE_CREATED, event.GetEvent())
	}
}
