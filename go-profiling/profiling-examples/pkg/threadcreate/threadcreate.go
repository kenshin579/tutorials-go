package threadcreate

func CreateGoroutine1000() {
	for i := 0; i < 100000; i++ {
		go innerFunc()
	}
}

func innerFunc() {
	n := 0
	for i := 0; i < 1000000; i++ {
		n++
	}
}
