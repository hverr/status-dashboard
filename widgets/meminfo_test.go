package widgets

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGatherMeminfo(t *testing.T) {
	m, err := GatherMeminfo()
	require.Nil(t, err, "Did not expect error: %v", err)

	assert.NotEqual(t, 0, m.MainTotalKB, "MainTotalKB is zero")
	assert.NotEqual(t, 0, m.MainFreeKB, "MainFreeKB is zero")
	assert.NotEqual(t, 0, m.MainBuffersKB, "MainBuffersKB is zero")
	assert.NotEqual(t, 0, m.MainCachedKB, "MainCachedKB is zero")
	assert.NotEqual(t, 0, m.MainUsedKB, "MainUsedKB is zero")
}
