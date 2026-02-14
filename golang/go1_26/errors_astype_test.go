package go1_26_test

import (
	"errors"
	"fmt"
	"testing"
)

// AppError - 커스텀 에러 타입
type AppError struct {
	Code    int
	Message string
}

func (e *AppError) Error() string {
	return fmt.Sprintf("[%d] %s", e.Code, e.Message)
}

// ValidationError - 또 다른 커스텀 에러 타입
type ValidationError struct {
	Field   string
	Message string
}

func (e *ValidationError) Error() string {
	return fmt.Sprintf("validation error on %s: %s", e.Field, e.Message)
}

func TestErrorsAs_OldWay(t *testing.T) {
	err := fmt.Errorf("wrapped: %w", &AppError{Code: 404, Message: "not found"})

	// 기존 방식: 타겟 변수를 미리 선언해야 함
	var target *AppError
	if errors.As(err, &target) {
		fmt.Printf("기존 방식 - Code: %d, Message: %s\n", target.Code, target.Message)
	}

	if target.Code != 404 {
		t.Errorf("expected 404, got %d", target.Code)
	}
}

func TestErrorsAsType_NewWay(t *testing.T) {
	err := fmt.Errorf("wrapped: %w", &AppError{Code: 404, Message: "not found"})

	// 새로운 방식: errors.AsType[T]() - 타입 안전, 변수 선언 불필요
	if target, ok := errors.AsType[*AppError](err); ok {
		fmt.Printf("새로운 방식 - Code: %d, Message: %s\n", target.Code, target.Message)
		if target.Code != 404 {
			t.Errorf("expected 404, got %d", target.Code)
		}
	} else {
		t.Error("expected to find AppError")
	}
}

func TestErrorsAsType_NotFound(t *testing.T) {
	err := fmt.Errorf("wrapped: %w", &AppError{Code: 500, Message: "server error"})

	// ValidationError 타입으로 찾기 시도 → 실패
	if _, ok := errors.AsType[*ValidationError](err); ok {
		t.Error("should not find ValidationError")
	} else {
		fmt.Println("ValidationError not found (expected)")
	}
}

func TestErrorsAsType_ChainedErrors(t *testing.T) {
	innerErr := &ValidationError{Field: "email", Message: "invalid format"}
	wrappedErr := fmt.Errorf("outer: %w", fmt.Errorf("inner: %w", innerErr))

	// 중첩된 에러 체인에서도 동작
	if target, ok := errors.AsType[*ValidationError](wrappedErr); ok {
		fmt.Printf("Found in chain - Field: %s, Message: %s\n", target.Field, target.Message)
		if target.Field != "email" {
			t.Errorf("expected 'email', got '%s'", target.Field)
		}
	} else {
		t.Error("expected to find ValidationError in chain")
	}
}
