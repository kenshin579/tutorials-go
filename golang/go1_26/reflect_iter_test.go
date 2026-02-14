package go1_26_test

import (
	"fmt"
	"net/http"
	"reflect"
	"testing"
)

func TestReflectFields(t *testing.T) {
	type User struct {
		Name  string
		Email string
		Age   int
	}

	typ := reflect.TypeFor[User]()
	fmt.Printf("Type: %s\n", typ.Name())
	fmt.Println("--- Fields ---")

	// 새로운 방식: Type.Fields() 반복자
	count := 0
	for f := range typ.Fields() {
		fmt.Printf("  %s: %s\n", f.Name, f.Type)
		count++
	}

	if count != 3 {
		t.Errorf("expected 3 fields, got %d", count)
	}
}

func TestReflectMethods(t *testing.T) {
	// 포인터 타입으로 메서드 조회 (http.Client의 메서드는 포인터 리시버)
	typ := reflect.TypeFor[*http.Client]()
	fmt.Printf("Type: %s\n", typ)
	fmt.Println("--- Methods ---")

	// 새로운 방식: Type.Methods() 반복자
	count := 0
	for m := range typ.Methods() {
		fmt.Printf("  %s\n", m.Name)
		count++
	}

	if count == 0 {
		t.Error("expected at least one method")
	}
	fmt.Printf("Total methods: %d\n", count)
}

func TestReflectFieldsComparison(t *testing.T) {
	type Config struct {
		Host    string
		Port    int
		Debug   bool
		Timeout float64
	}

	typ := reflect.TypeFor[Config]()

	// 기존 방식: 인덱스로 순회
	fmt.Println("--- 기존 방식 (인덱스) ---")
	for i := 0; i < typ.NumField(); i++ {
		f := typ.Field(i)
		fmt.Printf("  %s: %s\n", f.Name, f.Type)
	}

	// 새로운 방식: range 반복자
	fmt.Println("--- 새로운 방식 (반복자) ---")
	for f := range typ.Fields() {
		fmt.Printf("  %s: %s\n", f.Name, f.Type)
	}
}

func TestReflectMethodsOfInterface(t *testing.T) {
	typ := reflect.TypeFor[error]()
	fmt.Printf("Interface: %s\n", typ.Name())
	fmt.Println("--- Methods ---")

	for m := range typ.Methods() {
		fmt.Printf("  %s: %s\n", m.Name, m.Type)
	}
}
