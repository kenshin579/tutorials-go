package channel

import "testing"

// BenchmarkUnbuffered - unbuffered channel 성능
func BenchmarkUnbuffered(b *testing.B) {
	ch := make(chan int)

	go func() {
		for range b.N {
			ch <- 1
		}
	}()

	for range b.N {
		<-ch
	}
}

// BenchmarkBuffered1 - buffered channel (크기 1)
func BenchmarkBuffered1(b *testing.B) {
	ch := make(chan int, 1)

	go func() {
		for range b.N {
			ch <- 1
		}
	}()

	for range b.N {
		<-ch
	}
}

// BenchmarkBuffered100 - buffered channel (크기 100)
func BenchmarkBuffered100(b *testing.B) {
	ch := make(chan int, 100)

	go func() {
		for range b.N {
			ch <- 1
		}
	}()

	for range b.N {
		<-ch
	}
}

// BenchmarkBuffered1000 - buffered channel (크기 1000)
func BenchmarkBuffered1000(b *testing.B) {
	ch := make(chan int, 1000)

	go func() {
		for range b.N {
			ch <- 1
		}
	}()

	for range b.N {
		<-ch
	}
}
