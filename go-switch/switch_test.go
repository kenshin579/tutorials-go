package go_switch

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

/*
golang에서는 case 구문에 break를 넣지 않아도 된다
*/
func Example_Switch() {
	n := 2

	switch n {
	case 3:
		fmt.Println("n is ", n)
	case 2:
		fmt.Println("n is ", n)
	case 1:
		fmt.Println("n is ", n)
	}

	//Output:
	//n is  2
}

func Example_Switch_Case_조건문에_여러_값이_있는_경우() {
	str := "blue"

	switch str {
	case "red":
		fmt.Println("Stop")
	case "yellow":
		fmt.Println("caution")
	case "green", "blue": //OR 역할을 한다
		fmt.Println("Go")
	default:
		fmt.Println("wrong")
	}

	//Output:
	//Go
}

type Status string

const (
	StatusReady = "Ready"
	StatusDone  = "Done"
)

func Example_Switch_Case에_() {
	status1 := StatusReady
	status2 := StatusDone

	//if-else 처럼 사용한다
	switch {
	case status1 == StatusReady && status2 == StatusDone:
		fmt.Println("ready and done")
	case status1 == StatusReady:
		fmt.Println("ready only")
	default:
		fmt.Println("no match")
	}

	//Output:
	//ready and done
}

/*
Fallthrough 사용하면 자바에서 break 없이 작성하는 것처럼 동작한다
- 그냥 아래 case를 실행한다
- fallthrough를 사용하면 그 다음 case 문에 작성한 조건문을 무시한다
*/
func Example_Switch_Fallthrough() {
	cnt := 6
	switch cnt {
	case 4:
		fmt.Println("was <= 4")
		fallthrough
	case 5:
		fmt.Println("was <= 5")
		fallthrough
	case 6:
		fmt.Println("was <= 6")
		fallthrough
	case 7:
		fmt.Println("was <= 7")
		fallthrough
	case 8:
		fmt.Println("was <= 8")
		fallthrough
	default:
		fmt.Println("default case")
	}

	//Output:
	//was <= 6
	//was <= 7
	//was <= 8
	//default case
}

func Example_Fallthrough2() {
	score := 3

	switch score {
	case 3:
		score *= 100
		fallthrough
	case 2:
		fmt.Println(score) //실행이 된다
	case 1:
		fmt.Println(score * 10) //실행되지 않는다
	}

	//Output:
	//300
}

func Example() {
	cmdTypes := []string{"type1", "type2", "type3"}

	switch {
	case containAllTypes(cmdTypes, "type1", "type2"):
		fmt.Println("type1 & type2")
	case containAllTypes(cmdTypes, "type1"):
		fmt.Println("type1")
	default:
		fmt.Println("unknown")
	}

	//Output:
	//type1 & type2
}

/*
https://www.cloudhadoop.com/golang-map-type-tutorial-examples/
*/
func containAllTypes(cmdTypes []string, containTypes ...string) bool {
	count := 0

	typeMap := make(map[string]bool)
	for _, x := range cmdTypes {
		typeMap[x] = true
	}

	for _, cmdType := range containTypes {
		if _, found := typeMap[cmdType]; found {
			count++
		}
	}
	return count == len(containTypes)
}

func Test_containAllTypes(t *testing.T) {
	types := containAllTypes([]string{"type1", "type2", "type3"}, "type1", "type2")
	assert.True(t, types)
}
