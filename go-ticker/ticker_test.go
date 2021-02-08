package main

import (
	"fmt"
	"testing"
	"time"
)

//https://mingrammer.com/gobyexample/tickers/
func TestTicker1(t *testing.T) {
	ticker := time.NewTicker(time.Millisecond * 500)
	go func() {
		for t := range ticker.C { //ticker 채널로 tick 값을 전달 받음
			fmt.Println("Tick at", t)
		}
	}()

	time.Sleep(time.Millisecond * 1600)
	ticker.Stop() //티커가 멈추면 ticker 채널로 더 이상 값을 받지 못함
	fmt.Println("Ticker stopped")
}

func TestTicker2(t *testing.T) {
	ticker := time.NewTicker(500 * time.Millisecond)
	done := make(chan bool)

	go func() {
		for {
			select {
			case <-done:
				return
			case t := <-ticker.C:
				fmt.Println("Tick at", t)
			}
		}
	}()

	time.Sleep(1600 * time.Millisecond)
	ticker.Stop()
	done <- true
	fmt.Println("Ticker stopped")
}
