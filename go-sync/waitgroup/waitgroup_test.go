package waitgroup

import (
	"fmt"
	"sync"
	"testing"
	"time"
)

func TestWaitGroup(t *testing.T) {
	var wg sync.WaitGroup

	for i := 1; i <= 5; i++ {
		wg.Add(1) //
		go worker(i, &wg)

	}

	//모든 goroutine이 끝날 때까직 기다린다 - wg counter가 0이 될때까지 기다린다
	wg.Wait()
}

func worker(id int, wg *sync.WaitGroup) {
	defer wg.Done()
	fmt.Printf("Woker %d starting\n", id)
	time.Sleep(time.Second)

	fmt.Printf("worker %d done\n", id)
}
