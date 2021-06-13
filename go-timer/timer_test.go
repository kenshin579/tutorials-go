package go_timer

import (
	"fmt"
	"testing"
	"time"
)

//https://mingrammer.com/gobyexample/timers/
func TestTimeout_NewTimer_timerC로_expired될때까지_블록이_된다(t *testing.T) {
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

//https://yourbasic.org/golang/time-reset-wait-stop-timeout-cancel-interval/
func TestTimeout_TimeAfter_사용해서_타임아웃_예제(t *testing.T) {
	CNN := make(chan string)
	go func() { //goroutine에서 news article을 쓴다
		time.Sleep(time.Second * 1)
		fmt.Println("write to news channel : ")
		CNN <- "2021/5/29 - politic news!!!"
	}()

	select {
	case news := <-CNN:
		fmt.Println(news)
	case <-time.After(time.Second * 3): //waits for a specified duration and
		fmt.Println("No news")
	}
}

//time.Timer (?) 는 timer가 fires 되지 않으면 GC에 의해서 처리가 안되어서
func TestTimeout_NewTimer으로_Stop_시킬수_있다(t *testing.T) {
	CNN := make(chan string)
	go func() { //goroutine에서 news article을 쓴다
		time.Sleep(time.Second * 1)
		fmt.Println("write to news channel")
		CNN <- "2021/5/29 - politic news!!!"
	}()

	for alive := true; alive; {
		timer := time.NewTimer(time.Second * 3)
		select {
		case news := <-CNN:
			timer.Stop() //timer를 취소시킨다
			fmt.Println(news)
		case <-timer.C:
			alive = false
			fmt.Println("No news. Service aborting.")
		}
	}
}

func TestTimeout_AfterFunc(t *testing.T) {
	//AfterFunc으로 duration이 지난는 경우 실행할 수 함수를 인자로 넘겨준다
	timer := time.AfterFunc(time.Second*2, func() {
		fmt.Println("Foo run for two seconds")
	})

	defer timer.Stop() //취소시킬 수 있음

	// Do heavy work
	time.Sleep(time.Second * 1)
}

func command(maxWaitInSecond time.Duration) string {
	fmt.Println("start working")
	time.Sleep(time.Second * maxWaitInSecond)
	fmt.Println("finish working")
	return "done"
}

func TestCommand_실행후_Timeout되기전에_실행완료되는_경우(t *testing.T) {
	msg := make(chan string)
	go func() {
		msg <- command(1)
	}()

	for working := true; working; {
		timer := time.NewTimer(time.Second * 3)
		select {
		case <-msg:
			working = false
			timer.Stop()
			fmt.Println("stopping timer")
		case <-timer.C:
			working = false
			fmt.Println("command 완료")
		}
	}
}

func TestCommand_완료전에_timeout되는_경우(t *testing.T) {
	msg := make(chan string)
	go func() {
		msg <- command(4) //command 실행이 4초 이상되는 경우
	}()

	for working := true; working; {
		timer := time.NewTimer(time.Second * 3) //timeout 3초이후에 됨
		select {
		case <-msg:
			working = false
			timer.Stop()
			fmt.Println("stopping timer")
		case <-timer.C:
			working = false
			fmt.Println("command 완료")
		}
	}
}
