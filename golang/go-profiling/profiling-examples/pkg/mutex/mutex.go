package mutex

import (
	"fmt"
	"sync"
)

var mutex = sync.Mutex{}

func Mutex01() {
	for {
		mutex.Lock()
		fmt.Printf("Mutex01\n")
		mutex.Unlock()
	}
}

func Mutex02() {
	for {
		mutex.Lock()
		fmt.Printf("Mutex02\n")
		mutex.Unlock()
	}
}

func Mutex03() {
	for {
		mutex.Lock()
		fmt.Printf("Mutex03\n")
		mutex.Unlock()
	}
}
