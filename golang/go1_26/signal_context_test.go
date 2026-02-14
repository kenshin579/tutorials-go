package go1_26_test

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"testing"
)

func TestSignalNotifyContext(t *testing.T) {
	// signal.NotifyContext: 시그널을 컨텍스트 취소로 연결
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
	defer stop()

	// stop()을 호출하면 컨텍스트가 취소됨
	stop()

	// 컨텍스트가 취소되었는지 확인
	select {
	case <-ctx.Done():
		fmt.Println("Context cancelled")
		// Go 1.26: context.Cause로 원인 확인 가능
		cause := context.Cause(ctx)
		fmt.Printf("Cause: %v\n", cause)
		if cause == nil {
			t.Log("Cause is nil (stop() was called, not signal)")
		}
	default:
		t.Error("context should be cancelled after stop()")
	}
}

func TestContextCauseBasic(t *testing.T) {
	// context.Cause 기본 사용법
	ctx, cancel := context.WithCancelCause(context.Background())

	// 특정 원인으로 취소
	cancel(fmt.Errorf("user requested cancellation"))

	cause := context.Cause(ctx)
	fmt.Printf("Context cause: %v\n", cause)

	if cause == nil {
		t.Error("expected cause to be non-nil")
	}
}
