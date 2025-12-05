package main

import (
	"bytes"
	"fmt"
	"runtime"
	"strconv"
	"sync"
	"testing"
)

func getGID() uint64 {
	b := make([]byte, 64)
	b = b[:runtime.Stack(b, false)]
	b = bytes.TrimPrefix(b, []byte("goroutine "))
	b = b[:bytes.IndexByte(b, ' ')]
	n, _ := strconv.ParseUint(string(b), 10, 64)
	return n
}

func Test_Print_GID(t *testing.T) {
	fmt.Println("main", getGID())
	var wg sync.WaitGroup
	for i := 0; i < 10; i++ {
		i := i
		wg.Add(1)
		go func() {
			defer wg.Done()
			fmt.Println(i, getGID())
		}()
	}
	wg.Wait()
}
