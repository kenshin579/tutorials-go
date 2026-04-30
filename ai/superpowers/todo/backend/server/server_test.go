package server_test

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/kenshin579/tutorials-go/ai/superpowers/todo/backend/server"
	"github.com/kenshin579/tutorials-go/ai/superpowers/todo/backend/todo"
)

func setupServer(t *testing.T) *httptest.Server {
	t.Helper()
	s := todo.NewStore()
	e := server.New(s)
	ts := httptest.NewServer(e)
	t.Cleanup(ts.Close)
	return ts
}

func TestServer_Health(t *testing.T) {
	t.Parallel()
	ts := setupServer(t)

	resp, err := http.Get(ts.URL + "/api/health")
	assert.NoError(t, err)
	defer resp.Body.Close()
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	var body map[string]string
	assert.NoError(t, json.NewDecoder(resp.Body).Decode(&body))
	assert.Equal(t, "ok", body["status"])
}

func TestServer_CRUDLifecycle(t *testing.T) {
	t.Parallel()
	ts := setupServer(t)

	// 1. Create
	create := bytes.NewBufferString(`{"title":"buy milk","priority":"high"}`)
	resp, err := http.Post(ts.URL+"/api/todos", "application/json", create)
	assert.NoError(t, err)
	defer resp.Body.Close()
	assert.Equal(t, http.StatusCreated, resp.StatusCode)

	var created todo.Todo
	assert.NoError(t, json.NewDecoder(resp.Body).Decode(&created))
	assert.NotEmpty(t, created.ID)
	assert.Equal(t, "buy milk", created.Title)

	// 2. List
	resp, err = http.Get(ts.URL + "/api/todos")
	assert.NoError(t, err)
	defer resp.Body.Close()
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	var list []todo.Todo
	assert.NoError(t, json.NewDecoder(resp.Body).Decode(&list))
	assert.Len(t, list, 1)

	// 3. Update (toggle complete)
	patch := bytes.NewBufferString(`{"completed":true}`)
	req, _ := http.NewRequest(http.MethodPatch, ts.URL+"/api/todos/"+created.ID, patch)
	req.Header.Set("Content-Type", "application/json")
	resp, err = http.DefaultClient.Do(req)
	assert.NoError(t, err)
	defer resp.Body.Close()
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	var updated todo.Todo
	assert.NoError(t, json.NewDecoder(resp.Body).Decode(&updated))
	assert.True(t, updated.Completed)

	// 4. Delete
	req, _ = http.NewRequest(http.MethodDelete, ts.URL+"/api/todos/"+created.ID, nil)
	resp, err = http.DefaultClient.Do(req)
	assert.NoError(t, err)
	defer resp.Body.Close()
	assert.Equal(t, http.StatusNoContent, resp.StatusCode)

	// 5. List again — should be empty
	resp, err = http.Get(ts.URL + "/api/todos")
	assert.NoError(t, err)
	defer resp.Body.Close()
	body, _ := io.ReadAll(resp.Body)
	assert.Equal(t, "[]\n", string(body))
}

func TestServer_CORSPreflight(t *testing.T) {
	t.Parallel()
	ts := setupServer(t)

	req, _ := http.NewRequest(http.MethodOptions, ts.URL+"/api/todos", nil)
	req.Header.Set("Origin", "http://localhost:5173")
	req.Header.Set("Access-Control-Request-Method", "POST")
	resp, err := http.DefaultClient.Do(req)
	assert.NoError(t, err)
	defer resp.Body.Close()
	assert.Contains(t, resp.Header.Get("Access-Control-Allow-Origin"), "http://localhost:5173")
}
