package go1_26_test

import (
	"encoding/json"
	"fmt"
	"testing"
	"time"
)

func TestNewWithExpression(t *testing.T) {
	// 기존 방식: 변수를 먼저 선언한 후 주소를 가져옴
	x := 42
	pOld := &x
	fmt.Println("기존 방식:", *pOld) // 42

	// 새로운 방식: new(expr)로 직접 포인터 생성
	pNew := new(42)
	fmt.Println("새로운 방식:", *pNew) // 42

	if *pOld != *pNew {
		t.Errorf("expected %d, got %d", *pOld, *pNew)
	}
}

func TestNewWithStringExpression(t *testing.T) {
	s := new("hello world")
	fmt.Println(*s) // hello world

	if *s != "hello world" {
		t.Errorf("expected 'hello world', got '%s'", *s)
	}
}

func TestNewWithSliceExpression(t *testing.T) {
	// 슬라이스 포인터 생성
	ps := new([]int{11, 12, 13})
	fmt.Println(*ps) // [11 12 13]

	if len(*ps) != 3 || (*ps)[0] != 11 {
		t.Errorf("unexpected slice: %v", *ps)
	}
}

func TestNewWithMapExpression(t *testing.T) {
	// 맵 포인터 생성
	pm := new(map[string]int{"a": 1, "b": 2})
	fmt.Println(*pm) // map[a:1 b:2]

	if (*pm)["a"] != 1 {
		t.Errorf("expected 1, got %d", (*pm)["a"])
	}
}

type Person struct {
	Name string
	Age  *int
}

func yearsSince(t time.Time) int {
	return int(time.Since(t).Hours() / 8766)
}

func TestNewInStructField(t *testing.T) {
	born := time.Date(1990, 1, 1, 0, 0, 0, 0, time.UTC)

	// 기존 방식: 별도 변수 필요
	age := yearsSince(born)
	p1 := Person{Name: "Alice", Age: &age}

	// 새로운 방식: new(expr)로 직접 할당
	p2 := Person{Name: "Alice", Age: new(yearsSince(born))}

	fmt.Printf("기존: %s, %d세\n", p1.Name, *p1.Age)
	fmt.Printf("신규: %s, %d세\n", p2.Name, *p2.Age)

	if *p1.Age != *p2.Age {
		t.Errorf("ages should match: %d vs %d", *p1.Age, *p2.Age)
	}
}

func TestNewWithJSON(t *testing.T) {
	born := time.Date(1990, 1, 1, 0, 0, 0, 0, time.UTC)

	// JSON 마샬링에서 new(expr) 활용
	data, err := json.Marshal(Person{
		Name: "Bob",
		Age:  new(yearsSince(born)),
	})
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(string(data))

	var result Person
	if err := json.Unmarshal(data, &result); err != nil {
		t.Fatal(err)
	}
	if result.Name != "Bob" || result.Age == nil {
		t.Error("unexpected unmarshal result")
	}
}
