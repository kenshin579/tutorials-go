package context_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

// 타입 안전한 context key 정의
type contextKey string

const (
	userIDKey    contextKey = "userID"
	requestIDKey contextKey = "requestID"
)

// TestWithValue - context.WithValue 기본 사용법
func TestWithValue(t *testing.T) {
	ctx := context.Background()
	ctx = context.WithValue(ctx, userIDKey, "user-123")
	ctx = context.WithValue(ctx, requestIDKey, "req-456")

	// 값 조회
	userID := ctx.Value(userIDKey).(string)
	requestID := ctx.Value(requestIDKey).(string)

	assert.Equal(t, "user-123", userID)
	assert.Equal(t, "req-456", requestID)
}

// TestWithValueNotFound - 존재하지 않는 키 조회
func TestWithValueNotFound(t *testing.T) {
	ctx := context.Background()

	val := ctx.Value(userIDKey)
	assert.Nil(t, val, "존재하지 않는 키는 nil 반환")
}

// TestWithValueChain - parent context의 값도 child에서 접근 가능
func TestWithValueChain(t *testing.T) {
	parent := context.WithValue(context.Background(), userIDKey, "parent-user")
	child := context.WithValue(parent, requestIDKey, "child-req")

	// child에서 parent의 값 접근 가능
	assert.Equal(t, "parent-user", child.Value(userIDKey))
	assert.Equal(t, "child-req", child.Value(requestIDKey))

	// parent에서 child의 값은 접근 불가
	assert.Nil(t, parent.Value(requestIDKey))
}
