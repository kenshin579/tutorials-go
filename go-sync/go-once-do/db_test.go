package go_once_do

import (
	"bytes"
	"log"
	"os"
	"testing"
)

func TestGetConnectionSimple(t *testing.T) {
	for i := 0; i < 10; i++ {
		go GetConnectionSimple()
	}

	//여러번 실행될 수 있음
	//todo : unit test에서 한번만 실행되고 있다는 것을 체크하려면?
}

func TestGetConnectionWithLock(t *testing.T) {
	for i := 0; i < 10; i++ {
		go GetConnectionWithLock()
	}
	//한번만 샐행된다
}

func TestGetConnectionWithLockV2(t *testing.T) {
	for i := 0; i < 10; i++ {
		go GetConnectionWithLockV2()
	}
	//한번만 샐행된다
}

func TestGetConnectionWithOnceDo(t *testing.T) {
	for i := 0; i < 10; i++ {
		go GetConnectionWithOnceDo()
	}
	//한번만 샐행된다
}

func captureOutput(f func() *DbConnection) string {
	var buf bytes.Buffer
	log.SetOutput(&buf)
	f()
	log.SetOutput(os.Stderr)
	return buf.String()
}
