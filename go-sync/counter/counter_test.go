package counter

import (
	"sync"
	"testing"
)

func prepareCounter(key string, lim int64) *Counter {
	counter := Counter{}

	wg := sync.WaitGroup{}

	for i := int64(0); i < 100; i++ {
		icopy := i
		wg.Add(1)
		go func() {
			counter.Add(key, icopy)
			wg.Done()
		}()
	}

	wg.Wait()

	return &counter
}

func testCounter(t *testing.T) (testCount int64, counterCount int64) {
	counter := prepareCounter("key1", 100)

	testCount = int64(0)

	for i := int64(0); i < 100; i++ {
		testCount += i
	}

	var ok bool
	counterCount, ok = counter.Get("key1")

	if !ok {
		t.Error("error: expected counter.Get to succeed but failed")
	}
	return testCount, counterCount
}

func TestCounter(t *testing.T) {
	testCount, counterCount := testCounter(t)
	if testCount != counterCount {
		t.Error("error : possible race condition. Expected Counter to return ", testCount, " but got ", counterCount)
	}

}

func TestDeleteAndGetLastValue(t *testing.T) {
	counter := prepareCounter("key", 100)

	expectedResult, receiveSucceeded := counter.Get("key")

	if !receiveSucceeded {
		t.Error("counter.Get expected to succeed but failed")
	}

	starter := sync.WaitGroup{}
	wg := sync.WaitGroup{}

	var finalResult *int64
	muFinalResult := sync.Mutex{}

	starter.Add(1)
	for i := 0; i < 100; i++ {
		wg.Add(1)
		go func() {

			// ensures that all the goroutines start simultaneously.
			starter.Wait()

			if result, ok := counter.DeleteAndGetLastValue("key"); ok {
				muFinalResult.Lock()
				if finalResult != nil {
					t.Error("error: attempted to modify finalResult more than once")
				} else {
					finalResult = &result
				}
				muFinalResult.Unlock()
			}
			wg.Done()
		}()
	}

	// start all goroutines
	starter.Done()

	// wait for the goroutines to complete
	wg.Wait()

	if expectedResult != *finalResult {
		t.Error("error: unexpected final result. Expected ", expectedResult, " but got ", *finalResult)
	}
}

func TestCountableIncorrectGet(t *testing.T) {
	counter := Counter{}
	if _, ok := counter.Get("some key"); ok {
		t.Error("error: expected Get to fail because `some key` does not exists. ",
			"Instead it succeeded ")
	}
}
