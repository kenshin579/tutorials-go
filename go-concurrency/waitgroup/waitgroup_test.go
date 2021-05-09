package waitgroup

import (
	"fmt"
	"sync"
	"testing"
	"time"
)

//https://tutorialedge.net/golang/go-waitgroup-tutorial/
func TestGoroutine(t *testing.T) {
	fmt.Println("Hello World")
	go myFunc()
	fmt.Println("Finished Execution")
}

func myFunc() {
	time.Sleep(time.Second)
	fmt.Println("Inside my goroutine") //호출하는 메서드가 먼저 종료되어서 출력이 안될 수 있음
}

func TestWaitGroup_Simple(t *testing.T) {
	fmt.Println("Hello World")

	var wg sync.WaitGroup
	wg.Add(1)
	go myFuncSimple(&wg)
	wg.Wait()

	fmt.Println("Finished Execution")
}

func myFuncSimple(wg *sync.WaitGroup) {
	fmt.Println("Inside my goroutine")
	wg.Done()
}

func TestWaitGroup_Anonymous_Func(t *testing.T) {
	fmt.Println("Hello World")

	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		fmt.Println("Inside my goroutine")
		wg.Done()
	}()
	wg.Wait()

	fmt.Println("Finished Execution")
}

func TestWaitGroup(t *testing.T) {
	var wg sync.WaitGroup

	for i := 1; i <= 5; i++ {
		wg.Add(1) //1만큼씩 wg counter를 증가시킨다
		go worker(i, &wg)

	}

	//모든 goroutine이 끝날 때까직 기다린다 - wg counter가 0이 될때까지 기다린다
	wg.Wait()
}

func worker(id int, wg *sync.WaitGroup) {
	defer wg.Done() //wg counter 1만큼 줄인다
	fmt.Printf("Woker %d starting\n", id)
	time.Sleep(time.Second)

	fmt.Printf("worker %d done\n", id)
}
