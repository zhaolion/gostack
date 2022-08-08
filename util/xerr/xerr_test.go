package xerr

import (
	"io"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMyWrap(t *testing.T) {
	assert.Error(t, foo())

	assert.Error(t, Wrap(foo(), "wrap_foo"))
	assert.Error(t, WrapWithLog(foo(), "wrap_foo_log"))

	assert.Error(t, Wrap(io.EOF))
	assert.Error(t, WrapWithLog(io.EOF))

	assert.Error(t, WithMessage(foo()))
	assert.Error(t, WithMessage(foo(), "with message"))
}

func foo() error {
	return Wrap(io.EOF, "read file")
}
