package go1_26_test

import (
	"fmt"
	"testing"
)

// Go 1.26에서 fmt.Errorf 포맷 없는 문자열 최적화: 0 할당 (기존 2 할당), 92% 빨라짐

func BenchmarkErrorf_NoFormat(b *testing.B) {
	// 포맷 인자 없는 경우 → Go 1.26에서 0 할당으로 최적화
	for b.Loop() {
		_ = fmt.Errorf("simple error message")
	}
}

func BenchmarkErrorf_WithFormat(b *testing.B) {
	// 포맷 인자 있는 경우
	for b.Loop() {
		_ = fmt.Errorf("error: %s at line %d", "parse failed", 42)
	}
}

func BenchmarkErrorf_WithWrap(b *testing.B) {
	inner := fmt.Errorf("inner error")
	for b.Loop() {
		_ = fmt.Errorf("outer: %w", inner)
	}
}

func TestErrorfNoFormat(t *testing.T) {
	err := fmt.Errorf("simple error")
	if err == nil {
		t.Error("expected non-nil error")
	}
	if err.Error() != "simple error" {
		t.Errorf("expected 'simple error', got '%s'", err.Error())
	}
}
