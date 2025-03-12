package utils

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func setup() []string {
	return []string{"1", "2", "3"}
}

func TestIndex(t *testing.T) {
	strList := setup()

	index := Index(strList, "2")
	assert.Equal(t, 1, index)
}

func TestInclude(t *testing.T) {
	strList := setup()
	assert.True(t, Include(strList, "3"))
}

func TestArrayToString(t *testing.T) {
	assert.Equal(t, "2,6,7,8", ArrayToString([]string{"2", "6", "7", "8"}, ","))
}
