package passwordhash

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestHashAndVerify(t *testing.T) {
	hash, err := Hash("password")
	require.NoError(t, err)
	assert.NotEqual(t, "password", hash)
	assert.True(t, Verify("password", hash))
	assert.False(t, Verify("wrong", hash))
}
