package util

import (
	"fmt"
	"testing"
)

type TestCaseInfo struct {
	T             *testing.T
	TestFnName    string
	TestCaseIndex int
}

func (tci *TestCaseInfo) Helper() {
	tci.T.Helper()
}

func (tci *TestCaseInfo) Fatal(expected, actual string) {
	tci.T.Fatalf("fn = [%s], index = [%d] expected [%s], actual [%s]", tci.TestFnName, tci.TestCaseIndex, expected, actual)
}

func (tci *TestCaseInfo) Fatalf(pattern string, fields ...interface{}) {
	tci.T.Fatalf("fn = [%s], index = [%d] : %s", tci.TestFnName, tci.TestCaseIndex, fmt.Sprintf(pattern, fields...))
}
