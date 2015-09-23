package widgets

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewNetworkDevice(t *testing.T) {
	{
		line := "  lo:   22104     236    0    0    0     0          0         0    22104     236    0    0    0     0       0          0"
		dev := NewNetworkDevice(line)
		assert.NotNil(t, dev, "Could not parse line")
		if dev != nil {
			assert.EqualValues(t, "lo", dev.Interface, "Interface doesn't match.")
			assert.EqualValues(t, 22104, dev.ReceivedBytes, "Received bytes don't match.")
			assert.EqualValues(t, 22104, dev.TransmittedBytes, "Transmitted bytes don't match.")
		}
	}
	{
		line := "eth0:  313772     676    0    0    0     0          0         0    75840     674    0    0    0     0       0          0"
		dev := NewNetworkDevice(line)
		assert.NotNil(t, dev, "Could not parse line")
		if dev != nil {
			assert.EqualValues(t, "eth0", dev.Interface, "Interface doesn't match.")
			assert.EqualValues(t, 313772, dev.ReceivedBytes, "Received bytes don't match.")
			assert.EqualValues(t, 75840, dev.TransmittedBytes, "Transmitted bytes don't match.")
		}
	}
}
