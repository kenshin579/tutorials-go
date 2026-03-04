package todolist_test

import (
	"context"
	"log"
	"net"
	"testing"

	todopb "github.com/kenshin579/tutorials-go/grpc/todolist/gen/todo/v1"
	"github.com/kenshin579/tutorials-go/grpc/todolist/server"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/status"
	"google.golang.org/grpc/test/bufconn"
)

const bufSize = 1024 * 1024

var lis *bufconn.Listener

func init() {
	lis = bufconn.Listen(bufSize)
	s := grpc.NewServer()
	todopb.RegisterTodoServiceServer(s, server.NewTodoService())
	go func() {
		if err := s.Serve(lis); err != nil {
			log.Fatalf("server exited with error: %v", err)
		}
	}()
}

func bufDialer(ctx context.Context, _ string) (net.Conn, error) {
	return lis.DialContext(ctx)
}

func newTestClient(t *testing.T) todopb.TodoServiceClient {
	t.Helper()
	conn, err := grpc.NewClient("passthrough:///bufnet",
		grpc.WithContextDialer(bufDialer),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	require.NoError(t, err)
	t.Cleanup(func() { conn.Close() })
	return todopb.NewTodoServiceClient(conn)
}

func TestTodoService_CRUD(t *testing.T) {
	client := newTestClient(t)
	ctx := context.Background()

	// Create
	createResp, err := client.CreateTodo(ctx, &todopb.CreateTodoRequest{Title: "Test Todo"})
	require.NoError(t, err)
	assert.Equal(t, "Test Todo", createResp.GetTodo().GetTitle())
	assert.False(t, createResp.GetTodo().GetCompleted())

	todoID := createResp.GetTodo().GetId()

	// Get
	getResp, err := client.GetTodo(ctx, &todopb.GetTodoRequest{Id: todoID})
	require.NoError(t, err)
	assert.Equal(t, "Test Todo", getResp.GetTodo().GetTitle())

	// List
	listResp, err := client.ListTodos(ctx, &todopb.ListTodosRequest{})
	require.NoError(t, err)
	assert.GreaterOrEqual(t, len(listResp.GetTodos()), 1)

	// Update
	updateResp, err := client.UpdateTodo(ctx, &todopb.UpdateTodoRequest{
		Id:        todoID,
		Title:     "Updated Todo",
		Completed: true,
	})
	require.NoError(t, err)
	assert.Equal(t, "Updated Todo", updateResp.GetTodo().GetTitle())
	assert.True(t, updateResp.GetTodo().GetCompleted())

	// Delete
	_, err = client.DeleteTodo(ctx, &todopb.DeleteTodoRequest{Id: todoID})
	require.NoError(t, err)

	// Get after delete → NotFound
	_, err = client.GetTodo(ctx, &todopb.GetTodoRequest{Id: todoID})
	require.Error(t, err)
	st, ok := status.FromError(err)
	assert.True(t, ok)
	assert.Equal(t, codes.NotFound, st.Code())
}

func TestTodoService_CreateTodo_EmptyTitle(t *testing.T) {
	client := newTestClient(t)
	ctx := context.Background()

	_, err := client.CreateTodo(ctx, &todopb.CreateTodoRequest{Title: ""})
	require.Error(t, err)
	st, ok := status.FromError(err)
	assert.True(t, ok)
	assert.Equal(t, codes.InvalidArgument, st.Code())
}

func TestTodoService_GetTodo_NotFound(t *testing.T) {
	client := newTestClient(t)
	ctx := context.Background()

	_, err := client.GetTodo(ctx, &todopb.GetTodoRequest{Id: "nonexistent"})
	require.Error(t, err)
	st, ok := status.FromError(err)
	assert.True(t, ok)
	assert.Equal(t, codes.NotFound, st.Code())
}
