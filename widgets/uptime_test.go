package widgets

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestUptimeWidget(t *testing.T) {
	w := &UptimeWidget{time.Now(), 10, 5, 3, 1}
	assert.NotEqual(t, "", w.Name())
	assert.NotEqual(t, "", w.Type())
	assert.Equal(t, w.Type(), w.Identifier())
	assert.True(t, w.HasData())
	assert.Nil(t, w.Configure(nil))
	assert.Nil(t, w.Configuration())
	assert.Nil(t, w.Start())
	assert.Nil(t, w.Update())
}
func TestGatherUptime(t *testing.T) {
	_, _, _, _, err := GatherUptime()
	assert.Nil(t, err, "Did not expect error: %v", err)
}
