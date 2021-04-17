package main

import (
	"fmt"
	"log"
	"regexp"
	"runtime"
	"testing"
	"time"
)

func TestTimestamp(t *testing.T) {
	now := time.Now()       //current local time
	timestamp := now.Unix() //number of seconds since Jan. 1, 1970 UTC
	timestampInSec := now.UnixNano()

	fmt.Println(now)            //time.Time (2021-04-17 14:21:16.254465 +0900 KST m=+0.000414223)
	fmt.Println(timestamp)      //1618636876
	fmt.Println(timestampInSec) //1618636876254465000
}

func expensiveCall() {
	time.Sleep(time.Second * 2)
}

func TestDuration1(t *testing.T) {
	start := time.Now()
	expensiveCall()
	end := time.Now()
	duration := time.Since(end)

	fmt.Printf("The call took %v to run.\n", end.Sub(start))
	fmt.Printf("The call took %v to run.\n", duration.Seconds())
}

func TestDuration_Seconds(t *testing.T) {
	start := time.Now()
	expensiveCall()
	end := time.Now()

	fmt.Printf("The call took %v to run.\n", end.Sub(start).Seconds())
}

func TestDuration_Measure_Execution_Time(t *testing.T) {
	defer endTime(startTime("expensiveCall"))
	expensiveCall()
}

func startTime(message string) (string, time.Time) {
	return message, time.Now()
}

func endTime(msg string, start time.Time) {
	log.Printf("%v: %v\n", msg, time.Since(start))
}

func TestDuration_Measure_Execution_호출함수_이름도_출력하기(t *testing.T) {
	defer executionTime(time.Now())
	expensiveCall()
}

//https://stackoverflow.com/questions/45766572/is-there-an-efficient-way-to-calculate-execution-time-in-golang
func executionTime(start time.Time) {
	elapsed := time.Since(start)

	// Skip this function, and fetch the PC and file for its parent.
	pc, _, _, _ := runtime.Caller(1)

	// Retrieve a function object this functions parent.
	funcObj := runtime.FuncForPC(pc)

	// Regex to extract just the function name (and not the module path).
	runtimeFunc := regexp.MustCompile(`^.*\.(.*)$`)
	name := runtimeFunc.ReplaceAllString(funcObj.Name(), "$1")

	log.Println(fmt.Sprintf("%s took %s", name, elapsed))
}
