package go_jq

import (
	"encoding/json"
	"fmt"
	"log"
	"testing"

	_ "embed"

	"github.com/itchyny/gojq"
)

//go:embed json/ex1.json
var ex1JsonFile []byte

type keyValue struct {
	Key   string
	Value string
}

func Example_Gojq() {
	input := map[string]interface{}{"foo": []interface{}{1, 2, 3}}
	//	jsonObject := `
	//	"foo": [1,2,3]
	//`

	query, err := gojq.Parse(".foo | ..")
	if err != nil {
		log.Fatalln(err)
	}

	iter := query.Run(input) // or query.RunWithContext
	for {
		v, ok := iter.Next()
		if !ok {
			break
		}
		if err, ok := v.(error); ok {
			log.Fatalln(err)
		}
		fmt.Printf("%#v\n", v)
	}

	//Output:
	//[]interface {}{1, 2, 3}
	//1
	//2
	//3
}

// todo: 이거 잘 안됨
func Test_Parse_Key(t *testing.T) {
	var keyValue keyValue
	if err := json.Unmarshal(ex1JsonFile, &keyValue); err != nil {
		panic(err)
	}
	fmt.Println("keyvalue", keyValue)

	query, err := gojq.Parse(".key")
	if err != nil {
		log.Fatalln(err)
	}

	iter := query.Run(keyValue)

	for {
		v, ok := iter.Next()
		if !ok {
			break
		}
		if err, ok := v.(error); ok {
			log.Fatalln(err)
		}
		fmt.Printf("%+v\n", v)
	}
}
