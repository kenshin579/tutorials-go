package go_custom_errors

import (
	"errors"
	"fmt"
	"os"
	"testing"

	"github.com/kenshin579/tutorials-go/go-custom-errors/model"
)

func ExampleCreating_Error_New() {
	error1 := errors.New("error occured")
	fmt.Println(error1)

	//Output: error occured
}

func ExampleCreating_Error_fmt_Errorf() {
	error1 := fmt.Errorf("err is: %s", "database connection issue")
	fmt.Println(error1)

	//Output: err is: database connection issue
}

func TestIgnoreError_처리(t *testing.T) {
	file, _ := os.Open("non-existing.txt") //_로 error를 무시하도록 처리함
	fmt.Println(file)
}

func ExampleCreating_Custom_Error() {
	fmt.Println(model.ErrRequestUser)
	//Output: 400:10100:Request is invalid
}

func ExampleCreating_Custom_Error2() {
	err := validate("", "")
	if err != nil {
		if err, ok := err.(*model.InputError); ok {
			fmt.Println(err)
			fmt.Printf("Missing Field is %s\n", err.GetMissingField())
		}
	}

	//Output:
	//Name is mandatory
	//Missing Field is name
}

func validate(name, gender string) error {
	if name == "" {
		return &model.InputError{Message: "Name is mandatory", MissingField: "name"}
	}
	if gender == "" {
		return &model.InputError{Message: "Gender is mandatory", MissingField: "gender"}
	}
	return nil
}
