package go_funk

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/thoas/go-funk"
)

/*
filter
map, toMap, set, mustset,
get, getorelse,
reduce
prune, prunebytag
forEachRight, chunk, flattendeep, uniq, initial, tail, shuffle, tail,sum
reduce
shard
substract
intersect
difference
zip
subset
equal
contains
keys
values
indexOf, lastIndexOf

*/

type test struct {
	ID       int
	Password string
}

func TestContains(t *testing.T) {
	assert.True(t, funk.Contains([]string{"a", "b"}, "a"))
	assert.True(t, funk.ContainsInt([]int{1, 2}, 1))
}

func TestEvery(t *testing.T) {
	strArr := []string{"go", "java", "c", "python"}
	assert.True(t, funk.Every(strArr, "go", "c"))
	assert.False(t, funk.Every(strArr, "php", "go"))
	assert.False(t, funk.Every(strArr, "php", "c++"))
}

func TestSome(t *testing.T) {
	strArr := []string{"go", "java", "c", "python"}
	assert.True(t, funk.Every(strArr, "go", "c"))
	assert.False(t, funk.Every(strArr, "php", "go"))
	assert.False(t, funk.Every(strArr, "php", "c++"))
}

func TestAny(t *testing.T) {
	strArr := []string{"go", "java", "c", "python"}
	assert.True(t, funk.Every(strArr, "go", "c"))
	assert.False(t, funk.Every(strArr, "php", "go"))
	assert.False(t, funk.Every(strArr, "php", "c++"))
}

func TestAll(t *testing.T) {

}

func Test10(t *testing.T) {
	assert.True(t, funk.Contains([]string{"hello", "world"}, "hello"))

	testStruct := &test{
		ID:       1,
		Password: "0000",
	}
	contains := funk.Contains([]*test{testStruct}, testStruct)
	fmt.Println(contains)
	contains = funk.Contains([]*test{testStruct}, nil)
	fmt.Println(contains)

	testStruct2 := &test{
		ID:       2,
		Password: "1111",
	}

	get := funk.Get(testStruct2, "ID")
	fmt.Println(get)

	contains = funk.Contains("hello", "lo")
	fmt.Println(contains)
	contains = funk.Contains("hello", "world")
	fmt.Println(contains)

	contains = funk.Contains(map[int]string{1: "test"}, 1)
	fmt.Println(contains)

	indexOf := funk.IndexOf([]string{"hello", "world"}, "hello")
	fmt.Println(indexOf)
	indexOf = funk.IndexOf([]string{"hello", "world"}, "test")
	fmt.Println(indexOf)

	indexOf = funk.LastIndexOf([]string{"hello", "hello", "world"}, "hello")
	fmt.Println(indexOf)
	indexOf = funk.LastIndexOf([]string{"hello", "world"}, "test")
	fmt.Println(indexOf)

	//ToMap으로 특정 키를 기반으로 객체를 map 변환 - todo: 이거는 어떻게 사용할 수 있나?
	results := []*test{testStruct, testStruct2}
	mapping := funk.ToMap(results, "ID")
	fmt.Println(mapping)

}

func Test2(t *testing.T) {

	filter := funk.Filter([]int{1, 2, 3, 4}, func(x int) bool {
		return x%2 == 0
	})
	fmt.Println(filter)

	find := funk.Find([]int{1, 2, 3, 4}, func(x int) bool {
		return x%2 == 0
	})
	fmt.Println(find)

	intMap := funk.Map([]int{1, 2, 3, 4}, func(x int) int {
		return x * 2
	})
	fmt.Println(intMap)

	stringMap := funk.Map([]int{1, 2, 3, 4}, func(x int) string {
		return "Hello"
	})
	fmt.Println(stringMap)

	intMap = funk.Map([]int{1, 2, 3, 4}, func(x int) (int, int) {
		return x, x
	})
	fmt.Println(intMap)

}

func Test3(t *testing.T) {
	mapping := map[int]string{
		1: "hello",
		2: "world",
	}

	intMap := funk.Map(mapping, func(k int, v string) int {
		return k
	})
	fmt.Println(intMap)

	stringMap := funk.Map(mapping, func(k int, v string) (string, string) {
		return fmt.Sprintf("%d", k), v
	})
	fmt.Println(stringMap)
}
