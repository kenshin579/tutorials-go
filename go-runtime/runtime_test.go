package go_runtime

import (
	"fmt"
	"runtime"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func bar() string {
	callDepth := 1
	caller, _, _, _ := runtime.Caller(callDepth)
	callerName := runtime.FuncForPC(caller).Name()
	return callerName
}

func Test_실행중인_함수이름_가져오기(t *testing.T) {
	result := bar()
	assert.Equal(t,
		"github.com/kenshin579/tutorials-go/go-runtime.Test_실행중인_함수이름_가져오기",
		result)
}

/*
https://www.golangprograms.com/example-stack-and-caller-from-runtime-package.html
https://vimsky.com/examples/detail/golang-ex-runtime---Callers-function.html
*/
func Test_runtime_Stack(t *testing.T) {
	fmt.Println("######### STACK ################")
	stackExample()
	fmt.Println("\n\n######### CALLER ################")
	First()
}

func stackExample() {
	stackSlice := make([]byte, 512)
	s := runtime.Stack(stackSlice, false)
	fmt.Printf("\n%s", stackSlice[0:s])
}

func First() {
	Second()
}

func Second() {
	//Third()
	Third2()
}

func Third() {
	for depthIndex := 0; depthIndex < 5; depthIndex++ {
		caller, _, line, _ := runtime.Caller(depthIndex)
		callerName := runtime.FuncForPC(caller).Name()

		fmt.Printf("[%d] callerName: %s :%d\n", depthIndex, callerName, line)
	}
}

const (
	maxDepth = 50
)

func Third2() {
	var pcs [maxDepth]uintptr
	numStack := runtime.Callers(0, pcs[:])
	fmt.Println("numStack", numStack)
	for i := 0; i < numStack; i++ {
		function := runtime.FuncForPC(pcs[i])
		if function == nil {
			break
		}

		_, line := function.FileLine(pcs[i])
		fmt.Printf("%d callerName: %s :%d\n", i, function.Name(), line)

	}
}

func Test_runtime_callers(t *testing.T) {
	First()
}

/*
https://golang.org/src/runtime/example_test.go
*/
func Test_runtime_Frames(t *testing.T) {
	c := func() {
		// Ask runtime.Callers for up to 10 PCs, including runtime.Callers itself.
		pc := make([]uintptr, 10)
		n := runtime.Callers(0, pc)
		if n == 0 {
			// No PCs available. This can happen if the first argument to
			// runtime.Callers is large.
			//
			// Return now to avoid processing the zero Frame that would
			// otherwise be returned by frames.Next below.
			return
		}

		pc = pc[:n] // pass only valid pcs to runtime.CallersFrames
		frames := runtime.CallersFrames(pc)

		// Loop to get frames.
		// A fixed number of PCs can expand to an indefinite number of Frames.
		for {
			frame, more := frames.Next()

			// Process this frame.
			//
			// To keep this example's output stable
			// even if there are changes in the testing package,
			// stop unwinding when we leave package runtime.
			if !strings.Contains(frame.File, "runtime/") {
				break
			}
			fmt.Printf("- more:%v | %s\n", more, frame.Function)

			// Check whether there are more frames to process after this one.
			if !more {
				break
			}
		}
	}

	b := func() { c() }
	a := func() { b() }

	a()
}
