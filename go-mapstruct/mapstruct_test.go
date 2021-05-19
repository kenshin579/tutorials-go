package go_mapstruct

import (
	"fmt"

	"github.com/mitchellh/mapstructure"
)

func Example_Convert_Map_To_Struct() {
	type Person struct {
		Name   string
		Age    int
		Emails []string
		Extra  map[string]string
	}

	inputMap := map[string]interface{}{
		"name":   "Mitchell",
		"age":    91,
		"emails": []string{"one", "two", "three"},
		"extra": map[string]string{
			"twitter": "mitchellh",
		},
	}

	var person Person
	err := mapstructure.Decode(inputMap, &person)
	if err != nil {
		panic(err)
	}

	fmt.Printf("%#v\n", person)
	// Output:
	// go_mapstruct.Person{Name:"Mitchell", Age:91, Emails:[]string{"one", "two", "three"}, Extra:map[string]string{"twitter":"mitchellh"}}
}

func Example_Convert_Map_To_Struct_Weakly() {
	type Person struct {
		Name   string
		Age    int
		Emails []string
	}

	inputMap := map[string]interface{}{
		"name":   123,
		"age":    "42", //WeaklyTypedInput:true 설정으로 인해 string 값이지만, 숫자로 변환이 가능함
		"emails": map[string]interface{}{},
	}

	var person Person
	config := &mapstructure.DecoderConfig{
		WeaklyTypedInput: true,
		Result:           &person,
	}

	decoder, err := mapstructure.NewDecoder(config)
	if err != nil {
		panic(err)
	}

	err = decoder.Decode(inputMap)
	if err != nil {
		panic(err)
	}

	fmt.Printf("%#v\n", person)
	// Output:
	//go_mapstruct.Person{Name:"123", Age:42, Emails:[]string{}}
}

func Example_Tag와_Map_key값_이름_매팅으로_변환하기() {
	type Person struct {
		Name string `mapstructure:"person_name"`
		Age  int    `mapstructure:"person_age"`
	}

	inputMap := map[string]interface{}{
		"person_name1": "Mitchell",
		"person_age":   91,
	}

	var person Person
	err := mapstructure.Decode(inputMap, &person)
	if err != nil {
		panic(err)
	}

	fmt.Printf("%#v\n", person)
	// Output:
	// go_mapstruct.Person{Name:"Mitchell", Age:91}
}

func Example_Squash_Tag를_이용해서_Nested_Struct로_변환하기() {
	type Family struct {
		LastName string
	}
	type Location struct {
		City string
	}
	type Person struct {
		Family    `mapstructure:",squash"`
		Location  `mapstructure:",squash"`
		FirstName string
	}

	input := map[string]interface{}{
		"FirstName": "Mitchell",
		"LastName":  "Hashimoto",
		"City":      "San Francisco",
	}

	var person Person
	err := mapstructure.Decode(input, &person)
	if err != nil {
		panic(err)
	}

	fmt.Println(person)
	// Output:
	// {{Hashimoto} {San Francisco} Mitchell}
}

func Example_Remain_Tag() {
	type Person struct {
		Name  string
		Age   int
		Other map[string]interface{} `mapstructure:",remain"`
	}

	inputMap := map[string]interface{}{
		"name":  "Mitchell",
		"age":   91,
		"email": "mitchell@example.com",
	}

	var person Person
	err := mapstructure.Decode(inputMap, &person)
	if err != nil {
		panic(err)
	}

	fmt.Printf("%#v\n", person)
	// Output:
	//go_mapstruct.Person{Name:"Mitchell", Age:91, Other:map[string]interface {}{"email":"mitchell@example.com"}}
}

//Omitempty tag가 없으면 기본 값으로 초기화가 됨
func Example_Omitempty_Tag로_Map에_없는_값은_skip함() {
	type Family struct {
		LastName string
	}
	type Location struct {
		City string
	}
	type Person struct {
		*Family   `mapstructure:",omitempty"`
		*Location `mapstructure:",omitempty"`
		Age       int `mapstructure:",omitempty"`
		FirstName string
	}

	result := &map[string]interface{}{}
	input := Person{FirstName: "Somebody"}
	err := mapstructure.Decode(input, &result)
	if err != nil {
		panic(err)
	}

	fmt.Printf("%+v\n", result)

	// Output:
	//&map[FirstName:Somebody]
}
