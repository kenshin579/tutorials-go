package go_errors

import (
	"errors"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

var (
	ErrNotFound     = errors.New("err not found")
	ErrInvalidQuery = errors.New("err invalid query")
)

func TestIs(t *testing.T) {
	err := ErrNotFound
	assert.True(t, errors.Is(err, ErrNotFound))
	assert.False(t, errors.Is(err, ErrInvalidQuery))

	//%w의 경우에는 errors.Is로 같은 오류인지 확인할 수 있다
	e := fmt.Errorf("adding more context: %w", ErrNotFound)
	assert.True(t, errors.Is(e, ErrNotFound))
}

type QueryError struct {
	Query string
	Err   error
}

func (q *QueryError) Error() string {
	return fmt.Sprintf("%s: %s", q.Query, q.Err)
}

/*
errors.As()는 interface나 error를 구현한 타입인 경우에 사용할 수 있음
*/
func TestAs(t *testing.T) {
	q := &QueryError{
		Query: "query 1",
		Err:   ErrInvalidQuery,
	}

	var err *QueryError
	result := errors.As(q, &err)
	assert.True(t, result)
	assert.Equal(t, ErrInvalidQuery, err.Err)
}
