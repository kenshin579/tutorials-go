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

func (q *QueryError) Unwrap() error {
	return q.Err
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

func Test_Wrap_Unwrap(t *testing.T) {
	q1 := &QueryError{
		Query: "query 1",
		Err:   errors.New("err 1"),
	}

	q2 := fmt.Errorf("q2: %w", q1)
	q3 := fmt.Errorf("q3 : %w", q2)

	fmt.Println(q2)
	fmt.Println(q3)

	//Unwrap
	fmt.Println(errors.Unwrap(q3))
	fmt.Println(errors.Unwrap(q2))
	fmt.Println(errors.Unwrap(q1))
}

const badInput = "abc"

var ErrBadInput = errors.New("bad input")

func validateInputForIs(input string) error {
	if input == badInput {
		return fmt.Errorf("validateInput: %w", ErrBadInput) //ErrBadInput error는 fmt.Errorf()의해 생성된 오류에 wrapped됨
	}
	return nil
}

/*
Is() 함수는 주어진 오류가 다른 특정 오류와 일치하는지 확인하려는 경우에 사용한다
*/
func TestError_Is(t *testing.T) {
	input := badInput

	err := validateInputForIs(input)

	assert.True(t, errors.Is(err, ErrBadInput)) //bad input error
}

type BadInputError struct {
	input string
}

func (e *BadInputError) Error() string {
	return fmt.Sprintf("bad input: %s", e.input)
}

func validateInputForAs(input string) error {
	if input == badInput {
		return fmt.Errorf("validateInput: %w", &BadInputError{input: input})
	}
	return nil
}

func TestError_As(t *testing.T) {
	input := badInput

	err := validateInputForAs(input)
	var badInputErr *BadInputError

	if errors.As(err, &badInputErr) {
		fmt.Printf("bad input error occured: %#v\n", badInputErr)
	}
}
