package main

import (
	"fmt"
	"time"
)

//https://mingrammer.com/gobyexample/tickers/
func main() {
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
