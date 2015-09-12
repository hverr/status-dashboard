package widgets

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewTCPInfo(t *testing.T) {
	{
		tcp6 := " 1: 00000000000000000000000001000000:0277 00000000000000000000000000000000:0000 0A 00000000:00000000 00:00000000 00000000     0        0 13365 1 0000000000000000 100 0 0 10 0"
		i, err := NewTCPInfo(tcp6)
		assert.Nil(t, err, "Could not parse TCP6 info string.")
		if i != nil {
			assert.EqualValues(t, 1, i.Entry, "Entry did not match.")
			assert.EqualValues(t, "::1", i.LocalIP.String(), "Local IP did not match.")
			assert.EqualValues(t, 631, i.LocalPort, "Local port did not match.")
			assert.EqualValues(t, "::", i.RemoteIP.String(), "Remote IP did not match.")
			assert.EqualValues(t, 0, i.RemotePort, "Remote port did not match.")
			assert.EqualValues(t, 10, i.State, "Remote state did not match.")
		}
	}

	{
		tcp6 := "85: 0000000000000000FFFF00001401A8C0:008B 0000000000000000FFFF00000228C19D:B6A5 01 00000000:00000000 02:000A8DC1 00000000    33        0 31881895 2 c0423280 22 4 28 10 189"
		i, err := NewTCPInfo(tcp6)
		assert.Nil(t, err, "Could not parse TCP6 info string.")
		if i != nil {
			assert.EqualValues(t, 85, i.Entry, "Entry did not match.")
			assert.EqualValues(t, "192.168.1.20", i.LocalIP.String(), "Local IP did not match.")
			assert.EqualValues(t, 139, i.LocalPort, "Local port did not match.")
			assert.EqualValues(t, "157.193.40.2", i.RemoteIP.String(), "Remote IP did not match.")
			assert.EqualValues(t, 46757, i.RemotePort, "Remote port did not match.")
			assert.EqualValues(t, 1, i.State, "Remote state did not match.")
		}
	}
	{
		tcp4 := "2: 0F02000A:C85E 0101144E:0006 01 00000000:00000000 02:0007BBC6 00000000  1000        0 26093 2 0000000000000000 20 4 26 10 -1"
		i, err := NewTCPInfo(tcp4)
		assert.Nil(t, err, "Could not parse TCP4 info string.")
		if i != nil {
			assert.EqualValues(t, 2, i.Entry, "Entry did not match.")
			assert.EqualValues(t, "10.0.2.15", i.LocalIP.String(), "Local IP did not match.")
			assert.EqualValues(t, 51294, i.LocalPort, "Local port did not match.")
			assert.EqualValues(t, "78.20.1.1", i.RemoteIP.String(), "Remote IP did not match.")
			assert.EqualValues(t, 6, i.RemotePort, "Remote port did not match.")
			assert.EqualValues(t, 1, i.State, "Remote state did not match.")
		}
	}
}
