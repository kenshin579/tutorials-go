package go_mockery

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type SomeStruct struct {
	Name string
}

type MockSomeStruct struct {
	mock.Mock
}

func (m *MockSomeStruct) WithName(name string) bool {
	args := m.Called(name)
	return args.Bool(0)
}

// MatchedBy() 함수는 실행되는 함수의 인자 값을 확인하여 user-defined 함수를 정의해서 사용할수 있다
func Test_MatchedBy(t *testing.T) {
	mockStruct := new(MockSomeStruct)

	// Expect the WithName method to be called with a SomeStruct instance that has a Name field
	// that starts with "John"
	mockStruct.On("WithName", mock.MatchedBy(func(arg string) bool {
		return strings.HasPrefix(arg, "John")
	})).Return(true)

	// Call the WithName method with a string that starts with "John"
	result := mockStruct.WithName("John Smith")
	assert.True(t, result)

	// Assert that the expectation was met
	mockStruct.AssertExpectations(t)
}

type MyInterface interface {
	DoSomething(arg1 string, arg2 int) bool
}

type MyMock struct {
	mock.Mock
}

func (m *MyMock) DoSomething(arg1 string, arg2 int) bool {
	args := m.Called(arg1, arg2)
	return args.Bool(0)
}

func Test_MatchedBy2(t *testing.T) {
	myMock := new(MyMock)

	// Set up an expectation that the DoSomething method will be called with a string argument
	// that starts with "foo" and an int argument that is greater than 10.
	myMock.On("DoSomething", mock.MatchedBy(func(arg1 string) bool {
		return strings.HasPrefix(arg1, "foo")
	}), mock.MatchedBy(func(arg2 int) bool {
		return arg2 > 10
	})).Return(true)

	// Call the function with an argument that matches the expectation.
	result := myMock.DoSomething("foo bar", 20)

	assert.True(t, result)

	// Assert that the expectation was met.
	myMock.AssertExpectations(t)
}
