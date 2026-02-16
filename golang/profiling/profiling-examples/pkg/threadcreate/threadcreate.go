// threadcreate 패키지는 대량 고루틴 생성을 통한 OS 스레드 생성 프로파일링 예제를 제공한다.
// Go 런타임은 고루틴이 시스템 콜 등으로 블로킹되면 새로운 OS 스레드를 생성한다.
// 대량의 고루틴을 동시에 실행하면 스레드 생성 패턴을 pprof threadcreate 프로파일에서 확인할 수 있다.
package threadcreate

// CreateGoroutine1000은 100,000개의 고루틴을 생성하여 대량 동시 실행을 시뮬레이션한다.
// 각 고루틴은 CPU 연산을 수행하며, Go 스케줄러가 이를 OS 스레드에 분배한다.
// 고루틴 수가 GOMAXPROCS보다 훨씬 많으므로 스케줄링 오버헤드가 발생한다.
func CreateGoroutine1000() {
	for i := 0; i < 100000; i++ {
		go innerFunc()
	}
}

// innerFunc는 단순 연산을 수행하는 고루틴의 작업 단위이다.
func innerFunc() {
	n := 0
	for i := 0; i < 1000000; i++ {
		n++
	}
}
