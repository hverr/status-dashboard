package widgets

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGatherCores(t *testing.T) {
	_, err := GatherCores()
	assert.Nil(t, err, "Did not expect error:", err)
}

func TestGatherLoadAverage(t *testing.T) {
	_, _, _, err := GatherLoadAverage()
	assert.Nil(t, err, "Did not expect error:", err)
}
