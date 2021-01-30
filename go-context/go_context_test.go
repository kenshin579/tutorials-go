package go_context

import (
	"context"
	"fmt"
	"testing"
	"time"
)

func Test_Context_WithCancel(t *testing.T) {
	ctx := context.Background()
	cancelCtx, cancelFunc := context.WithCancel(ctx)
	go task(cancelCtx)
	time.Sleep(time.Second * 3)
	cancelFunc()
	time.Sleep(time.Second * 1)
}

func Test_Context_WithTimeout(t *testing.T) {
	ctx := context.Background()
	cancelCtx, cancel := context.WithTimeout(ctx, time.Second*3)
	defer cancel()
	go task(cancelCtx)
	time.Sleep(time.Second * 4)
}

func Test_Context_WithDeadline(t *testing.T) {
	ctx := context.Background()
	cancelCtx, cancel := context.WithDeadline(ctx, time.Now().Add(time.Second*5))
	defer cancel()
	go task(cancelCtx)
	time.Sleep(time.Second * 6)
}

func task(ctx context.Context) {
	i := 1
	for {
		select {
		case <-ctx.Done(): //cancelFunc() 호출시 실행된다
			fmt.Println("Gracefully exit")
			fmt.Println(ctx.Err())
			return
		default:
			fmt.Println(i)
			time.Sleep(time.Second * 1)
			i++
		}
	}
}
