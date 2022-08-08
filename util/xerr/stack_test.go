package xerr

import (
	"errors"
	"fmt"
	"runtime"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/zhaolion/gostack/util/runtimeutil"
)

func TestWithStackRepeat(t *testing.T) {
	err := WithStack(errors.New("test error"))
	err = wrapErr(err)
	stackErr, ok := err.(*withStack)
	require.True(t, ok)
	assert.Equal(t, runtimeutil.CallerFuncName(0), funcStack(stackErr.stack))
}

func TestWrapStack(t *testing.T) {
	err := wrapStack(nil, 0)
	require.Nil(t, err, "wrap nil")

	err = wrapStack(errors.New("test error"), 0)
	require.NotNil(t, err)
	assert.Equal(t, "test error", err.Error(), "wrap normal error")
	assert.Equal(t, runtimeutil.CallerFuncName(0), funcStack(err.(*withStack).stack), "wrap normal error")
}

func wrapErr(err error) error {
	return wrapStack(err, 1)
}

func funcStack(s *stack) string {
	frame, _ := runtime.CallersFrames(*s).Next()
	return fmt.Sprintf("pkg.xerr.%s", funcname(frame.Function))
}
