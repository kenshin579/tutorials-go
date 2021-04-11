package go_reflect

import (
	"fmt"
	"reflect"

	"github.com/kenshin579/tutorials-go/go-reflect/model"
)

func Example_구조체_필드_순회하기() {
	cat := &model.Cat{
		Name:  "nabi",
		Age:   5,
		Child: []string{"nyang", "kong"},
	}
	IterateStructField(cat)

	//Output:
	//Name: Name / Type: string / Value: nabi / Tag: name
	//Name: Age / Type: int / Value: 5 / Tag: age
	//Name: Child / Type: []string / Value: [nyang kong] / Tag: child
}

func IterateStructField(object interface{}) {
	e := reflect.ValueOf(object).Elem()
	fieldNum := e.NumField()
	for i := 0; i < fieldNum; i++ {
		v := e.Field(i)
		t := e.Type().Field(i)
		fmt.Printf("Name: %s / Type: %s / Value: %v / Tag: %s\n",
			t.Name, t.Type, v.Interface(), t.Tag.Get("custom"))
	}
}
