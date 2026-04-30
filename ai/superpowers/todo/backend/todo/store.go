package todo

import (
	"sort"
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
		DueDate:   copyTimePointer(input.DueDate),
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
		t.DueDate = copyTimePointer(p.DueDate)
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

// StatusFilter narrows List results by completion state.
type StatusFilter string

const (
	StatusAll       StatusFilter = "all"
	StatusActive    StatusFilter = "active"
	StatusCompleted StatusFilter = "completed"
)

// IsValid reports whether s is one of the defined StatusFilter constants.
func (s StatusFilter) IsValid() bool {
	switch s {
	case StatusAll, StatusActive, StatusCompleted:
		return true
	default:
		return false
	}
}

// SortKey selects the field used to order List results.
type SortKey string

const (
	SortCreatedAt SortKey = "createdAt"
	SortDueDate   SortKey = "dueDate"
	SortPriority  SortKey = "priority"
)

// IsValid reports whether k is one of the defined SortKey constants.
func (k SortKey) IsValid() bool {
	switch k {
	case SortCreatedAt, SortDueDate, SortPriority:
		return true
	default:
		return false
	}
}

// OrderDir is ascending or descending.
type OrderDir string

const (
	OrderAsc  OrderDir = "asc"
	OrderDesc OrderDir = "desc"
)

// IsValid reports whether o is one of the defined OrderDir constants.
func (o OrderDir) IsValid() bool {
	switch o {
	case OrderAsc, OrderDesc:
		return true
	default:
		return false
	}
}

// Query is the input to Store.List. Zero values mean defaults:
// Status defaults to StatusAll, Sort to SortCreatedAt, Order to OrderDesc.
type Query struct {
	Status StatusFilter
	Sort   SortKey
	Order  OrderDir
}

func (q Query) withDefaults() Query {
	if q.Status == "" {
		q.Status = StatusAll
	}
	if q.Sort == "" {
		q.Sort = SortCreatedAt
	}
	if q.Order == "" {
		q.Order = OrderDesc
	}
	return q
}

// List returns todos matching the query, sorted as requested.
// nil dueDate values are placed last regardless of Order (deterministic).
// Tie-breaks fall back to createdAt ascending (stable sort).
func (s *Store) List(q Query) []Todo {
	q = q.withDefaults()
	s.mu.RLock()
	out := make([]Todo, 0, len(s.todos))
	for _, t := range s.todos {
		if !matchesStatus(t, q.Status) {
			continue
		}
		out = append(out, t)
	}
	s.mu.RUnlock()
	sortTodos(out, q.Sort, q.Order)
	return out
}

func matchesStatus(t Todo, f StatusFilter) bool {
	switch f {
	case StatusActive:
		return !t.Completed
	case StatusCompleted:
		return t.Completed
	default:
		return true
	}
}

// sortTodos sorts ts in place by the given key/order.
// nil dueDate values always sort last regardless of order (deterministic).
// Tie-breaks fall back to createdAt ascending (stable sort).
func sortTodos(ts []Todo, key SortKey, order OrderDir) {
	asc := order == OrderAsc
	sort.SliceStable(ts, func(i, j int) bool {
		a, b := ts[i], ts[j]

		// dueDate 정렬: nil은 항상 마지막 (asc/desc 무관)
		if key == SortDueDate {
			if a.DueDate == nil && b.DueDate == nil {
				return a.CreatedAt.Before(b.CreatedAt)
			}
			if a.DueDate == nil {
				return false
			}
			if b.DueDate == nil {
				return true
			}
		}

		cmp := primaryCompare(a, b, key)
		if cmp == 0 {
			return a.CreatedAt.Before(b.CreatedAt) // tie-break: createdAt asc
		}
		if asc {
			return cmp < 0
		}
		return cmp > 0
	})
}

func primaryCompare(a, b Todo, key SortKey) int {
	switch key {
	case SortPriority:
		return priorityRank(a.Priority) - priorityRank(b.Priority)
	case SortDueDate:
		// 호출 시점에 둘 다 non-nil 보장됨 (sortTodos에서 nil 분기)
		return compareTime(*a.DueDate, *b.DueDate)
	default: // SortCreatedAt
		return compareTime(a.CreatedAt, b.CreatedAt)
	}
}

func priorityRank(p Priority) int {
	switch p {
	case PriorityLow:
		return 1
	case PriorityMedium:
		return 2
	case PriorityHigh:
		return 3
	default:
		return 0
	}
}

func compareTime(a, b time.Time) int {
	if a.Before(b) {
		return -1
	}
	if a.After(b) {
		return 1
	}
	return 0
}

// copyTimePointer returns a fresh *time.Time pointing to a copy of *t,
// or nil when t is nil. Used in Add/Update to prevent external mutation
// of stored DueDate via pointer aliasing.
func copyTimePointer(t *time.Time) *time.Time {
	if t == nil {
		return nil
	}
	v := *t
	return &v
}
