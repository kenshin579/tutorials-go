package server

import (
	"context"
	"sync"
	"time"

	todopb "github.com/kenshin579/tutorials-go/grpc/todolist/gen/todo/v1"
	"github.com/google/uuid"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type TodoService struct {
	todopb.UnimplementedTodoServiceServer
	mu    sync.Mutex
	todos map[string]*todopb.Todo
}

func NewTodoService() *TodoService {
	return &TodoService{
		todos: make(map[string]*todopb.Todo),
	}
}

func (s *TodoService) CreateTodo(ctx context.Context, req *todopb.CreateTodoRequest) (*todopb.CreateTodoResponse, error) {
	if req.GetTitle() == "" {
		return nil, status.Error(codes.InvalidArgument, "title is required")
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	todo := &todopb.Todo{
		Id:        uuid.New().String(),
		Title:     req.GetTitle(),
		Completed: false,
		CreatedAt: timestamppb.New(time.Now()),
	}
	s.todos[todo.Id] = todo

	return &todopb.CreateTodoResponse{Todo: todo}, nil
}

func (s *TodoService) GetTodo(ctx context.Context, req *todopb.GetTodoRequest) (*todopb.GetTodoResponse, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	todo, ok := s.todos[req.GetId()]
	if !ok {
		return nil, status.Errorf(codes.NotFound, "todo %q not found", req.GetId())
	}

	return &todopb.GetTodoResponse{Todo: todo}, nil
}

func (s *TodoService) ListTodos(ctx context.Context, req *todopb.ListTodosRequest) (*todopb.ListTodosResponse, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	todos := make([]*todopb.Todo, 0, len(s.todos))
	for _, todo := range s.todos {
		todos = append(todos, todo)
	}

	return &todopb.ListTodosResponse{Todos: todos}, nil
}

func (s *TodoService) UpdateTodo(ctx context.Context, req *todopb.UpdateTodoRequest) (*todopb.UpdateTodoResponse, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	todo, ok := s.todos[req.GetId()]
	if !ok {
		return nil, status.Errorf(codes.NotFound, "todo %q not found", req.GetId())
	}

	if req.GetTitle() != "" {
		todo.Title = req.GetTitle()
	}
	todo.Completed = req.GetCompleted()

	return &todopb.UpdateTodoResponse{Todo: todo}, nil
}

func (s *TodoService) DeleteTodo(ctx context.Context, req *todopb.DeleteTodoRequest) (*emptypb.Empty, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	if _, ok := s.todos[req.GetId()]; !ok {
		return nil, status.Errorf(codes.NotFound, "todo %q not found", req.GetId())
	}

	delete(s.todos, req.GetId())
	return &emptypb.Empty{}, nil
}
