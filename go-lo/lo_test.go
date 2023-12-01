package go_lo

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/samber/lo"
	lop "github.com/samber/lo/parallel"
	"github.com/stretchr/testify/assert"

	"testing"
)

// https://github.com/samber/lo

func Test_Supported_helpers_for_slices(t *testing.T) {
	t.Run("Filter", func(t *testing.T) {
		even := lo.Filter([]int{1, 2, 3, 4}, func(x int, index int) bool {
			return x%2 == 0
		})
		assert.Equal(t, []int{2, 4}, even)
	})

	t.Run("FilterMap", func(t *testing.T) {
		// Returns a slice which obtained after both filtering and mapping using the given callback function.

		matching := lo.FilterMap([]string{"cpu", "gpu", "mouse", "keyboard"}, func(x string, _ int) (string, bool) {
			if strings.HasSuffix(x, "pu") {
				return "xpu", true
			}
			return "", false
		})
		assert.Equal(t, []string{"xpu", "xpu"}, matching)
	})

	t.Run("Map", func(t *testing.T) {
		// Manipulates a slice of one type and transforms it into a slice of another type:
		result := lo.Map([]int64{1, 2, 3, 4}, func(x int64, index int) string {
			return strconv.FormatInt(x, 10)
		})
		assert.Equal(t, []string{"1", "2", "3", "4"}, result)

		// Parallel processing: like lo.Map(), but the mapper function is called in a goroutine. Results are returned in the same order.
		result = lop.Map([]int64{1, 2, 3, 4}, func(x int64, _ int) string {
			return strconv.FormatInt(x, 10)
		})
		assert.Equal(t, []string{"1", "2", "3", "4"}, result)
	})

	t.Run("Uniq, UniqBy", func(t *testing.T) {
		uniqValues := lo.Uniq([]int{1, 2, 2, 1})
		assert.Equal(t, []int{1, 2}, uniqValues)

		uniqValues = lo.UniqBy([]int{0, 1, 2, 3, 4, 5}, func(i int) int {
			return i % 3
		})
		assert.Equal(t, []int{0, 1, 2}, uniqValues)
	})

	t.Run("Reduce", func(t *testing.T) {
		sum := lo.Reduce([]int{1, 2, 3, 4}, func(agg int, item int, _ int) int {
			return agg + item
		}, 0)
		assert.Equal(t, 10, sum)
	})

	t.Run("GroupBy", func(t *testing.T) {
		// Returns an object composed of keys generated from the results of running each element of collection through iteratee.
		groups := lo.GroupBy([]int{0, 1, 2, 3, 4, 5}, func(i int) int {
			return i % 3
		})
		assert.Equal(t, map[int][]int{0: {0, 3}, 1: {1, 4}, 2: {2, 5}}, groups)
	})

	t.Run("Compact", func(t *testing.T) {
		// Returns a slice of all non-zero elements.
		in := []string{"", "foo", "", "bar", ""}

		slice := lo.Compact[string](in)
		assert.Equal(t, []string{"foo", "bar"}, slice)
	})

	t.Run("Flatten", func(t *testing.T) {
		flat := lo.Flatten([][]int{{0, 1}, {2, 3, 4, 5}})
		assert.Equal(t, []int{0, 1, 2, 3, 4, 5}, flat)
	})

	t.Run("Reverse", func(t *testing.T) {
		reverseOrder := lo.Reverse([]int{0, 1, 2, 3, 4, 5})
		assert.Equal(t, []int{5, 4, 3, 2, 1, 0}, reverseOrder)
	})

	t.Run("Associate(alias:SliceToMap)", func(t *testing.T) {
		type foo struct {
			baz string
			bar int
		}
		in := []*foo{{baz: "apple", bar: 1}, {baz: "banana", bar: 2}}

		aMap := lo.Associate(in, func(f *foo) (string, int) {
			return f.baz, f.bar
		})

		assert.Equal(t, map[string]int{"apple": 1, "banana": 2}, aMap)
	})

	t.Run("SliceToMap", func(t *testing.T) {
		type foo struct {
			baz string
			bar int
		}
		in := []*foo{{baz: "apple", bar: 1}, {baz: "banana", bar: 2}}

		aMap := lo.SliceToMap(in, func(f *foo) (string, int) {
			return f.baz, f.bar
		})

		assert.Equal(t, map[string]int{"apple": 1, "banana": 2}, aMap)
	})
}

