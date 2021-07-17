package go_testing

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/suite"
)

//https://brunoscheufler.com/blog/2020-04-12-building-go-test-suites-using-testify
type ExampleTestSuite struct {
	suite.Suite
}

func (suite *ExampleTestSuite) BeforeTest(suiteName, testName string) {
	fmt.Println(suiteName, testName)
	fmt.Println("BeforeTest - before test")
}

func (suite *ExampleTestSuite) AfterTest(suiteName, testName string) {
	fmt.Println(suiteName, testName)
	fmt.Println("AfterTest - after test")
}

// This will run before before the tests in the suite are run
func (suite *ExampleTestSuite) SetupSuite() {
	fmt.Println("SetupSuite - before once")
}

// This will run before each test in the suite
func (suite *ExampleTestSuite) SetupTest() {
	fmt.Println("SetupTest - setup test")
}

// This is an example test that will always succeed
func (suite *ExampleTestSuite) TestExample1() {
	fmt.Println("test1")
}

func (suite *ExampleTestSuite) TestExample2() {
	fmt.Println("test2")
}

// We need this function to kick off the test suite, otherwise
// "go test" won't know about our tests
func TestExampleTestSuite(t *testing.T) {
	if isSkip() {
		t.Skip("skipping")
	}
	suite.Run(t, new(ExampleTestSuite))
}

func isSkip() bool {
	return false
}
