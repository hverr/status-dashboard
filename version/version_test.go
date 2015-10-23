package version

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestIsDirty(t *testing.T) {
	Commit = "2041633c7e4b7f576dbfd513b40e97f58b5b2503"

	Dirty = "yes"
	assert.True(t, IsDirty())

	Dirty = ""
	assert.False(t, IsDirty())

	Dirty = "no"
	assert.False(t, IsDirty())

	Dirty = "qsdf"
	assert.False(t, IsDirty())
}

func TestHasVersionInformation(t *testing.T) {
	Commit = "2041633c7e4b7f576dbfd513b40e97f58b5b2503"
	assert.True(t, HasVersionInformation())

	Commit = ""
	assert.False(t, HasVersionInformation())
}

func TestPrintVersionInformation(t *testing.T) {
	{
		Commit = ""
		buf := bytes.NewBuffer(nil)
		flag := PrintVersionInformation(buf)
		assert.False(t, flag)
		assert.Contains(t, buf.String(), "No version information")
	}
	{
		Commit = "2041633c7e4b7f576dbfd513b40e97f58b5b2503"
		Dirty = "yes"
		buf := bytes.NewBuffer(nil)
		flag := PrintVersionInformation(buf)
		assert.True(t, flag)
		assert.Contains(t, buf.String(), Commit)
		assert.Contains(t, buf.String(), "dirty")
	}
	{
		Commit = "2041633c7e4b7f576dbfd513b40e97f58b5b2503"
		Dirty = ""
		buf := bytes.NewBuffer(nil)
		flag := PrintVersionInformation(buf)
		assert.True(t, flag)
		assert.Contains(t, buf.String(), Commit)
		assert.NotContains(t, buf.String(), "dirty")
	}
}
