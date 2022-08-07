package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMyFunc(t *testing.T) {
	resp, err := MyFunc()
	assert.Error(t, err)
	assert.Equal(t, Response{
		Message: "Failure",
	}, resp)
}
