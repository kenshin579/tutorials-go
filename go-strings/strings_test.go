package go_strings

import (
	"fmt"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestStrings(t *testing.T) {
	assert.True(t, strings.Contains("test", "st"))
	assert.True(t, strings.ContainsAny("test", "s"))
	assert.True(t, strings.HasPrefix("test", "te"))
	assert.True(t, strings.HasSuffix("test", "st"))

	assert.Equal(t, 2, strings.Count("test", "t"))
	assert.Equal(t, 1, strings.Index("test", "e"))
	assert.Equal(t, "a-b", strings.Join([]string{"a", "b"}, "-"))
	assert.Equal(t, "AAAAA", strings.Repeat("A", 5))
	assert.Equal(t, "f00", strings.Replace("foo", "o", "0", -1))
	assert.Equal(t, []string{"a", "b", "c"}, strings.Split("a,b,c", ","))

	assert.Equal(t, "test", strings.ToLower("TEST"))
	assert.Equal(t, "TEST", strings.ToUpper("test"))

	assert.Equal(t, []string{"t", "e", "s", "t"}, strings.Fields("t e s t"))

	assert.Equal(t, "Test", strings.Trim(" Test  ", " "))
	assert.Equal(t, "Test", strings.TrimSpace(" Test  "))

	assert.Equal(t, "hello world", "hello"+" world")
	assert.Equal(t, "hello world", fmt.Sprintf("%s %s", "hello", "world"))

	var b strings.Builder
	for i := 3; i >= 1; i-- {
		fmt.Fprintf(&b, "%d...", i)
	}
	b.WriteString("end")
	assert.Equal(t, "3...2...1...end", b.String())

}
