package go_testing

import (
	"fmt"
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_IsSlice(t *testing.T) {
	strSlice := []string{"a", "b", "c"}
	assert.True(t, isSlice(strSlice))

	assert.False(t, isSlice("str"))
}

//https://stackoverflow.com/questions/62750483/check-if-interface-is-slice-of-something-interface
func isSlice(slice interface{}) bool {
	typ := reflect.TypeOf(slice)
	if typ.Kind() == reflect.Slice {
		elemType := typ.Elem()
		fmt.Println("slice of", elemType.Name())
		return true
	}
	return false
}
