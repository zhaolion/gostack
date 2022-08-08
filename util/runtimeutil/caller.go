package runtimeutil

import (
	"fmt"
	"runtime"
	"strings"
)

// CallerFuncPos get caller's func name
func CallerFuncPos(skip int) string {
	_, line, name := Caller(skip + 1)
	return fmt.Sprintf("%s:%d", name, line)
}

// CallerFuncName get caller's func name
func CallerFuncName(skip int) string {
	_, _, name := Caller(skip + 1)
	arr := strings.Split(name, "/")
	if len(arr) >= 3 {
		return strings.Join(arr[2:], ".")
	}
	return strings.Join(arr, ".")
}

// Caller file, file line, function name
func Caller(skip int) (file string, line int, functionName string) {
	var (
		pc uintptr
		ok bool
	)
	pc, file, line, ok = runtime.Caller(skip + 1)
	if !ok {
		return
	}

	return file, line, runtime.FuncForPC(pc).Name()
}
