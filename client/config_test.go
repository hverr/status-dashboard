package client

import (
	"bytes"
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestParseConfiguration(t *testing.T) {
	// Setup
	testConfig := map[string]interface{}{
		"api":        "http://localhost/api",
		"identifier": "webserver",
		"name":       "Web Server",
		"secret":     "supersecret",
		"widgets": map[string]interface{}{
			"load": make(map[string]interface{}),
		},
	}

	b, err := json.Marshal(&testConfig)
	require.Nil(t, err)
	buf := bytes.NewBuffer(b)
	require.Nil(t, err)

	// Test
	c := Configuration{}
	err = c.ParseConfiguration(buf)
	require.Nil(t, err, "Unexpected error: %v", err)

	assert.Equal(t, "http://localhost/api", c.API)
	assert.Equal(t, "webserver", c.Identifier)
	assert.Equal(t, "Web Server", c.Name)
	assert.Equal(t, "supersecret", c.Secret)
	assert.NotNil(t, c.Widgets)
}

func TestParseConfigurationWithInvalidJSON(t *testing.T) {
	// Test
	c := Configuration{}
	err := c.ParseConfiguration(bytes.NewBufferString("invalid json"))
	assert.Error(t, err)
}

func TestParseConfigurationWithMissingAPI(t *testing.T) {
	// Setup
	testConfig := map[string]interface{}{
		"identifier": "webserver",
		"name":       "Web Server",
		"secret":     "supersecret",
		"widgets": map[string]interface{}{
			"load": make(map[string]interface{}),
		},
	}

	b, err := json.Marshal(&testConfig)
	require.Nil(t, err)
	buf := bytes.NewBuffer(b)
	require.Nil(t, err)

	// Test
	c := Configuration{}
	err = c.ParseConfiguration(buf)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "API")
}

func TestParseConfigurationWithMissingWidgets(t *testing.T) {
	// Setup
	testConfig := map[string]interface{}{
		"api":        "http://localhost/api",
		"identifier": "webserver",
		"name":       "Web Server",
		"secret":     "supersecret",
	}

	b, err := json.Marshal(&testConfig)
	require.Nil(t, err)
	buf := bytes.NewBuffer(b)
	require.Nil(t, err)

	// Test
	c := Configuration{}
	err = c.ParseConfiguration(buf)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "widgets")
}

func TestParseConfigurationWithInvalidWidget(t *testing.T) {
	// Setup
	testConfig := map[string]interface{}{
		"api":        "http://localhost/api",
		"identifier": "webserver",
		"name":       "Web Server",
		"secret":     "supersecret",
		"widgets": map[string]interface{}{
			"invalid widget": make(map[string]interface{}),
		},
	}

	b, err := json.Marshal(&testConfig)
	require.Nil(t, err)
	buf := bytes.NewBuffer(b)
	require.Nil(t, err)

	// Test
	c := Configuration{}
	err = c.ParseConfiguration(buf)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "Unsupported widget")
}
