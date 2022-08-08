package xerr

import (
	"io"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestCustomError(t *testing.T) {
	_, got := IsCustomError(Custom("test"))
	assert.True(t, got)

	_, got = IsCustomError(Customf("test %d", time.Now().Unix()))
	assert.True(t, got)

	_, got = IsCustomError(NewNotFoundError(&dummy{}, "test1"))
	assert.True(t, got)

	_, got = IsCustomError(io.EOF)
	assert.False(t, got)

	assert.Error(t, IgnoreDuplicateEntryError(io.EOF))
	assert.NoError(t, IgnoreDuplicateEntryError(Custom("Duplicate entry UQE-x-y"), "a"))
	assert.NoError(t, IgnoreDuplicateEntryError(nil))
}
