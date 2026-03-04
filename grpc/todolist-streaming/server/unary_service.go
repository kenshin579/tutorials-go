package server

import (
	"context"
	"sync"
	"time"

	"github.com/google/uuid"
	todopb "github.com/kenshin579/tutorials-go/grpc/todolist-streaming/gen/todo/v1"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
)

// TodoUnaryService - 벤치마크 비교용 Unary 서비스
type TodoUnaryService struct {
	todopb.UnimplementedTodoUnaryServiceServer
	mu    sync.Mutex
	todos map[string]*todopb.Todo
}

func NewTodoUnaryService() *TodoUnaryService {
	return &TodoUnaryService{
		todos: make(map[string]*todopb.Todo),
	}
}

func (s *TodoUnaryService) CreateTodo(ctx context.Context, req *todopb.CreateTodoRequest) (*todopb.Todo, error) {
	if req.GetTitle() == "" {
		return nil, status.Error(codes.InvalidArgument, "title is required")
	}

	todo := &todopb.Todo{
		Id:        uuid.New().String(),
		Title:     req.GetTitle(),
		Completed: false,
		CreatedAt: timestamppb.New(time.Now()),
	}

	s.mu.Lock()
	s.todos[todo.Id] = todo
	s.mu.Unlock()

	return todo, nil
}

func (s *TodoUnaryService) GetTodos(ctx context.Context, req *todopb.ListTodosRequest) (*todopb.TodoList, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	todos := make([]*todopb.Todo, 0, len(s.todos))
	for _, todo := range s.todos {
		if req.GetCompletedOnly() && !todo.GetCompleted() {
			continue
		}
		todos = append(todos, todo)
	}

	return &todopb.TodoList{Todos: todos}, nil
}

func (s *TodoUnaryService) AddTestData(todos []*todopb.Todo) {
	s.mu.Lock()
	defer s.mu.Unlock()
	for _, todo := range todos {
		s.todos[todo.Id] = todo
	}
}
