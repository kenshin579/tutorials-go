package go_defaults

import (
	"testing"

	"github.com/creasty/defaults"
	"github.com/stretchr/testify/assert"
)

type Sample struct {
	Name string `default:"John Smith"`
	Age  int    `default:"27"`
}

func Test(t *testing.T) {
	sample := Sample{}
	defaults.Set(&sample)
	assert.Equal(t, 27, sample.Age)
}
