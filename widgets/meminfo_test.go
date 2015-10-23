package widgets

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestMeminfoWidget(t *testing.T) {
	w := &MeminfoWidget{2048, 1024}
	assert.NotEqual(t, "", w.Name())
	assert.NotEqual(t, "", w.Type())
	assert.Equal(t, w.Type(), w.Identifier())
	assert.True(t, w.HasData())
	assert.Nil(t, w.Configure(nil))
	assert.Nil(t, w.Configuration())
	assert.Nil(t, w.Start())
	assert.Nil(t, w.Update())
}
func TestGatherMeminfo(t *testing.T) {
	m, err := GatherMeminfo()
	require.Nil(t, err, "Did not expect error: %v", err)

	assert.NotEqual(t, 0, m.MainTotalKB, "MainTotalKB is zero")
	assert.NotEqual(t, 0, m.MainFreeKB, "MainFreeKB is zero")
	assert.NotEqual(t, 0, m.MainBuffersKB, "MainBuffersKB is zero")
	assert.NotEqual(t, 0, m.MainCachedKB, "MainCachedKB is zero")
	assert.NotEqual(t, 0, m.MainUsedKB, "MainUsedKB is zero")
}
