package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

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

func TestCpuWork(t *testing.T) {
	// panic 없이 실행되는지 확인
	assert.NotPanics(t, func() {
		cpuWork()
	})
}

func TestMemoryWork(t *testing.T) {
	assert.NotPanics(t, func() {
		memoryWork()
	})
}

func TestMutexWork(t *testing.T) {
	assert.NotPanics(t, func() {
		mutexWork()
	})
}

func TestHostname(t *testing.T) {
	h := hostname()
	assert.NotEmpty(t, h)
	assert.NotEqual(t, "unknown", h)
}

func TestGetEnv(t *testing.T) {
	assert.Equal(t, "fallback", getEnv("NON_EXISTENT_KEY_12345", "fallback"))

	t.Setenv("TEST_PYROSCOPE_KEY", "value123")
	assert.Equal(t, "value123", getEnv("TEST_PYROSCOPE_KEY", "fallback"))
}
