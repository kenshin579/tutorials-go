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

	// Output:
	// 2020-01-04 12:26:00 +0000 UTC
	// 2020 January 4 12 26 0
	// Saturday 1578140760 0
}

func ExampleDate_AddDate() {
	t := time.Date(2021, time.Month(4), 10, 0, 0, 0, 0, time.UTC)
	fmt.Println(t.AddDate(0, 1, 0))

	// Output: 2021-05-10 00:00:00 +0000 UTC
}

func TestBefore_After_Equal(t *testing.T) {
	t1 := time.Date(2021, time.Month(4), 10, 0, 0, 0, 0, time.UTC)
	t2 := time.Date(2021, time.Month(4), 11, 0, 0, 0, 0, time.UTC)

	// t1 < t2 (now)
	assert.True(t, t1.Before(t2))
	assert.True(t, t2.After(t1))
	assert.True(t, t1.Equal(t1))
}

func Test_isExpired(t *testing.T) {
	expiration := time.Now().Add(3 * time.Second)
	assert.False(t, isExpired(expiration))

	time.Sleep(4 * time.Second)
	assert.True(t, isExpired(expiration))
}

func isExpired(expirationTime time.Time) bool {
	return time.Now().After(expirationTime)
}

func Test_Since(t *testing.T) {
	utc := time.Now()

	assert.False(t, isReadyToUpdate(utc.Add(-time.Minute*3)))
	assert.True(t, isReadyToUpdate(utc.Add(-time.Minute*100)))
}

func isReadyToUpdate(updatedAt time.Time) bool {
	duration := time.Since(updatedAt)

	// Check if the duration is greater than or equal to 5 minutes
	return duration >= 5*time.Minute
}
