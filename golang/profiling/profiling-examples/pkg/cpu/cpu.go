// cpu 패키지는 CPU 프로파일링을 위한 부하 생성 함수를 제공한다.
// pprof CPU 프로파일에서 increase1000과 increase2000의 실행 시간 비율(약 1:2)을 확인할 수 있다.
package cpu

// IncreaseInt는 무한 루프에서 CPU 연산을 수행하여 CPU 프로파일에 잡히도록 한다.
// pprof top 명령어에서 이 함수의 cum(누적) 시간이 높게 나타난다.
func IncreaseInt() {
	i := 0
	for {
		i = increase1000(i)
		i = increase2000(i)
	}
}

// IncreaseIntGoroutine은 별도 고루틴에서 CPU 부하를 생성한다.
// pprof에서 IncreaseInt와 별개의 고루틴으로 추적된다.
func IncreaseIntGoroutine() {
	go func() {
		i := 0
		for {
			i = increase1000(i)
			i = increase2000(i)
		}
	}()
}

// increase1000은 1000회 루프를 돌며 CPU 시간을 소비한다.
// increase2000 대비 절반의 CPU 시간을 차지하여, 프로파일에서 비율 차이를 확인할 수 있다.
func increase1000(n int) int {
	for n := 0; n < 1000; n++ {
		n = n + 1
	}
	return n
}

// increase2000은 2000회 루프를 돌며 CPU 시간을 소비한다.
// increase1000 대비 약 2배의 CPU 시간을 차지한다.
func increase2000(n int) int {
	for n := 0; n < 2000; n++ {
		n = n + 1
	}
	return n
}
