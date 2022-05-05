package separate_channel

import (
	"fmt"
	"testing"
	"time"
)

func Test_Separate_Channel_Both_Result_And_Error(t *testing.T) {
	fmt.Println("Using separate channels for error and result")

	errorIds := []string{
		"1001",
		"2001",
		"3001",
	}

	for _, e := range errorIds {
		r, err := getError(e)
		if err != nil {
			fmt.Printf("Failed: %s\n", err.Error())
			continue
		}
		fmt.Printf("Name: \"%s\" has occurred \"%d\" times\n", r.ErrorName, r.NumberOfOccurances)
	}
}

func getError(errorId string) (r *Result, err error) {
	nameOut, nameErr := getErrorName(errorId)
	occurancesOut, occurancesErr := getOccurances(errorId)

	var ok bool

	if err, ok = <-nameErr; ok {
		return
	}
	if err, ok = <-occurancesErr; ok {
		return
	}

	r = &Result{
		ErrorName:          <-nameOut,
		NumberOfOccurances: <-occurancesOut,
	}

	//todo: err는 어떻게 반환되는 건가?
	return
}

func getErrorName(errorId string) (<-chan string, <-chan error) {
	names := map[string]string{
		"1001": "a is undefined",
		"2001": "Cannot read property 'data' of undefined",
	}

	out := make(chan string, 1)
	errs := make(chan error, 1)

	go func() {
		time.Sleep(time.Second)
		if name, ok := names[errorId]; ok {
			out <- name
		} else {
			errs <- fmt.Errorf("getErrorName: %s errorId not found", errorId)
		}

		close(out)
		close(errs)
	}()

	return out, errs
}

func getOccurances(errorId string) (<-chan int64, <-chan error) {
	occurances := map[string]int64{
		"1001": 245,
		"2001": 10352,
	}

	out := make(chan int64, 1)
	errs := make(chan error, 1)

	go func() {
		time.Sleep(time.Second)
		if occ, ok := occurances[errorId]; ok {
			out <- occ
		} else {
			errs <- fmt.Errorf("getOccurances: %s errorId not found", errorId)
		}

		close(out)
		close(errs)
	}()

	return out, errs
}
