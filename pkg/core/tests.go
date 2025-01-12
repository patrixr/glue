package core

import (
	"fmt"
	"strings"

	"github.com/patrixr/q"
)

type TestFunc func()

type Testable interface {
	RegisterTest(name string, t TestFunc)
	Test()
	Testing() bool
	TestResults() []TestResult
}

type UnitTest struct {
	Id   string
	Name string
	Fn   TestFunc
}

type TestResult struct {
	Test  UnitTest
	Error error
}

type TestSuite struct {
	tests   []UnitTest
	testing bool
	results []TestResult
}

func NewTestSuite() Testable {
	return &TestSuite{
		testing: false,
		tests:   []UnitTest{},
		results: []TestResult{},
	}
}

func (ts *TestSuite) Testing() bool {
	return ts.testing
}

func (ts *TestSuite) RegisterTest(name string, fn TestFunc) {
	ts.tests = append(ts.tests, UnitTest{
		Name: name,
		Id:   strings.ReplaceAll(name, " ", "_"),
		Fn:   fn,
	})
}

func (ts *TestSuite) TestResults() []TestResult {
	return ts.results
}

func (ts *TestSuite) Test() {
	ts.testing = true

	defer func() {
		ts.testing = false
	}()

	ts.results = q.Map(ts.tests, test)
}

func test(t UnitTest) (result TestResult) {
	result.Test = t

	defer func() {
		if r := recover(); r != nil {
			result.Error = fmt.Errorf("%v", r)
		}
	}()

	t.Fn()

	return
}
