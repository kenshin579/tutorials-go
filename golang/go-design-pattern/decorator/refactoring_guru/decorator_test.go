package refactoring_guru

import "fmt"

func Example_Decorator() {
	/*
		pizza - 인터페이스
		veggieMania - 구조체
		tomatoTopping - 구조체 (인퍼테이스를 포함할 수 있음)
		cheeseTopping - 구조체 (인퍼테이스를 포함할 수 있음)
	*/

	pizza := &veggeMania{}

	fmt.Println("veggeMania pizza price", pizza.getPrice())

	//Add cheese topping
	pizzaWithCheese := &cheeseTopping{
		pizza: pizza,
	}

	fmt.Println("pizzaWithCheese pizza price", pizzaWithCheese.getPrice())

	//Add tomato topping
	pizzaWithCheeseAndTomato := &tomatoTopping{
		pizza: pizzaWithCheese,
	}

	fmt.Printf("Price of veggeMania with tomato and cheese topping is %d\n",
		pizzaWithCheeseAndTomato.getPrice())

	//Output:
	//veggeMania pizza price 15
	//pizzaWithCheese pizza price 25
	//Price of veggeMania with tomato and cheese topping is 32
}
