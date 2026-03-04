package backend

import (
	"context"
	"fmt"
	"time"

	"github.com/wailsapp/wails/v2/pkg/runtime"
)

type App struct {
	ctx   context.Context
	store *TodoStore
}

func NewApp() *App {
	return &App{}
}

func (a *App) Startup(ctx context.Context) {
	a.ctx = ctx
	a.store = NewTodoStore("todos.json")
}

// GetTodos returns all todos
func (a *App) GetTodos() []Todo {
	return a.store.Todos
}

// AddTodo creates a new todo
func (a *App) AddTodo(title string) Todo {
	todo := Todo{
		ID:        fmt.Sprintf("%d", time.Now().UnixNano()),
		Title:     title,
		Done:      false,
		CreatedAt: time.Now(),
	}
	a.store.Todos = append(a.store.Todos, todo)
	a.store.Save()
	return todo
}

// ToggleTodo toggles the done status
func (a *App) ToggleTodo(id string) []Todo {
	for i, t := range a.store.Todos {
		if t.ID == id {
			a.store.Todos[i].Done = !a.store.Todos[i].Done
			break
		}
	}
	a.store.Save()
	return a.store.Todos
}

// DeleteTodo deletes a todo with confirmation dialog
func (a *App) DeleteTodo(id string) []Todo {
	for i, t := range a.store.Todos {
		if t.ID == id {
			result, err := runtime.MessageDialog(a.ctx, runtime.MessageDialogOptions{
				Type:          runtime.QuestionDialog,
				Title:         "삭제 확인",
				Message:       fmt.Sprintf("'%s'을(를) 삭제하시겠습니까?", t.Title),
				DefaultButton: "No",
			})
			if err == nil && result == "Yes" {
				a.store.Todos = append(a.store.Todos[:i], a.store.Todos[i+1:]...)
				a.store.Save()
			}
			break
		}
	}
	return a.store.Todos
}

// ExportTodos exports todos to a file via save dialog
func (a *App) ExportTodos() error {
	path, err := runtime.SaveFileDialog(a.ctx, runtime.SaveDialogOptions{
		Title:           "Todo 목록 내보내기",
		DefaultFilename: "todos.json",
		Filters: []runtime.FileFilter{
			{DisplayName: "JSON Files", Pattern: "*.json"},
		},
	})
	if err != nil || path == "" {
		return err
	}
	return a.store.ExportToFile(path)
}

// ImportTodos imports todos from a file via open dialog
func (a *App) ImportTodos() ([]Todo, error) {
	path, err := runtime.OpenFileDialog(a.ctx, runtime.OpenDialogOptions{
		Title: "Todo 목록 불러오기",
		Filters: []runtime.FileFilter{
			{DisplayName: "JSON Files", Pattern: "*.json"},
		},
	})
	if err != nil || path == "" {
		return a.store.Todos, err
	}
	if err := a.store.LoadFromFile(path); err != nil {
		return a.store.Todos, err
	}
	return a.store.Todos, nil
}
