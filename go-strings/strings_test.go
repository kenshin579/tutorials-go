package go_strings

import (
	"fmt"
	"math"
	"strconv"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_String_Search(t *testing.T) {
	//Search (Contains, Prefix/Suffix, Index)
	assert.True(t, strings.Contains("test", "st"))
	assert.True(t, strings.ContainsAny("test", "s"))
	assert.True(t, strings.HasPrefix("test", "te"))
	assert.True(t, strings.HasSuffix("test", "st"))
	assert.Equal(t, 2, strings.Count("test", "t"))
	assert.Equal(t, 1, strings.Index("test", "e"))
}

func Test_String_Replace(t *testing.T) {
	//Replace (Uppercase/Lowercase, Trim)
	assert.Equal(t, "f00", strings.Replace("foo", "o", "0", -1))
	assert.Equal(t, "test", strings.ToLower("TEST"))
	assert.Equal(t, "TEST", strings.ToUpper("test"))
	assert.Equal(t, "Test", strings.Trim(" Test  ", " "))
	assert.Equal(t, "Test", strings.TrimSpace(" Test  "))
	f := func(r rune) rune {
		return r + 1
	}
	assert.Equal(t, "bc", strings.Map(f, "ab"))
}

func Test_String_Split(t *testing.T) {
	//Split (split, fields)
	assert.Equal(t, []string{"a", "b", "c"}, strings.Split("a,b,c", ","))
	assert.Equal(t, []string{"t", "e", "s", "t"}, strings.Fields("t\t   e s t"))
}

func Test_String_Concatenate(t *testing.T) {
	//Concatenate (+, Sprintf, builder)
	assert.Equal(t, "hello world", "hello"+" world")
	assert.Equal(t, "data: 123", fmt.Sprintf("%s %d", "data:", 123))
	assert.Equal(t, "3.1416", fmt.Sprintf("%.4f", math.Pi))

	var b strings.Builder
	for i := 3; i >= 1; i-- {
		fmt.Fprintf(&b, "%d...", i)
	}
	b.WriteString("end")
	assert.Equal(t, "3...2...1...end", b.String())
}

func Test_String_Join(t *testing.T) {
	//Join (Join, Repeat)
	assert.Equal(t, "a-b", strings.Join([]string{"a", "b"}, "-"))
	assert.Equal(t, "AAAAA", strings.Repeat("A", 5))
}

func Test_String_Format(t *testing.T) {
	//Format, Convert (strconv)
	assert.Equal(t, "23", strconv.Itoa(23))
	assert.Equal(t, "ff", strconv.FormatInt(255, 16))

	intValue, _ := strconv.Atoi("23")
	assert.Equal(t, 23, intValue)
}

func Test_String_Equal(t *testing.T) {
	assert.True(t, strings.EqualFold("QUIBI", "quibi"))
}
