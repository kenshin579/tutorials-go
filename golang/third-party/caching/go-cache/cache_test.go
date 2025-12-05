package main

import (
	"fmt"
	"testing"
	"time"

	"github.com/kenshin579/tutorials-go/golang/data-structure/sort/model"
	"github.com/patrickmn/go-cache"
	"github.com/stretchr/testify/assert"
)

func Test_cache(t *testing.T) {
	// Create a cache with a default expiration time of 5 minutes, and which
	// purges expired items every 10 minutes
	c := cache.New(1*time.Second, 10*time.Second)

	// Set the value of the key "foo" to "bar", with the default expiration time
	c.Set("test1", "value1", cache.DefaultExpiration)
	c.Set("test2", "value2", 2*time.Second)
	c.Set("test3", "value3", 5*time.Second)
	c.Set("test4", "value4", cache.NoExpiration)
	c.Set("test5", "value5", cache.NoExpiration)
	c.Set("test6", model.Employee{Name: "frank", Age: 20}, cache.NoExpiration)

	sleepInSec(12)

	// Set the value of the key "baz" to 42, with no expiration time
	// (the item won't be removed until it is re-set, or removed using
	c.Delete("test4")

	// Get the string associated with the key "foo" from the cache
	_, found := c.Get("test2")
	assert.False(t, found)

	value5, found := c.Get("test5")
	assert.True(t, found)
	assert.Equal(t, "value5", value5)

	value6, found := c.Get("test6")
	assert.True(t, found)
	employee := value6.(model.Employee)
	assert.Equal(t, employee.Name, "frank")
	assert.Equal(t, employee.Age, 20)
}

func sleepInSec(sleepTime int) {
	// sleepInSec for a few seconds and print dots
	for i := 0; i < sleepTime; i++ {
		time.Sleep(1 * time.Second)
		fmt.Print(".")
	}

	fmt.Print("done!")
}
