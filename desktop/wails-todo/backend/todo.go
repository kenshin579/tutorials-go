package backend

import (
	"encoding/json"
	"os"
	"time"
)

type Todo struct {
	ID        string    `json:"id"`
	Title     string    `json:"title"`
	Done      bool      `json:"done"`
	CreatedAt time.Time `json:"createdAt"`
}

type TodoStore struct {
	filePath string
	Todos    []Todo `json:"todos"`
}

func NewTodoStore(filePath string) *TodoStore {
	store := &TodoStore{filePath: filePath}
	store.Load()
	return store
}

func (s *TodoStore) Load() {
	data, err := os.ReadFile(s.filePath)
	if err != nil {
		s.Todos = []Todo{}
		return
	}
	if err := json.Unmarshal(data, &s.Todos); err != nil {
		s.Todos = []Todo{}
	}
}

func (s *TodoStore) Save() error {
	data, err := json.MarshalIndent(s.Todos, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(s.filePath, data, 0644)
}

func (s *TodoStore) LoadFromFile(path string) error {
	data, err := os.ReadFile(path)
	if err != nil {
		return err
	}
	if err := json.Unmarshal(data, &s.Todos); err != nil {
		return err
	}
	return s.Save()
}

func (s *TodoStore) ExportToFile(path string) error {
	data, err := json.MarshalIndent(s.Todos, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(path, data, 0644)
}
