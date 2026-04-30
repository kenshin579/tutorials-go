package todo

import (
	"strings"
	"sync"
	"time"

	"github.com/google/uuid"
)

// Store is an in-memory, concurrency-safe collection of Todo values.
// All read methods take an RLock; mutations take a Lock.
// Returned values are copies — callers cannot affect stored state by mutation.
type Store struct {
	mu    sync.RWMutex
	todos map[string]Todo
}

// NewStore creates an empty Store ready for use.
func NewStore() *Store {
	return &Store{todos: make(map[string]Todo)}
}

// Add assigns ID/timestamps/default priority, stores the todo, and returns a copy.
// Title is trimmed before storing. Add does NOT call NewTodo.Validate — callers
// must validate at the boundary (handler) before invoking Add.
func (s *Store) Add(input NewTodo) Todo {
	now := time.Now().UTC()
	priority := input.Priority
	if priority == "" {
		priority = PriorityMedium
	}
	t := Todo{
		ID:        uuid.NewString(),
		Title:     strings.TrimSpace(input.Title),
		Completed: false,
		Priority:  priority,
		DueDate:   input.DueDate,
		CreatedAt: now,
		UpdatedAt: now,
	}
	s.mu.Lock()
	s.todos[t.ID] = t
	s.mu.Unlock()
	return t
}

// Get returns the todo with the given id and true, or zero value and false.
func (s *Store) Get(id string) (Todo, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	t, ok := s.todos[id]
	return t, ok
}

// Update applies the patch to the todo identified by id.
// Returns ErrNotFound when id is unknown, or *ValidationError when patch is invalid.
// Validates the patch before applying.
func (s *Store) Update(id string, p Patch) (Todo, error) {
	if err := p.Validate(); err != nil {
		return Todo{}, err
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	t, ok := s.todos[id]
	if !ok {
		return Todo{}, ErrNotFound
	}
	if p.Title != nil {
		t.Title = strings.TrimSpace(*p.Title)
	}
	if p.Completed != nil {
		t.Completed = *p.Completed
	}
	if p.Priority != nil {
		t.Priority = *p.Priority
	}
	if p.ClearDueDate {
		t.DueDate = nil
	} else if p.DueDate != nil {
		t.DueDate = p.DueDate
	}
	t.UpdatedAt = time.Now().UTC()
	s.todos[id] = t
	return t, nil
}

// Delete removes the todo by id. Returns ErrNotFound when id is unknown.
func (s *Store) Delete(id string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, ok := s.todos[id]; !ok {
		return ErrNotFound
	}
	delete(s.todos, id)
	return nil
}
