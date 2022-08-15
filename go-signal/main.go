package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"
)

/*
Go의 signal 알람은 os.Signal 값을 채널에 보내는 방식으로 동작함
*/
func main() {
	sigs := make(chan os.Signal, 1)
	done := make(chan bool, 1)

	//signal.Notify는 우리가 지정한 signal을 받을 수 있는 채널을 받고 등록해줍니다.
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

	//signal을 받으면 프로그램을 종료하기 위한 코드임
	go func() {
		sig := <-sigs
		fmt.Println()
		fmt.Println(sig)
		done <- true
	}()

	fmt.Println("awaiting signal")
	<-done
	fmt.Println("exiting")
}
