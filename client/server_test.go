package client

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"net/url"
	"reflect"
	"testing"

	"github.com/hverr/status-dashboard/server"
	"github.com/hverr/status-dashboard/widgets"
	"github.com/jmcvetta/napping"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func newTestHttpServer(f func(w http.ResponseWriter, r *http.Request)) (*httptest.Server, *http.Client) {
	ts := httptest.NewServer(http.HandlerFunc(f))
	transport := &http.Transport{
		Proxy: func(req *http.Request) (*url.URL, error) {
			return url.Parse(ts.URL)
		},
	}
	client := &http.Client{
		Transport: transport,
	}

	return ts, client
}

func newTestSession(f func(w http.ResponseWriter, r *http.Request)) (*httptest.Server, napping.Session) {
	ts, client := newTestHttpServer(f)

	session := napping.Session{
		Client: client,
	}

	return ts, session
}

const (
	TestServerIdentifier = "webserver"
	TestServerAPI        = "http://localhost/api/"
	TestServerName       = "Web Server"
	TestServerSecret     = "testsecret"
)

func newTestServer(f func(w http.ResponseWriter, r *http.Request)) (*httptest.Server, *Server) {
	ts, session := newTestSession(f)

	s := &Server{
		Session: session,
		Configuration: Configuration{
			Identifier: TestServerIdentifier,
			API:        TestServerAPI,
			Name:       TestServerName,
			Secret:     TestServerSecret,
		},
	}

	return ts, s
}

func TestRegister(t *testing.T) {
	// Setup
	networkConfiguration, err := json.Marshal(map[string]interface{}{
		"interface": "lo",
	})
	require.Nil(t, err)

	availableWidgets := []server.WidgetRegistration{
		{
			Type:          "network",
			Configuration: networkConfiguration,
		},
		{
			Type:          "load",
			Configuration: []byte("null"),
		},
	}

	expectedPayload := server.ClientRegistration{
		Name:             TestServerName,
		Identifier:       TestServerIdentifier,
		AvailableWidgets: availableWidgets,
	}

	// Test
	ts, s := newTestServer(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, TestServerAPI+"/clients/"+TestServerIdentifier+"/register", r.URL.String())

		var payload server.ClientRegistration
		err := json.NewDecoder(r.Body).Decode(&payload)
		assert.Nil(t, err, "Could not decode body: %v", err)
		assert.Equal(t, TestServerName, payload.Name)
		assert.Equal(t, TestServerIdentifier, payload.Identifier)
		assert.True(
			t,
			reflect.DeepEqual(payload, expectedPayload),
			"Payload did not match: %v != %v (expected)",
			payload, availableWidgets,
		)

		w.WriteHeader(http.StatusOK)
	})
	defer ts.Close()

	err = s.Register(availableWidgets)
	assert.Nil(t, err, "Unexpected error: %v", err)
}

func TestRegisterWithError(t *testing.T) {
	// Simulate a JSON marshal error
	availableWidgets := []server.WidgetRegistration{
		{
			Type: "network",
		},
	}

	// Test
	ts, s := newTestServer(func(w http.ResponseWriter, r *http.Request) {
		assert.True(t, false, "Should not make server request.")
	})
	defer ts.Close()

	err := s.Register(availableWidgets)
	assert.Error(t, err, "Expected a JSON marshal error.")
}

func TestRegisterWithHttpError(t *testing.T) {
	// Setup
	availableWidgets := []server.WidgetRegistration{
		{
			Type:          "load",
			Configuration: []byte("null"),
		},
	}

	// Test
	ts, s := newTestServer(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	})
	defer ts.Close()

	err := s.Register(availableWidgets)
	assert.Error(t, err, "Expected a custom error.")
	assert.Contains(t, err.Error(), "Internal Server Error")
}

func TestGetRequestedWidgets(t *testing.T) {
	// Setup
	expectedResponse := RequestedWidgets{
		Widgets: []string{"load", "network_eth0"},
	}

	// Test
	ts, s := newTestServer(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, TestServerAPI+"/clients/"+TestServerIdentifier+"/requested_widgets", r.URL.String())
		err := json.NewEncoder(w).Encode(&expectedResponse)
		require.Nil(t, err, "Unexpected error: %v", err)
	})
	defer ts.Close()

	response, err := s.GetRequestedWidgets()
	assert.Nil(t, err, "Unexpected error: %v", err)
	assert.True(
		t,
		reflect.DeepEqual(expectedResponse, response),
		"Response doesn't match: %v == %v (expected)",
		response,
		expectedResponse,
	)
}

func TestGetRequestedWidgetsWithError(t *testing.T) {
	// Simulate malformed server response
	ts, s := newTestServer(func(w http.ResponseWriter, r *http.Request) {
		_, err := w.Write([]byte("malformed body"))
		require.Nil(t, err, "Unexpected error: %v", err)
	})
	defer ts.Close()

	// Test
	_, err := s.GetRequestedWidgets()
	assert.Error(t, err)
}

func TestGetRequestedWidgetsWithHttpError(t *testing.T) {
	// Test
	ts, s := newTestServer(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	})
	defer ts.Close()

	_, err := s.GetRequestedWidgets()
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "Internal Server Error")
}

func TestPostWidgetBulkUpdate(t *testing.T) {
	// Setup
	widget, err := json.Marshal(map[string]interface{}{
		"key": "value",
	})
	require.Nil(t, err)

	expectedUpdates := []widgets.BulkElement{
		{
			Type:       "load",
			Identifier: "load",
			Widget:     []byte("null"),
		},
		{
			Type:       "meminfo",
			Identifier: "meminfo",
			Widget:     widget,
		},
	}

	// Test
	ts, s := newTestServer(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, TestServerAPI+"/clients/"+TestServerIdentifier+"/bulk_update", r.URL.String())

		var payload []widgets.BulkElement
		err := json.NewDecoder(r.Body).Decode(&payload)
		assert.Nil(t, err, "Could not decode body: %v", err)
		assert.True(
			t,
			reflect.DeepEqual(expectedUpdates, payload),
			"Payload did not match: %v == %v (expected).",
			payload,
			expectedUpdates,
		)

		w.WriteHeader(http.StatusOK)
	})
	defer ts.Close()

	err = s.PostWidgetBulkUpdate(expectedUpdates)
	assert.Nil(t, err, "Unexpected error: %v", err)
}

func TestPostWidgetBulkUpdateWithError(t *testing.T) {
	// Simulate a JSON marshal error
	expectedUpdates := []widgets.BulkElement{
		{
			Type:       "load",
			Identifier: "load",
			Widget:     nil,
		},
	}

	// Test
	ts, s := newTestServer(func(w http.ResponseWriter, r *http.Request) {
		assert.True(t, false, "Should not make server request.")
	})
	defer ts.Close()

	err := s.PostWidgetBulkUpdate(expectedUpdates)
	assert.Error(t, err, "Expected a JSON marshal error.")
}

func TestPostWidgetBulkUpdateWithHttpError(t *testing.T) {
	// Setup
	expectedUpdates := []widgets.BulkElement{
		{
			Type:       "load",
			Identifier: "load",
			Widget:     []byte("null"),
		},
	}

	// Test
	ts, s := newTestServer(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	})
	defer ts.Close()

	err := s.PostWidgetBulkUpdate(expectedUpdates)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "Internal Server Error")
}
