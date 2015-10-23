package widgets

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLoadWidget(t *testing.T) {
	w := &LoadWidget{"0.00", "0.01", "0.05", 4}
	assert.NotEqual(t, "", w.Name())
	assert.NotEqual(t, "", w.Type())
	assert.Equal(t, w.Type(), w.Identifier())
	assert.True(t, w.HasData())
	assert.Nil(t, w.Configure(nil))
	assert.Nil(t, w.Configuration())
	assert.Nil(t, w.Start())
	assert.Nil(t, w.Update())
}

func TestGatherCores(t *testing.T) {
	_, err := GatherCores()
	assert.Nil(t, err, "Did not expect error: %v", err)
}

func TestGatherLoadAverage(t *testing.T) {
	_, _, _, err := GatherLoadAverage()
	assert.Nil(t, err, "Did not expect error: %v", err)
}
