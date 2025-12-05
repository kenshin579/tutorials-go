package once_do

import (
	"fmt"
	"sync"
	"testing"
)

func TestOnce_여러번_실행해도_단_한번만_실행한다(t *testing.T) {
	var once sync.Once
	onceBodyFunc := func() {
		fmt.Printf("  Body : Only once\n")
	}
	done := make(chan bool)
	for i := 0; i < 10; i++ {
		go func(num int) {
			fmt.Printf("[%d] go 실행\n", num)
			once.Do(onceBodyFunc)
			done <- true
		}(i)
	}
	for i := 0; i < 10; i++ {
		fmt.Printf("[%d] done\n", i)
		<-done
	}
}
