package main

import (
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func ExampleTimeDate() {
	t := time.Date(2020, time.January, 4, 12, 26, 0, 0, time.UTC)
	fmt.Println(t)
	fmt.Println(t.Year(), t.Month(), t.Day(), t.Hour(), t.Minute(), t.Second())
	fmt.Println(t.Weekday(), t.Unix(), t.Nanosecond())

	//Output:
	//2020-01-04 12:26:00 +0000 UTC
	//2020 January 4 12 26 0
	//Saturday 1578140760 0
}

func ExampleDate_AddDate() {
	t := time.Date(2021, time.Month(4), 10, 0, 0, 0, 0, time.UTC)
	fmt.Println(t.AddDate(0, 1, 0))

	//Output: 2021-05-10 00:00:00 +0000 UTC
}

func TestBefore_After_Equal(t *testing.T) {
	t1 := time.Date(2021, time.Month(4), 10, 0, 0, 0, 0, time.UTC)
	t2 := time.Date(2021, time.Month(4), 11, 0, 0, 0, 0, time.UTC)

	assert.True(t, t1.Before(t2))
	assert.True(t, t2.After(t1))
	assert.True(t, t1.Equal(t1))
}
