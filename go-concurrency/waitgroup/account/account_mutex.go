package account

import (
	"sync"
)

type AccountMutex struct {
	Mutex           sync.Mutex
	customerBalance map[string]int
}

func (a *AccountMutex) add(customerName string) {
	a.Mutex.Lock()
	a.customerBalance[customerName]++
	a.Mutex.Unlock()
}

func (a *AccountMutex) GetBalance(customerName string) int {
	return a.customerBalance[customerName]
}
