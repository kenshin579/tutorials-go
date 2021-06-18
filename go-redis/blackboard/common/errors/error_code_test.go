package errors

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestReformatMessageWithParams(t *testing.T) {
	err := NewBlackBoardError(http.StatusNotFound, 100, "{0} + {1} = {2}").
		WithParams("1", "2", "3")

	assert.Equal(t, "1 + 2 = 3", err.Error())
}
