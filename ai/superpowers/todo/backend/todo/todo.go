// Package todo defines the domain model and (in later phases) an
// in-memory, concurrency-safe store for the superpowers todo learning sample.
// All stored values are passed by value so callers cannot mutate stored
// state through returned references.
package todo

import (
	"errors"
	"fmt"
	"strings"
	"time"
	"unicode/utf8"
)

// Priority is the importance level of a todo item.
type Priority string

const (
	PriorityLow    Priority = "low"
	PriorityMedium Priority = "medium"
	PriorityHigh   Priority = "high"
)

// IsValid reports whether p is one of the defined Priority constants.
func (p Priority) IsValid() bool {
	switch p {
	case PriorityLow, PriorityMedium, PriorityHigh:
		return true
	default:
		return false
	}
}

// Todo is the domain entity. ID, CreatedAt, UpdatedAt are server-managed.
type Todo struct {
	ID        string     `json:"id"`
	Title     string     `json:"title"`
	Completed bool       `json:"completed"`
	Priority  Priority   `json:"priority"`
	DueDate   *time.Time `json:"dueDate"`
	CreatedAt time.Time  `json:"createdAt"`
	UpdatedAt time.Time  `json:"updatedAt"`
}

// NewTodo is the input type for Store.Add. Title is required.
// Priority defaults to PriorityMedium when empty. DueDate is optional.
type NewTodo struct {
	Title    string
	Priority Priority
	DueDate  *time.Time
}

// Validate returns a *ValidationError when input is rejected, otherwise nil.
// Title is trimmed before length check (1-200 chars).
// Priority must be empty (defaulted later) or one of the defined constants.
func (n NewTodo) Validate() error {
	title := strings.TrimSpace(n.Title)
	if title == "" {
		return &ValidationError{Field: "title", Message: "title is required"}
	}
	if utf8.RuneCountInString(title) > 200 {
		return &ValidationError{Field: "title", Message: "title must be at most 200 characters"}
	}
	if n.Priority != "" && !n.Priority.IsValid() {
		return &ValidationError{Field: "priority", Message: fmt.Sprintf("priority %q is invalid", n.Priority)}
	}
	return nil
}

// Patch describes a partial update to a Todo.
// A nil pointer means "field absent in request" (no change).
// ClearDueDate true means the request explicitly set dueDate to null.
type Patch struct {
	Title        *string
	Completed    *bool
	Priority     *Priority
	DueDate      *time.Time
	ClearDueDate bool
}

// Validate returns *ValidationError when patch is invalid.
// An empty patch (no fields set) is considered an error to match the API contract.
func (p Patch) Validate() error {
	if p.Title == nil && p.Completed == nil && p.Priority == nil && p.DueDate == nil && !p.ClearDueDate {
		return &ValidationError{Field: "", Message: "request body must contain at least one field"}
	}
	if p.Title != nil {
		t := strings.TrimSpace(*p.Title)
		if t == "" {
			return &ValidationError{Field: "title", Message: "title is required"}
		}
		if utf8.RuneCountInString(t) > 200 {
			return &ValidationError{Field: "title", Message: "title must be at most 200 characters"}
		}
	}
	if p.Priority != nil && !p.Priority.IsValid() {
		return &ValidationError{Field: "priority", Message: fmt.Sprintf("priority %q is invalid", *p.Priority)}
	}
	return nil
}

// ValidationError indicates a field-level validation failure.
// Field may be empty for whole-body errors.
type ValidationError struct {
	Field   string
	Message string
}

// Error implements error.
func (e *ValidationError) Error() string { return e.Message }

// ErrNotFound is returned by Store.Update and Store.Delete when the id is unknown.
var ErrNotFound = errors.New("todo not found")
