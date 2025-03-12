package go_testing

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/suite"
	"gopkg.in/go-playground/assert.v1"
)

// https://brunoscheufler.com/blog/2020-04-12-building-go-test-suites-using-testify
type ExampleTestSuite struct {
	suite.Suite
	TestValue int
}

// suite에서 전체 테스트 실행 전에 한번만 실행된다
func (ets *ExampleTestSuite) SetupSuite() {
	fmt.Println("SetupSuite :: run once")
}

// suite에서 각 테스트 실행 전에 실행된다
func (ets *ExampleTestSuite) SetupTest() {
	fmt.Println("SetupTest :: run setup test")
	ets.TestValue = 5
}

// 테스트가 실행하기 전에 suiteName testName 인자를 받아 실행하는 함수이다
func (ets *ExampleTestSuite) BeforeTest(suiteName, testName string) {
	fmt.Printf("BeforeTest :: run before test - suiteName:%s testName: %s\n", suiteName, testName)
}

// TEST ----------------------------------------
func (ets *ExampleTestSuite) TestExample1() {
	fmt.Println("TestExample1")
	ets.Equal(ets.TestValue, 5)
	assert.Equal(ets.T(), ets.TestValue, 5)
}

func (ets *ExampleTestSuite) TestExample2() {
	fmt.Println("TestExample2")
}

//TEST ----------------------------------------

// 테스트가 실행후에 suiteName testName 인자를 받아 실행하는 함수이다
func (ets *ExampleTestSuite) AfterTest(suiteName, testName string) {
	fmt.Printf("AfterTest :: suiteName:%s testName: %s\n", suiteName, testName)
}

// suite에서 각 테스트 실행후에 실행된다
func (ets *ExampleTestSuite) TearDownTest() {
	fmt.Println("TearDownTest :: run before after test")
}

// suite에서 모든 테스트가 실행된 후에 실행된다
func (ets *ExampleTestSuite) TearDownSuite() {
	fmt.Println("TearDownSuite :: run once")
}

func TestExampleTestSuite(t *testing.T) {
	//특정 조건인 경우에 전체 테스트를 실행하지 않을 수 있다
	if isSkip() {
		t.Skip("skipping")
	}
	suite.Run(t, new(ExampleTestSuite))
}

func isSkip() bool {
	return false
}
