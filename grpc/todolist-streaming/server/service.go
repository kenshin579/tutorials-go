package server

import (
	"io"
	"sync"
	"time"

	"github.com/google/uuid"
	todopb "github.com/kenshin579/tutorials-go/grpc/todolist-streaming/gen/todo/v1"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type TodoStreamingService struct {
	todopb.UnimplementedTodoStreamingServiceServer
	mu    sync.Mutex
	todos map[string]*todopb.Todo
}

func NewTodoStreamingService() *TodoStreamingService {
	return &TodoStreamingService{
		todos: make(map[string]*todopb.Todo),
	}
}

// ListTodos - Server Streaming: 필터 조건에 맞는 Todo를 하나씩 스트림 전송
func (s *TodoStreamingService) ListTodos(req *todopb.ListTodosRequest, stream todopb.TodoStreamingService_ListTodosServer) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	for _, todo := range s.todos {
		if req.GetCompletedOnly() && !todo.GetCompleted() {
			continue
		}
		if err := stream.Send(todo); err != nil {
			return status.Errorf(codes.Internal, "failed to send todo: %v", err)
		}
	}
	return nil
}

// BatchCreateTodos - Client Streaming: 클라이언트가 여러 Todo를 보내면 배치로 생성
func (s *TodoStreamingService) BatchCreateTodos(stream todopb.TodoStreamingService_BatchCreateTodosServer) error {
	var ids []string

	for {
		req, err := stream.Recv()
		if err == io.EOF {
			return stream.SendAndClose(&todopb.BatchCreateResponse{
				CreatedCount: int32(len(ids)),
				Ids:          ids,
			})
		}
		if err != nil {
			return status.Errorf(codes.Internal, "failed to receive: %v", err)
		}

		if req.GetTitle() == "" {
			return status.Error(codes.InvalidArgument, "title is required")
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

		ids = append(ids, todo.Id)
	}
}

// TodoUpdates - Bidirectional Streaming: 클라이언트가 작업을 보내면 서버가 실시간 이벤트 응답
func (s *TodoStreamingService) TodoUpdates(stream todopb.TodoStreamingService_TodoUpdatesServer) error {
	for {
		action, err := stream.Recv()
		if err == io.EOF {
			return nil
		}
		if err != nil {
			return status.Errorf(codes.Internal, "failed to receive: %v", err)
		}

		event := s.processAction(action)
		if err := stream.Send(event); err != nil {
			return status.Errorf(codes.Internal, "failed to send event: %v", err)
		}
	}
}

func (s *TodoStreamingService) processAction(action *todopb.TodoAction) *todopb.TodoEvent {
	switch action.GetAction() {
	case todopb.ActionType_ACTION_TYPE_CREATE:
		todo := &todopb.Todo{
			Id:        uuid.New().String(),
			Title:     action.GetTitle(),
			Completed: false,
			CreatedAt: timestamppb.New(time.Now()),
		}
		s.mu.Lock()
		s.todos[todo.Id] = todo
		s.mu.Unlock()

		return &todopb.TodoEvent{
			Event:   todopb.EventType_EVENT_TYPE_CREATED,
			TodoId:  todo.Id,
			Message: "created: " + todo.Title,
		}

	case todopb.ActionType_ACTION_TYPE_COMPLETE:
		s.mu.Lock()
		todo, ok := s.todos[action.GetTodoId()]
		if !ok {
			s.mu.Unlock()
			return &todopb.TodoEvent{
				Event:   todopb.EventType_EVENT_TYPE_ERROR,
				TodoId:  action.GetTodoId(),
				Message: "todo not found",
			}
		}
		todo.Completed = true
		s.mu.Unlock()

		return &todopb.TodoEvent{
			Event:   todopb.EventType_EVENT_TYPE_COMPLETED,
			TodoId:  todo.Id,
			Message: "completed: " + todo.Title,
		}

	case todopb.ActionType_ACTION_TYPE_DELETE:
		s.mu.Lock()
		todo, ok := s.todos[action.GetTodoId()]
		if !ok {
			s.mu.Unlock()
			return &todopb.TodoEvent{
				Event:   todopb.EventType_EVENT_TYPE_ERROR,
				TodoId:  action.GetTodoId(),
				Message: "todo not found",
			}
		}
		delete(s.todos, action.GetTodoId())
		s.mu.Unlock()

		return &todopb.TodoEvent{
			Event:   todopb.EventType_EVENT_TYPE_DELETED,
			TodoId:  todo.Id,
			Message: "deleted: " + todo.Title,
		}

	default:
		return &todopb.TodoEvent{
			Event:   todopb.EventType_EVENT_TYPE_ERROR,
			Message: "unknown action type",
		}
	}
}

// AddTestData 는 테스트/벤치마크에서 초기 데이터를 추가할 때 사용
func (s *TodoStreamingService) AddTestData(todos []*todopb.Todo) {
	s.mu.Lock()
	defer s.mu.Unlock()
	for _, todo := range todos {
		s.todos[todo.Id] = todo
	}
}
