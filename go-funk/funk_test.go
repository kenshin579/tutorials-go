package go_funk

import (
	"database/sql"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/thoas/go-funk"
)

// https://pkg.go.dev/github.com/thoas/go-funk#Every
type Account struct {
	ID       int
	Password string
}

type Bar struct {
	Name string `tag_name:"BarName"`
	Bar  *Bar
	Bars []*Bar
}

type Foo struct {
	ID         int
	FirstName  string `tag_name:"tag 1"`
	LastName   string `tag_name:"tag 2"`
	Age        int    `tag_name:"tag 3"`
	Bar        *Bar   `tag_name:"tag 4"`
	Bars       []*Bar
	EmptyValue sql.NullInt64

	BarInterface     interface{}
	BarPointer       interface{}
	GeneralInterface interface{}

	ZeroBoolValue   bool
	ZeroIntValue    int
	ZeroIntPtrValue *int
}

/*
*
Returns true if an element is present in a iteratee (slice, map, string).
*/
func TestContains(t *testing.T) {
	assert.True(t, funk.Contains([]string{"a", "b"}, "a"))
	assert.True(t, funk.ContainsInt([]int{1, 2}, 1))
}

/*
*
ContainsString returns true if a string is present in a iteratee.
*/
func TestContainString(t *testing.T) {
	result := funk.ContainsString([]string{"flo", "gilles"}, "flo")
	assert.True(t, result)

	result = funk.ContainsString([]string{"flo", "gilles"}, "alex")
	assert.False(t, result)
}

/*
*
Every returns true if every element is present in a iteratee.
*/
func TestEvery(t *testing.T) {
	strArr := []string{"go", "java", "c", "python"}
	assert.True(t, funk.Every(strArr, "go", "c"))
	assert.False(t, funk.Every(strArr, "php", "go"))
	assert.False(t, funk.Every(strArr, "php", "c++"))
}

/*
*
Some returns true if at least one element is present in an iteratee.
*/
func TestSome(t *testing.T) {
	strArr := []string{"go", "java", "c", "python"}
	assert.True(t, funk.Some(strArr, "go", "c"))
	assert.True(t, funk.Some(strArr, "php", "go"))
	assert.False(t, funk.Some(strArr, "php", "c++"))
}

/*
*
함수 조건이 참인 값만 필터링한다
*/
func TestFilter(t *testing.T) {
	filter := funk.Filter([]int{1, 2, 3, 4}, func(x int) bool {
		return x%2 == 0
	})
	assert.Equal(t, []int{2, 4}, filter)
}

/*
*
중복된 값은 제거하고 uniq 한 값을 반한환다
*/
func TestUniq(t *testing.T) {
	uniq := funk.Uniq([]int{0, 1, 1, 2, 3, 0, 0, 12})
	assert.Equal(t, []int{0, 1, 2, 3, 12}, uniq)
}

/*
*
2 collection에서 Interaction 값을 반환한다
*/
func TestIntersectString(t *testing.T) {
	result := funk.IntersectString([]string{"foo", "bar", "hello", "bar"}, []string{"foo", "bar"})
	assert.Equal(t, []string{"foo", "bar"}, result)
}

/*
*
맨 앞에 있는 값을 반환한다
*/
func TestHead(t *testing.T) {
	assert.Equal(t, 1, funk.Head([]int{1, 2, 3, 4}))
}

/*
*
Creates an array/slice with n elements dropped from the beginning.
- collection에서 첫 N은 drop하고 나머지를 반환한다
*/
func TestDrop(t *testing.T) {
	drop := funk.Drop([]int{0, 1, 1, 2, 3, 0, 0, 12}, 3)
	assert.Equal(t, []int{2, 3, 0, 0, 12}, drop)
}

/*
*
함수 조건이 참이 요소를 찾는다
*/
func TestFind(t *testing.T) {
	find := funk.Find([]int{1, 2, 3, 4}, func(x int) bool {
		return x%2 == 0
	})
	assert.Equal(t, 2, find)

	// list에 없는 경우 nil을 반환한다
	find = funk.Find([]int{1, 2, 3, 4}, func(x int) bool {
		return x == 5
	})
	assert.Equal(t, nil, find)
}

