package main

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

// https://quii.gitbook.io/learn-go-with-tests/build-an-application/http-server
func TestHelloWorldHandler(t *testing.T) {
	request, _ := http.NewRequest(http.MethodGet, "/welcome", nil)
	response := httptest.NewRecorder()

	HelloWorldHandler(response, request)

	assert.Equal(t, "Hello, world", response.Body.String())
	assert.Equal(t, http.StatusOK, response.Code)
}
