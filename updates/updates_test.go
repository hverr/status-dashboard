package updates

import (
	"io"
	"runtime"
	"testing"

	"github.com/hverr/status-dashboard/version"
	"github.com/stretchr/testify/assert"
)

func TestUpdater(t *testing.T) {
	version.Commit = "da39a3ee5e6b4b0d3255bfef95601890afd80709"

	u := appUpdater("dashboard-client")
	assert.NotNil(t, u.App)
	assert.Equal(t, version.Commit, u.CurrentReleaseIdentifier)
	assert.NotNil(t, u.WriterForAsset)
}

func TestAssetWriter(t *testing.T) {
	{
		a := &testAsset{name: "not an asset"}
		f := assetWriter("dashboard-client")
		u, err := f(a)
		assert.Nil(t, u)
		assert.Nil(t, err)
	}
	{
		a := &testAsset{name: assetName("dashboard-client")}
		f := assetWriter("dashboard-client")
		u, err := f(a)
		assert.NotNil(t, u)
		assert.Nil(t, err)
	}
}

func TestAssetName(t *testing.T) {
	n := assetName("dashboard-client")
	assert.Contains(t, "dashboard-client_"+runtime.GOOS+"_"+runtime.GOARCH, n)
}

type testAsset struct{ name string }

func (a *testAsset) Name() string          { return a.name }
func (a *testAsset) Write(io.Writer) error { return nil }
