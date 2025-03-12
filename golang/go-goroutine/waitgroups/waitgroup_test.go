package waitgroups

import (
	"fmt"
	"sync"
	"testing"
	"time"
)

/*
 */
func TestWait_Group(t *testing.T) {
	var wg sync.WaitGroup

	for i := 1; i <= 5; i++ {
		wg.Add(1) //counter가 delta 값을 더한다
		i := i
		go func() {
			defer wg.Done() //counter에서 1 값을 decrement한다
			worker(i)
		}()
	}

	wg.Wait() //counter가 0이 되면 끝난다
}

func worker(id int) {
	fmt.Printf("Worker %d starting\n", id)
	time.Sleep(time.Second)
	fmt.Printf("Worker %d done\n", id)
}

func Test_Execute(t *testing.T) {
	wg := new(sync.WaitGroup)
	wg.Add(2)

	// We are increasing the counter by 2
	// because we have 2 goroutines
	go runner1(wg)
	go runner2(wg)

	// This Blocks the execution
	// until its counter become 0
	wg.Wait()
}

func runner1(wg *sync.WaitGroup) {
	defer wg.Done() // This decreases counter by 1
	fmt.Print("\nI am first runner")

}

func runner2(wg *sync.WaitGroup) {
	defer wg.Done()
	fmt.Print("\nI am second runner")
}
