package go_timer

import (
	"fmt"
	"testing"
	"time"
)

//https://mingrammer.com/gobyexample/timers/
func TestTicker_Timer(t *testing.T) {
	//대기 시간만큼 타이머에게 지정해줌
	timer1 := time.NewTimer(time.Second * 2)
	<-timer1.C //blocked (timer1.C는 타이머가 만료되었음을 알려주는 값을 보내기전까지 타이머의 C 채널을 블로킹시킨다)
	fmt.Println("Timer1 expired")

	timer2 := time.NewTimer(time.Second)
	go func() {
		<-timer2.C
		fmt.Println("Timer2 expired")
	}()

	stop2 := timer2.Stop() //타미어가 만료되기전에 취소 시킬 수 있음
	if stop2 {
		fmt.Println("Timer2 stopped")
	}
}
