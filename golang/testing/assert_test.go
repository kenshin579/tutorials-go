package go_testing

import (
	"fmt"
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/thoas/go-funk"
)

func Test_IsSlice(t *testing.T) {
	strSlice := []string{"a", "b", "c"}
	assert.True(t, isSlice(strSlice))

	assert.False(t, isSlice("str"))
}

func Test_Assert_Slices_Not_Inorder(t *testing.T) {
	expected := []string{"a", "b", "c"}
	actual := []string{"c", "a", "b"}

	// Check if actual slice contains all elements of expected slice
	if !funk.Contains(funk.Uniq(expected), func(x string) bool {
		return funk.Contains(actual, x)
	}) {
		t.Errorf("Slices do not contain the same elements")
	}

}

// https://stackoverflow.com/questions/62750483/check-if-interface-is-slice-of-something-interface
func isSlice(slice interface{}) bool {
	typ := reflect.TypeOf(slice)
	if typ.Kind() == reflect.Slice {
		elemType := typ.Elem()
		fmt.Println("slice of", elemType.Name())
		return true
	}
	return false
}
