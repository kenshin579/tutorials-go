package account

import (
	"fmt"

	"github.com/go-redsync/redsync/v4"
)

type AccountRedSync struct {
	Mutex           *redsync.Mutex
	CustomerBalance map[string]int
}

func (a *AccountRedSync) add(customerName string) {
	if err := a.Mutex.Lock(); err != nil {
		fmt.Println(err)
	}

	a.CustomerBalance[customerName]++

	if ok, err := a.Mutex.Unlock(); !ok || err != nil {
		fmt.Println(err)
	}
}

func (a *AccountRedSync) GetBalance(customerName string) int {
	return a.CustomerBalance[customerName]
}
