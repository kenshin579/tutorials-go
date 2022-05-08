package same_channel

import (
	"fmt"
	"testing"
	"time"
)

type Result struct {
	ErrorName          string
	NumberOfOccurances int64
}

type ResultError struct {
	res Result
	err error
}

/*
여러 channel에서 기다리기보다 한 channel을 기다린다
*/
func Test_Same_Channel_Both_Result_And_Error(t *testing.T) {
	fmt.Println("Using same channels for error and result")
	errorIds := []string{
		"1001",
		"2001",
		"3001",
	}
	for _, e := range errorIds {
		r := getError(e)
		if r.err != nil {
			fmt.Printf("Failed: %s\n", r.err.Error())
			continue
		}
		fmt.Printf("Name: \"%s\" has occurred \"%d\" times\n", r.res.ErrorName, r.res.NumberOfOccurances)
	}
}

func getError(errorId string) (r ResultError) {
	errors := map[string]Result{
		"1001": {"a is undefined", 245},
		"2001": {"Cannot read property 'data' of undefined", 10352},
	}
	outputChannel := make(chan ResultError)
	go func() {
		time.Sleep(time.Second)
		if r, ok := errors[errorId]; ok {
			outputChannel <- ResultError{res: r, err: nil}
		} else {
			outputChannel <- ResultError{res: Result{}, err: fmt.Errorf("getErrorName: %s errorId not found", errorId)}
		}
	}()

	return <-outputChannel
}
