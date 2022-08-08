package runtimeutil

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCallerFuncPos(t *testing.T) {
	assert.True(t, strings.Contains(moo(), "TestCallerFuncPos:"))
}

func TestCaller(t *testing.T) {
	assert.True(t, strings.Contains(foo(1), "TestCaller"))
	assert.True(t, strings.Contains(bar(2), "TestCaller"))

	file, line, funcName := Caller(0)
	assert.True(t, strings.Contains(file, "caller_test.go"))
	assert.True(t, strings.Contains(funcName, ".TestCaller"))
	assert.True(t, line != 0)
}

func TestCallerFuncName(t *testing.T) {
	assert.True(t, strings.Contains(baz(), "TestCallerFuncName"))
}

func foo(skip int) string {
	_, _, s := Caller(skip)
	return s
}

func bar(skip int) string {
	return foo(skip)
}

func moo() string {
	return CallerFuncPos(1)
}

func baz() string {
	return CallerFuncName(1)
}
