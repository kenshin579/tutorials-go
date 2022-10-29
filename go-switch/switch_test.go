package go_switch

import "fmt"

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
