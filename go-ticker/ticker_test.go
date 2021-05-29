package main

import (
	"fmt"
	"testing"
	"time"
)

//time.Ticker는 GC로 처리가 안되기 때문에 time.NewTimer 사용하고 필요 없을 경우 Stop() 호출하자
func Test_Repeat_timeTick으로_반복작업(t *testing.T) {
	//time.Tick은 clock tick
	for now := range time.Tick(time.Second * 2) { //2초마다 반복
		fmt.Println(now, "done")
	}
}

//https://mingrammer.com/gobyexample/tickers/
func TestRepeat_NewTicker로_반복작업(t *testing.T) {
	//NewTicker()이용해서 0.5초마다 메시지를 전송하는 채널을 만듬
	ticker := time.NewTicker(time.Millisecond * 500)
	go func() {
		//이 goroutine에서는 0.5초 간격으로 채널로부터 메시지를 받음
		for t := range ticker.C { //ticker 채널로 tick 값을 전달 받음
			fmt.Println("Tick at", t)
		}
	}()

	time.Sleep(time.Millisecond * 1600)
	ticker.Stop() //티커가 멈추면 ticker 채널로 더 이상 값을 받지 못함
	fmt.Println("Ticker stopped")
}

func TestTicker_Select(t *testing.T) {
	ticker := time.NewTicker(500 * time.Millisecond)
	done := make(chan bool)

	go func() {
		for {
			select {
			case <-done:
				fmt.Println("Ticker done")
				return
			case t := <-ticker.C:
				fmt.Println("Tick at", t)
			}
		}
	}()

	time.Sleep(1600 * time.Millisecond)
	ticker.Stop()
	done <- true //여기서 ticker stop시 추가 작업을 할 수 있음
	fmt.Println("Ticker stopped")
}
