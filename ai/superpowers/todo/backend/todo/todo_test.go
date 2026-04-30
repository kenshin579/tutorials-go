package todo

import (
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestPriority_IsValid(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name string
		p    Priority
		want bool
	}{
		{"low", PriorityLow, true},
		{"medium", PriorityMedium, true},
		{"high", PriorityHigh, true},
		{"empty", Priority(""), false},
		{"unknown", Priority("urgent"), false},
	}
	for _, tc := range tests {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			assert.Equal(t, tc.want, tc.p.IsValid())
		})
	}
}

func TestNewTodo_Validate(t *testing.T) {
	t.Parallel()
	future := time.Now().Add(24 * time.Hour)
	tooLong := strings.Repeat("a", 201)

	tests := []struct {
		name      string
		input     NewTodo
		wantField string // 빈 문자열이면 통과 기대
	}{
		{"title only", NewTodo{Title: "buy milk"}, ""},
		{"title with priority", NewTodo{Title: "x", Priority: PriorityHigh}, ""},
		{"title with future due", NewTodo{Title: "x", DueDate: &future}, ""},
		{"empty title", NewTodo{Title: ""}, "title"},
		{"whitespace title", NewTodo{Title: "   "}, "title"},
		{"title length 201", NewTodo{Title: tooLong}, "title"},
		{"korean title at boundary", NewTodo{Title: strings.Repeat("가", 200)}, ""},
		{"korean title over boundary", NewTodo{Title: strings.Repeat("가", 201)}, "title"},
		{"invalid priority", NewTodo{Title: "x", Priority: Priority("urgent")}, "priority"},
	}
	for _, tc := range tests {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			err := tc.input.Validate()
			if tc.wantField == "" {
				assert.NoError(t, err)
				return
			}
			var verr *ValidationError
			if assert.ErrorAs(t, err, &verr) {
				assert.Equal(t, tc.wantField, verr.Field)
			}
		})
	}
}

func TestNewTodo_Validate_TitleLengthBoundaries(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name      string
		titleLen  int
		wantValid bool
	}{
		{"len 1", 1, true},
		{"len 200", 200, true},
		{"len 0", 0, false},
		{"len 201", 201, false},
	}
	for _, tc := range tests {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			title := strings.Repeat("a", tc.titleLen)
			err := NewTodo{Title: title}.Validate()
			if tc.wantValid {
				assert.NoError(t, err)
			} else {
				assert.Error(t, err)
			}
		})
	}
}
