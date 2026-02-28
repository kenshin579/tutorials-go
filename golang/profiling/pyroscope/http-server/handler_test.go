package main

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
)

func TestHandleFast(t *testing.T) {
	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/fast", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	err := handleFast(c)

	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, rec.Code)
	assert.Contains(t, rec.Body.String(), "fast response")
}

func TestFibonacci(t *testing.T) {
	tests := []struct {
		n    int
		want int
	}{
		{0, 0},
		{1, 1},
		{10, 55},
		{20, 6765},
	}

	for _, tt := range tests {
		got := fibonacci(tt.n)
		assert.Equal(t, tt.want, got)
	}
}

func TestAllocateMemory(t *testing.T) {
	assert.NotPanics(t, func() {
		allocateMemory()
	})
}
