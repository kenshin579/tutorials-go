package go_context

import (
	"context"
	"fmt"
	"testing"
	"time"
)

func Test_Context_WithCancel(t *testing.T) {
	ctx := context.Background()
	cancelCtx, cancelFunc := context.WithCancel(ctx) //취소신호를 수신할 수 있는 컨텍스트를 생성한다
	go task(cancelCtx)
	time.Sleep(time.Second * 3)
	cancelFunc()
	time.Sleep(time.Second * 1)
}

func Test_Context_WithTimeout(t *testing.T) {
	ctx := context.Background()
	cancelCtx, cancel := context.WithTimeout(ctx, time.Second*3) //컨텍스트가 특정 시간이후에 자동으로 취소되도록 생성함
	defer cancel()
	go task(cancelCtx)
	time.Sleep(time.Second * 4)
}

func Test_Context_WithDeadline(t *testing.T) {
	ctx := context.Background()
	cancelCtx, cancel := context.WithDeadline(ctx, time.Now().Add(time.Second*5)) // 현재+5분이후에 자동으로 취소되도록 생성함
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
			fmt.Println("error message:", ctx.Err())
			return
		default:
			fmt.Println(i)
			time.Sleep(time.Second * 1)
			i++
		}
	}
}
