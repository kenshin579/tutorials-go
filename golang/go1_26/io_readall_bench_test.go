package go1_26_test

import (
	"bytes"
	"io"
	"testing"
)

// Go 1.26에서 io.ReadAll 성능이 2배 향상, 메모리 사용 50% 감소
// 벤치마크로 확인

func BenchmarkReadAll_1KB(b *testing.B) {
	data := bytes.Repeat([]byte("x"), 1024)
	for b.Loop() {
		io.ReadAll(bytes.NewReader(data))
	}
}

func BenchmarkReadAll_1MB(b *testing.B) {
	data := bytes.Repeat([]byte("x"), 1024*1024)
	for b.Loop() {
		io.ReadAll(bytes.NewReader(data))
	}
}

func BenchmarkReadAll_10MB(b *testing.B) {
	data := bytes.Repeat([]byte("x"), 10*1024*1024)
	for b.Loop() {
		io.ReadAll(bytes.NewReader(data))
	}
}

func TestReadAllBasic(t *testing.T) {
	data := []byte("Go 1.26 io.ReadAll performance improvement")
	result, err := io.ReadAll(bytes.NewReader(data))
	if err != nil {
		t.Fatal(err)
	}
	if string(result) != string(data) {
		t.Errorf("expected %q, got %q", data, result)
	}
}
