package main

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestHelloHandler(t *testing.T) {
	request, _ := http.NewRequest(http.MethodGet, "/HelloHandler", nil)
	response := httptest.NewRecorder()

	HelloHandler(response, request)

	assert.Equal(t, "hello world!", response.Body.String())
	assert.Equal(t, http.StatusOK, response.Code)
}
