package widgets

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGatherUptime(t *testing.T) {
	_, _, _, _, err := GatherUptime()
	assert.Nil(t, err, "Did not expect error: %v", err)
}