func Test_Supported_helpers_for_intersection(t *testing.T) {

	t.Run("Contains, ContainsBy", func(t *testing.T) {
		present := lo.Contains([]int{0, 1, 2, 3, 4, 5}, 5)
		assert.True(t, present)

		present = lo.ContainsBy([]int{0, 1, 2, 3, 4, 5}, func(x int) bool {
			return x == 3
		})
		assert.True(t, present)
	})

	// Returns true if all elements of a subset are contained into a collection or if the subset is empty.
	t.Run("Every", func(t *testing.T) {
		ok := lo.Every([]int{0, 1, 2, 3, 4, 5}, []int{0, 2})
		assert.True(t, ok)

		ok = lo.Every([]int{0, 1, 2, 3, 4, 5}, []int{0, 6})
		assert.False(t, ok)
	})

	t.Run("Intersect", func(t *testing.T) {
		result := lo.Intersect([]int{0, 1, 2, 3, 4, 5}, []int{0, 2})
		assert.Equal(t, []int{0, 2}, result)

		result = lo.Intersect([]int{0, 1, 2, 3, 4, 5}, []int{0, 6})
		assert.Equal(t, []int{0}, result)

		result = lo.Intersect([]int{0, 1, 2, 3, 4, 5}, []int{-1, 6})
		assert.Equal(t, []int{}, result)
	})

	t.Run("Difference", func(t *testing.T) {
		left, right := lo.Difference([]int{0, 1, 2, 3, 4, 5}, []int{0, 2, 6})
		assert.Equal(t, []int{1, 3, 4, 5}, left)
		assert.Equal(t, []int{6}, right)

		left, right = lo.Difference([]int{0, 1, 2, 3, 4, 5}, []int{0, 1, 2, 3, 4, 5})
		assert.Equal(t, []int{}, left)
		assert.Equal(t, []int{}, right)
	})

	t.Run("Without", func(t *testing.T) {
		subset := lo.Without([]int{0, 2, 10}, 2)
		assert.Equal(t, []int{0, 10}, subset)

		subset = lo.Without([]int{0, 2, 10}, 0, 1, 2, 3, 4, 5)
		assert.Equal(t, []int{10}, subset)
	})

}

func Test_Searching_helpers(t *testing.T) {

	t.Run("", func(t *testing.T) {
		found := lo.IndexOf([]int{0, 1, 2, 1, 2, 3}, 2)
		assert.Equal(t, 2, found)

		notFound := lo.IndexOf([]int{0, 1, 2, 1, 2, 3}, 6)
		assert.Equal(t, -1, notFound)
	})

	t.Run("Find", func(t *testing.T) {
		str, ok := lo.Find([]string{"a", "b", "c", "d"}, func(i string) bool {
			return i == "b"
		})
		assert.Equal(t, "b", str)
		assert.True(t, ok)

		str, ok = lo.Find([]string{"foobar"}, func(i string) bool {
			return i == "b"
		})
		assert.Equal(t, "", str)
		assert.False(t, ok)
	})

}

func Test_Supported_helpers_for_maps(t *testing.T) {
	t.Run("Keys, Values", func(t *testing.T) {
		m := map[string]int{"foo": 1, "bar": 2}

		// Creates an array of the map keys.
		keys := lo.Keys[string, int](m)
		assert.Equal(t, []string{"foo", "bar"}, keys)

		// Creates an array of the map values.
		values := lo.Values[string, int](m)
		assert.Equal(t, []int{1, 2}, values)
	})

	// Returns the value of the given key or the fallback value if the key is not present.
	t.Run("ValueOr", func(t *testing.T) {
		value := lo.ValueOr[string, int](map[string]int{"foo": 1, "bar": 2}, "foo", 42)
		assert.Equal(t, 1, value)

		value = lo.ValueOr[string, int](map[string]int{"foo": 1, "bar": 2}, "baz", 42)
		assert.Equal(t, 42, value)
	})

	t.Run("MapToSlice", func(t *testing.T) {
		m := map[int]int64{1: 4, 2: 5, 3: 6}
		s := lo.MapToSlice(m, func(k int, v int64) string {
			return fmt.Sprintf("%d_%d", k, v)
		})
		assert.Equal(t, []string{"1_4", "2_5", "3_6"}, s)
	})

	t.Run("Range / RangeFrom / RangeWithSteps", func(t *testing.T) {
		assert.Equal(t, []int{0, 1, 2, 3}, lo.Range(4))
		assert.Equal(t, []int{0, -1, -2, -3}, lo.Range(-4))

		assert.Equal(t, []int{1, 2, 3, 4, 5}, lo.RangeFrom(1, 5))
		assert.Equal(t, []float64{1.0, 2.0, 3.0, 4.0, 5.0}, lo.RangeFrom[float64](1.0, 5))

		assert.Equal(t, []int{0, 5, 10, 15}, lo.RangeWithSteps(0, 20, 5))
		assert.Equal(t, []float32{-1.0, -2.0, -3.0}, lo.RangeWithSteps[float32](-1.0, -4.0, -1.0))
		assert.Equal(t, []int{}, lo.RangeWithSteps(1, 4, -1))
		assert.Equal(t, []int{}, lo.Range(0))
	})
}
