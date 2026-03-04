package go_reflect

import (
	"reflect"
	"testing"
)

type BenchStruct struct {
	Name string
	Age  int
}

func (b BenchStruct) GetName() string {
	return b.Name
}

// 필드 읽기: 직접 접근 vs reflect
func BenchmarkFieldDirect(b *testing.B) {
	s := BenchStruct{Name: "Go", Age: 10}
	var name string
	for i := 0; i < b.N; i++ {
		name = s.Name
	}
	_ = name
}

func BenchmarkFieldReflect(b *testing.B) {
	s := BenchStruct{Name: "Go", Age: 10}
	v := reflect.ValueOf(s)
	var name string
	for i := 0; i < b.N; i++ {
		name = v.Field(0).String()
	}
	_ = name
}

// 메서드 호출: 직접 호출 vs reflect
func BenchmarkMethodDirect(b *testing.B) {
	s := BenchStruct{Name: "Go", Age: 10}
	var name string
	for i := 0; i < b.N; i++ {
		name = s.GetName()
	}
	_ = name
}

func BenchmarkMethodReflect(b *testing.B) {
	s := BenchStruct{Name: "Go", Age: 10}
	v := reflect.ValueOf(s)
	method := v.MethodByName("GetName")
	for i := 0; i < b.N; i++ {
		method.Call(nil)
	}
}

// 구조체 생성: 리터럴 vs reflect.New
func BenchmarkCreateDirect(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_ = BenchStruct{Name: "Go", Age: 10}
	}
}

func BenchmarkCreateReflect(b *testing.B) {
	t := reflect.TypeOf(BenchStruct{})
	for i := 0; i < b.N; i++ {
		_ = reflect.New(t).Elem()
	}
}

// DeepEqual vs 직접 비교
func BenchmarkEqualDirect(b *testing.B) {
	s1 := BenchStruct{Name: "Go", Age: 10}
	s2 := BenchStruct{Name: "Go", Age: 10}
	for i := 0; i < b.N; i++ {
		_ = s1.Name == s2.Name && s1.Age == s2.Age
	}
}

func BenchmarkEqualDeepEqual(b *testing.B) {
	s1 := BenchStruct{Name: "Go", Age: 10}
	s2 := BenchStruct{Name: "Go", Age: 10}
	for i := 0; i < b.N; i++ {
		reflect.DeepEqual(s1, s2)
	}
}
