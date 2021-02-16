package user

import (
	mocks "github.com/kenshin579/tutorials-go/go-mockery/do_user/mocks/doer"
	"github.com/stretchr/testify/mock"
	"strings"
	"testing"
)

func TestUserWithTestifyMock(t *testing.T) {
	mockDoer := &mocks.Doer{}

	testUser := &User{Doer: mockDoer}

	//given
	mockDoer.On("Do", 1, "abc").Return(nil).Once()

	//when
	testUser.Use()

	//then
	mockDoer.AssertExpectations(t) //todo: 이거는 뭔가?
}

//testify에서는 메서드 실행시 On()으로 mocking해야 하는 인자 값을 구체적으로 더 자세히 보여준다
func TestUser_Testify_UnexpectedCall(t *testing.T) {
	mockDoer := &mocks.Doer{}
	testUser := &User{Doer: mockDoer}

	testUser.Use()
	mockDoer.AssertExpectations(t)
}

//testify에서는 기대하는 인자값이 아닌 내용더 상세히 알려준다
func TestUser_Testify_UnexpectedArgs(t *testing.T) {
	mockDoer := &mocks.Doer{}
	testUser := &User{Doer: mockDoer}

	mockDoer.On("Do", 2, "def")

	// Calls mockDoer with (1, "abc")
	testUser.Use()
	mockDoer.AssertExpectations(t)
}

func TestUser_Testify_MissingCall(t *testing.T) {
	mockDoer := &mocks.Doer{}

	mockDoer.On("Do", 1, "abc").Return(nil)
	mockDoer.On("Do", 2, "def").Return(nil)

	mockDoer.AssertExpectations(t)
}

func TestArgument_Matchers_Anything(t *testing.T) {
	mockDoer := &mocks.Doer{}

	testUser := &User{Doer: mockDoer}

	//given
	mockDoer.On("Do", mock.Anything, mock.AnythingOfType("string")).Return(nil).Once()

	//when
	testUser.Use()

	//then
	mockDoer.AssertExpectations(t)
}

func TestArgument_Matchers_MatchedBy(t *testing.T) {
	mockDoer := &mocks.Doer{}

	testUser := &User{Doer: mockDoer}

	//given

	// Custom matcher
	mockDoer.On("Do", 1,
		mock.MatchedBy(func(x string) bool {
			return strings.HasPrefix(x, "abc")
		})).Return(nil).Once()

	//when
	testUser.Use()

	//then
	mockDoer.AssertExpectations(t)
}

func TestCallFrequency(t *testing.T) {
	mockDoer := &mocks.Doer{}

	testUser := &User{Doer: mockDoer}

	//given

	//mockDoer.On("Do", mock.Anything, mock.Anything).Return(nil).Twice()
	mockDoer.On("Do", mock.Anything, mock.Anything).Return(nil).Times(1)

	//when
	testUser.Use()

	//then
	mockDoer.AssertExpectations(t)
}
