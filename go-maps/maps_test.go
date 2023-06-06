package go_maps

import (
	"fmt"
	"testing"

	"github.com/kenshin579/tutorials-go/go-maps/model"
	"github.com/stretchr/testify/assert"
	"golang.org/x/exp/maps"
)

func Example_Maps_Copy() {
	src := map[int]string{
		200: "foo",
		300: "bar",
	}
	dest := map[int]string{}

	maps.Copy(dest, src)
	fmt.Println(dest)

	dest2 := map[int]string{
		200: "will be overwritten",
	}
	maps.Copy(dest2, src)
	fmt.Println(dest2)

	// Output:
	// map[200:foo 300:bar]
	// map[200:foo 300:bar] //todo - 왜 overwrite가 안되나?
}

func ContainsKey[K, V comparable](m map[K]V, target V) bool {
	for _, v := range m {
		if v == target {
			return true
		}
	}
	return false
}

func Test_Maps_Key_Values(t *testing.T) {
	m := map[string]int{
		"a": 1,
		"b": 2,
		"c": 3,
	}

	assert.Equal(t, []string{"a", "b", "c"}, maps.Keys(m))
	assert.Equal(t, []int{1, 2, 3}, maps.Values(m))

	m2 := map[string]int{}
	maps.Copy(m2, m)

	fmt.Println(maps.Keys(m2))
	fmt.Println(maps.Values(m2))

	assert.True(t, ContainsKey(m2, 3))
	assert.False(t, ContainsKey(m2, 4))
}

func Test_Convert_Map_To_Struct(t *testing.T) {

}

func Test_Convert_Struct_To_Map(t *testing.T) {
	employee := model.Employee{
		Name:    "frank",
		Age:     20,
		Address: "address1",
	}

}
