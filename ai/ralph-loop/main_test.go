package main

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestHandleHello_정상응답(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/api/hello", nil)
	w := httptest.NewRecorder()

	handleHello(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, "application/json", w.Header().Get("Content-Type"))

	var resp Response
	err := json.NewDecoder(w.Body).Decode(&resp)
	assert.NoError(t, err)
	assert.Equal(t, "ok", resp.Status)
	assert.Equal(t, "Hello from Ralph Loop demo!", resp.Message)
	assert.NotEmpty(t, resp.Time)
}

func TestHandleHealth_정상응답(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/health", nil)
	w := httptest.NewRecorder()

	handleHealth(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var resp HealthResponse
	err := json.NewDecoder(w.Body).Decode(&resp)
	assert.NoError(t, err)
	assert.Equal(t, "healthy", resp.Status)
	assert.NotEmpty(t, resp.Uptime)
}

func TestNewServer_라우터설정(t *testing.T) {
	srv := NewServer("8080")
	assert.NotNil(t, srv.Router)
	assert.Equal(t, "8080", srv.Port)
}
