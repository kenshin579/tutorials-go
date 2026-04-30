package todo

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func newTestStore(t *testing.T) *Store {
	t.Helper()
	return NewStore()
}

func TestStore_Add_AssignsIDAndTimestamps(t *testing.T) {
	t.Parallel()
	s := newTestStore(t)
	before := time.Now()
	got := s.Add(NewTodo{Title: "buy milk"})
	after := time.Now()

	assert.NotEmpty(t, got.ID)
	assert.Equal(t, "buy milk", got.Title)
	assert.False(t, got.Completed)
	assert.Equal(t, PriorityMedium, got.Priority, "default priority")
	assert.Nil(t, got.DueDate)
	assert.False(t, got.CreatedAt.Before(before))
	assert.False(t, got.CreatedAt.After(after))
	assert.Equal(t, got.CreatedAt, got.UpdatedAt)
}

func TestStore_Add_TrimsTitle(t *testing.T) {
	t.Parallel()
	s := newTestStore(t)
	got := s.Add(NewTodo{Title: "  buy milk  "})
	assert.Equal(t, "buy milk", got.Title)
}

func TestStore_Add_RespectsPriorityWhenProvided(t *testing.T) {
	t.Parallel()
	s := newTestStore(t)
	got := s.Add(NewTodo{Title: "x", Priority: PriorityHigh})
	assert.Equal(t, PriorityHigh, got.Priority)
}

func TestStore_Get_ReturnsCopy(t *testing.T) {
	t.Parallel()
	s := newTestStore(t)
	added := s.Add(NewTodo{Title: "x"})
	got, ok := s.Get(added.ID)
	assert.True(t, ok)
	assert.Equal(t, added, got)

	got.Title = "mutated"
	again, _ := s.Get(added.ID)
	assert.Equal(t, "x", again.Title, "store copy must not be affected by external mutation")
}

func TestStore_Get_NotFound(t *testing.T) {
	t.Parallel()
	s := newTestStore(t)
	_, ok := s.Get("nope")
	assert.False(t, ok)
}

func TestStore_Update_PartialFields(t *testing.T) {
	t.Parallel()
	s := newTestStore(t)
	added := s.Add(NewTodo{Title: "x"})
	time.Sleep(time.Millisecond) // updatedAt 비교를 위한 1ms 차이

	completed := true
	newTitle := "y"
	got, err := s.Update(added.ID, Patch{Title: &newTitle, Completed: &completed})
	assert.NoError(t, err)
	assert.Equal(t, "y", got.Title)
	assert.True(t, got.Completed)
	assert.True(t, got.UpdatedAt.After(added.UpdatedAt))
	assert.Equal(t, added.CreatedAt, got.CreatedAt, "createdAt must not change")
}

func TestStore_Update_ClearDueDate(t *testing.T) {
	t.Parallel()
	s := newTestStore(t)
	due := time.Now().Add(time.Hour)
	added := s.Add(NewTodo{Title: "x", DueDate: &due})
	got, err := s.Update(added.ID, Patch{ClearDueDate: true})
	assert.NoError(t, err)
	assert.Nil(t, got.DueDate)
}

func TestStore_Update_NotFound(t *testing.T) {
	t.Parallel()
	s := newTestStore(t)
	completed := true
	_, err := s.Update("nope", Patch{Completed: &completed})
	assert.ErrorIs(t, err, ErrNotFound)
}

func TestStore_Update_ValidationError(t *testing.T) {
	t.Parallel()
	s := newTestStore(t)
	added := s.Add(NewTodo{Title: "x"})
	empty := ""
	_, err := s.Update(added.ID, Patch{Title: &empty})
	var verr *ValidationError
	assert.ErrorAs(t, err, &verr)
	assert.Equal(t, "title", verr.Field)
}

func TestStore_Delete_Success(t *testing.T) {
	t.Parallel()
	s := newTestStore(t)
	added := s.Add(NewTodo{Title: "x"})
	assert.NoError(t, s.Delete(added.ID))
	_, ok := s.Get(added.ID)
	assert.False(t, ok)
}

func TestStore_Delete_NotFound(t *testing.T) {
	t.Parallel()
	s := newTestStore(t)
	assert.ErrorIs(t, s.Delete("nope"), ErrNotFound)
}
