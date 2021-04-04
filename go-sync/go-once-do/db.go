package go_once_do

import (
	"fmt"
	"sync"
)

type DbConnection struct{}

var (
	dbConnOnce sync.Once
	mut        sync.Mutex
	conn       *DbConnection
)

func GetConnectionSimple() *DbConnection {
	//생성자 대신에 이미 생성한 값을 반환하도록 함
	//단점 : thread-safe하지 않은 코드임
	if conn == nil {
		fmt.Println("Creating DB Connection")
		conn = &DbConnection{}
	}
	return conn
}

func GetConnectionWithLock() *DbConnection {
	// Lock and unlock the entire GetInstance function
	mut.Lock()
	defer mut.Unlock()

	if conn == nil {
		fmt.Println("Creating DB Connection")
		conn = &DbConnection{}
	}
	return conn
}

func GetConnectionWithLockV2() *DbConnection {
	if conn == nil {
		mut.Lock()
		defer mut.Unlock()
		if conn == nil {
			fmt.Println("Creating DB Connection")
			conn = &DbConnection{}
		}
	}
	return conn
}

func GetConnectionWithOnceDo() *DbConnection {
	dbConnOnce.Do(func() {
		fmt.Println("Creating DB Connection")
		conn = &DbConnection{}
	})
	return conn
}
