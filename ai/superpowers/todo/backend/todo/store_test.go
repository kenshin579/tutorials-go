package todo

import (
	"sync"
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

func TestStore_List_FilterByStatus(t *testing.T) {
	t.Parallel()
	s := newTestStore(t)
	a := s.Add(NewTodo{Title: "active"})
	b := s.Add(NewTodo{Title: "done"})
	completed := true
	if _, err := s.Update(b.ID, Patch{Completed: &completed}); err != nil {
		t.Fatalf("update: %v", err)
	}
	_ = a

	tests := []struct {
		status    StatusFilter
		wantLen   int
		wantTitle string
	}{
		{StatusAll, 2, ""},
		{StatusActive, 1, "active"},
		{StatusCompleted, 1, "done"},
	}
	for _, tc := range tests {
		tc := tc
		t.Run(string(tc.status), func(t *testing.T) {
			t.Parallel()
			got := s.List(Query{Status: tc.status})
			assert.Len(t, got, tc.wantLen)
			if tc.wantTitle != "" && len(got) == 1 {
				assert.Equal(t, tc.wantTitle, got[0].Title)
			}
		})
	}
}

func TestStore_List_SortByPriority(t *testing.T) {
	t.Parallel()
	s := newTestStore(t)
	low := s.Add(NewTodo{Title: "low", Priority: PriorityLow})
	high := s.Add(NewTodo{Title: "high", Priority: PriorityHigh})
	mid := s.Add(NewTodo{Title: "mid", Priority: PriorityMedium})

	asc := s.List(Query{Sort: SortPriority, Order: OrderAsc})
	assert.Equal(t, []string{low.ID, mid.ID, high.ID}, ids(asc))

	desc := s.List(Query{Sort: SortPriority, Order: OrderDesc})
	assert.Equal(t, []string{high.ID, mid.ID, low.ID}, ids(desc))
}

func TestStore_List_SortByDueDate_NilLast(t *testing.T) {
	t.Parallel()
	s := newTestStore(t)
	t1 := time.Date(2026, 5, 1, 0, 0, 0, 0, time.UTC)
	t2 := time.Date(2026, 6, 1, 0, 0, 0, 0, time.UTC)
	a := s.Add(NewTodo{Title: "a", DueDate: &t2})
	b := s.Add(NewTodo{Title: "b", DueDate: nil})
	c := s.Add(NewTodo{Title: "c", DueDate: &t1})

	asc := s.List(Query{Sort: SortDueDate, Order: OrderAsc})
	assert.Equal(t, []string{c.ID, a.ID, b.ID}, ids(asc), "nil dueDate must be last")

	desc := s.List(Query{Sort: SortDueDate, Order: OrderDesc})
	assert.Equal(t, []string{a.ID, c.ID, b.ID}, ids(desc), "nil dueDate must be last regardless of order")
}

func TestStore_List_SortByCreatedAt_DescDefault(t *testing.T) {
	t.Parallel()
	s := newTestStore(t)
	first := s.Add(NewTodo{Title: "first"})
	time.Sleep(time.Millisecond)
	second := s.Add(NewTodo{Title: "second"})

	desc := s.List(Query{}) // default sort=createdAt, order=desc
	assert.Equal(t, []string{second.ID, first.ID}, ids(desc))
}

func TestStore_List_TieBreakByCreatedAtAsc(t *testing.T) {
	t.Parallel()
	s := newTestStore(t)
	a := s.Add(NewTodo{Title: "a", Priority: PriorityHigh})
	time.Sleep(time.Millisecond)
	b := s.Add(NewTodo{Title: "b", Priority: PriorityHigh})

	asc := s.List(Query{Sort: SortPriority, Order: OrderAsc})
	assert.Equal(t, []string{a.ID, b.ID}, ids(asc))
}

func TestStore_ConcurrentAddList(t *testing.T) {
	t.Parallel()
	s := newTestStore(t)
	var wg sync.WaitGroup
	for i := 0; i < 100; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			s.Add(NewTodo{Title: "t"})
		}(i)
	}
	done := make(chan struct{})
	go func() {
		for {
			select {
			case <-done:
				return
			default:
				_ = s.List(Query{})
			}
		}
	}()
	wg.Wait()
	close(done)
	assert.Len(t, s.List(Query{}), 100)
}

func ids(ts []Todo) []string {
	out := make([]string, len(ts))
	for i, t := range ts {
		out[i] = t.ID
	}
	return out
}
