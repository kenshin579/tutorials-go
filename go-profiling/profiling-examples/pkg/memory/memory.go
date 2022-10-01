package memory

import (
	"time"
)

func AllocMemory() {
	bytes1000 := alloc1000()
	bytes1000[0] = '0'

	for {
		time.Sleep(1 * time.Second)
	}

	//select {}
}

func alloc1000() []byte {
	return make([]byte, 1000)
}
