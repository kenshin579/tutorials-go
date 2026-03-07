package backend

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func newTestStore(t *testing.T) *TodoStore {
	t.Helper()
	tmpFile := filepath.Join(t.TempDir(), "todos.json")
	return NewTodoStore(tmpFile)
}

func TestNewTodoStore_FileNotExist(t *testing.T) {
	store := newTestStore(t)

	assert.Empty(t, store.Todos)
}

func TestNewTodoStore_LoadExistingFile(t *testing.T) {
	tmpFile := filepath.Join(t.TempDir(), "todos.json")
	data := `[{"id":"1","title":"test","done":false,"createdAt":"2024-01-01T00:00:00Z"}]`
	require.NoError(t, os.WriteFile(tmpFile, []byte(data), 0644))

	store := NewTodoStore(tmpFile)

	assert.Len(t, store.Todos, 1)
	assert.Equal(t, "test", store.Todos[0].Title)
}

func TestTodoStore_SaveAndLoad(t *testing.T) {
	store := newTestStore(t)
	store.Todos = []Todo{
		{ID: "1", Title: "할 일 1", Done: false},
		{ID: "2", Title: "할 일 2", Done: true},
	}

	require.NoError(t, store.Save())

	// 같은 파일로 새 store 생성하여 Load 검증
	loaded := NewTodoStore(store.filePath)
	assert.Len(t, loaded.Todos, 2)
	assert.Equal(t, "할 일 1", loaded.Todos[0].Title)
	assert.True(t, loaded.Todos[1].Done)
}

func TestTodoStore_LoadInvalidJSON(t *testing.T) {
	tmpFile := filepath.Join(t.TempDir(), "todos.json")
	require.NoError(t, os.WriteFile(tmpFile, []byte("invalid json"), 0644))

	store := NewTodoStore(tmpFile)

	assert.Empty(t, store.Todos)
}

func TestTodoStore_ExportAndImport(t *testing.T) {
	// Export
	store := newTestStore(t)
	store.Todos = []Todo{
		{ID: "1", Title: "export test", Done: false},
	}
	require.NoError(t, store.Save())

	exportPath := filepath.Join(t.TempDir(), "exported.json")
	require.NoError(t, store.ExportToFile(exportPath))

	// Import into new store
	newStore := newTestStore(t)
	require.NoError(t, newStore.LoadFromFile(exportPath))

	assert.Len(t, newStore.Todos, 1)
	assert.Equal(t, "export test", newStore.Todos[0].Title)
}

func TestTodoStore_LoadFromFile_NotExist(t *testing.T) {
	store := newTestStore(t)

	err := store.LoadFromFile("/nonexistent/path/file.json")

	assert.Error(t, err)
}
