package xerr

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_withStack_Unwrap(t *testing.T) {
	err := &CustomError{msg: "test error"}
	stackErr := WithStack(err)

	assert.NotEqual(t, stackErr, err)
	assert.Equal(t, errors.Unwrap(stackErr), err)

	assert.True(t, errors.Is(stackErr, err))
	assert.ErrorIs(t, stackErr, err)
}