/*
*
Range over an iteratee (map, slice).
*/
func TestForEach(t *testing.T) {
	funk.ForEach([]int{1, 2, 3, 4}, func(x int) {
		fmt.Println(x)
	})
}

/*
*
Manipulates an iteratee (map, slice) and transforms it to another type:

- map -> slice
- map -> map
- slice -> map
- slice -> slice
*/
func TestMap(t *testing.T) {
	result := funk.Map([]int{1, 2, 3, 4}, func(x int) int {
		return x * 2
	})
	assert.Equal(t, []int{2, 4, 6, 8}, result)

	result = funk.Map([]int{1, 2, 3, 4}, func(x int) string {
		return "Hello"
	})
	assert.Equal(t, []string{"Hello", "Hello", "Hello", "Hello"}, result)

	result = funk.Map([]int{1, 2, 3, 4}, func(x int) (int, int) {
		return x, x
	})
	assert.Equal(t, map[int]int{1: 1, 2: 2, 3: 3, 4: 4}, result)

	mapping := map[int]string{
		1: "Florent",
		2: "Gilles",
	}

	result = funk.Map(mapping, func(k int, v string) int {
		return k
	})
	assert.Equal(t, []int{1, 2}, result)

	result = funk.Map(mapping, func(k int, v string) (string, string) {
		return fmt.Sprintf("%d", k), v
	})
	assert.Equal(t, map[string]string{"1": "Florent", "2": "Gilles"}, result)

}

func TestToMapOfStruct(t *testing.T) {
	f1 := Foo{
		ID:        1,
		FirstName: "Dark",
		LastName:  "Vador",
		Age:       30,
		Bar: &Bar{
			Name: "Test",
		},
	}

	f2 := Foo{
		ID:        1,
		FirstName: "Light",
		LastName:  "Vador",
		Age:       30,
		Bar: &Bar{
			Name: "Test",
		},
	}

	arrayResults := [4]Foo{f1, f1, f2, f2}

	instanceMapByFirstName := funk.ToMap(arrayResults, "FirstName")
	mappingByFirstName, ok := instanceMapByFirstName.(map[string]Foo)
	assert.True(t, ok)
	assert.True(t, len(mappingByFirstName) == 2)

	for _, result := range arrayResults {
		item, ok := mappingByFirstName[result.FirstName]
		assert.True(t, ok)
		assert.Equal(t, result.FirstName, item.FirstName)
	}
}

func TestMap_Chunk(t *testing.T) {
	var elements = []string{"abc", "def", "fgi", "adi"}

	elementsMap := funk.Map(
		funk.Chunk(elements, 2),
		func(x []string) (string, string) { // Slice to Map
			return x[0], x[1]
		},
	)

	fmt.Println(elementsMap)
}

func TestGet(t *testing.T) {
	account := &Account{
		ID:       2,
		Password: "1111",
	}

	id := funk.Get(account, "ID")
	assert.Equal(t, 2, id)
}

func TestIndexOf(t *testing.T) {
	assert.Equal(t, 0, funk.IndexOf([]string{"hello", "world"}, "hello"))
	assert.Equal(t, -1, funk.IndexOf([]string{"hello", "world"}, "Account"))
}

func Test_Chain(t *testing.T) {
	ints := []int{0, 1, 2, 3, 4, 5, 6, 7, 8, 9}

	values := funk.Chain(ints).
		Filter(func(x int) bool {
			return x%2 == 0
		}).
		Map(func(x int) int {
			return x * 2
		}).
		Drop(2).Value()

	assert.Equal(t, []int{8, 12, 16}, values)
}

func Test_Last(t *testing.T) {
	ints := []int{0, 1, 2, 3, 4, 5, 6, 7, 8, 9}

	assert.Equal(t, 9, funk.Last(ints))
}
