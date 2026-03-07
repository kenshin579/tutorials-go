package backend

import (
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

func newTestApp(t *testing.T) *App {
	t.Helper()
	tmpFile := filepath.Join(t.TempDir(), "todos.json")
	app := &App{
		store: NewTodoStore(tmpFile),
	}
	return app
}

func TestApp_AddTodo(t *testing.T) {
	app := newTestApp(t)

	todo := app.AddTodo("새 할 일")

	assert.Equal(t, "새 할 일", todo.Title)
	assert.False(t, todo.Done)
	assert.NotEmpty(t, todo.ID)
	assert.Len(t, app.GetTodos(), 1)
}

func TestApp_AddMultipleTodos(t *testing.T) {
	app := newTestApp(t)

	app.AddTodo("할 일 1")
	app.AddTodo("할 일 2")
	app.AddTodo("할 일 3")

	assert.Len(t, app.GetTodos(), 3)
}

func TestApp_GetTodos_Empty(t *testing.T) {
	app := newTestApp(t)

	todos := app.GetTodos()

	assert.Empty(t, todos)
}

func TestApp_ToggleTodo(t *testing.T) {
	app := newTestApp(t)
	todo := app.AddTodo("토글 테스트")

	// Done: false -> true
	todos := app.ToggleTodo(todo.ID)
	assert.True(t, todos[0].Done)

	// Done: true -> false
	todos = app.ToggleTodo(todo.ID)
	assert.False(t, todos[0].Done)
}

func TestApp_ToggleTodo_NonExistentID(t *testing.T) {
	app := newTestApp(t)
	app.AddTodo("기존 할 일")

	todos := app.ToggleTodo("non-existent-id")

	assert.Len(t, todos, 1)
	assert.False(t, todos[0].Done)
}

func TestApp_AddTodo_PersistsToFile(t *testing.T) {
	tmpFile := filepath.Join(t.TempDir(), "todos.json")
	app := &App{store: NewTodoStore(tmpFile)}

	app.AddTodo("영속성 테스트")

	// 같은 파일로 새 store를 로드하여 저장 확인
	reloaded := NewTodoStore(tmpFile)
	assert.Len(t, reloaded.Todos, 1)
	assert.Equal(t, "영속성 테스트", reloaded.Todos[0].Title)
}
