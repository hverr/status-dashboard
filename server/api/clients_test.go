package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"reflect"
	"testing"
	"time"

	"github.com/hverr/status-dashboard/server"
	"github.com/hverr/status-dashboard/server/scheduler"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const (
	TestClientIdentifier = "webserver"
	TestClientName       = "Web Server"
	TestClientSecret     = "webserversecret"
)

func validApi() *API {
	conf := server.Configuration{
		UpdateInterval: 1,
		Clients: map[string]server.ClientConfiguration{
			"webserver": server.ClientConfiguration{"webserversecret"},
		},
		Users: map[string]string{},
	}
	serv := server.NewServer(conf)
	sched := scheduler.New(conf)
	userAuth := server.UserAuthenticator{conf}

	return &API{conf, serv, userAuth, sched}
}

func TestRegisterClient(t *testing.T) {
	// Setup
	testApi := validApi()
	reqUrl := fmt.Sprintf("/api/clients/%s/register", TestClientIdentifier)
	availableWidgets := []server.WidgetRegistration{
		{Type: "load", Configuration: []byte("null")},
	}
	m := map[string]interface{}{
		"name":             TestClientName,
		"identifier":       TestClientIdentifier,
		"availableWidgets": availableWidgets,
	}

	// Test all valid
	{
		resp := testApi.callJSON("POST", reqUrl, TestClientSecret, m)
		assert.Equal(t, http.StatusOK, resp.Code)
		assert.True(
			t,
			testApi.Scheduler.HasClient(TestClientIdentifier),
			"Expected succesful client registration in scheduler.",
		)

		registered := testApi.Server.AllRegisteredClients()
		found := false
		for _, r := range registered {
			if r.Identifier == TestClientIdentifier {
				found = true
			}
		}
		assert.True(t, found, "Expected successful client registration in server.")
	}

	// Test invalid json
	{
		resp := testApi.callJSON("POST", reqUrl, TestClientSecret, nil)
		assert.Equal(t, http.StatusBadRequest, resp.Code)
	}

	// Test invalid client secret
	{
		resp := testApi.callJSON("POST", reqUrl, "wrong secret", m)
		assert.Equal(t, http.StatusUnauthorized, resp.Code)
	}

	// Test invalid widget
	{
		im := m
		im["availableWidgets"] = []server.WidgetRegistration{
			{Type: "invalid type", Configuration: []byte("null")},
		}
		resp := testApi.callJSON("POST", reqUrl, TestClientSecret, im)
		assert.Equal(t, http.StatusBadRequest, resp.Code)
	}
}

func TestRequestedClientWidgets(t *testing.T) {
	// Setup
	testApi := validApi()
	reqUrl := fmt.Sprintf("/api/clients/%s/requested_widgets", TestClientIdentifier)

	// Mock scheduler
	updateRequests := make(chan []string, 1)
	sched := &mockScheduler{
		MockRequestUpdateRequest: func(client string) chan []string {
			return updateRequests
		},
	}
	testApi.Scheduler = sched

	// Register client
	r := &server.ClientRegistration{
		Name:             TestClientName,
		Identifier:       TestClientIdentifier,
		AvailableWidgets: nil,
	}
	err := testApi.Server.RegisterClient(r)
	require.Nil(t, err, "Unexpected error: %v", err)

	// Test all valid
	{
		request := []string{"load", "network"}
		updateRequests <- request
		resp := testApi.call("GET", reqUrl, TestClientSecret, nil)
		assert.Equal(t, http.StatusOK, resp.Code)

		var m struct {
			Widgets []string `json:"widgets"`
		}
		err := json.NewDecoder(resp.Body).Decode(&m)
		assert.Nil(t, err, "Could not decode body: %v", err)
		assert.True(t, reflect.DeepEqual(request, m.Widgets))
	}

	// Test scheduler timeout
	{
		testApi.Configuration.MaximumClientUpdateInterval = 1 * time.Nanosecond
		resp := testApi.call("GET", reqUrl, TestClientSecret, nil)
		assert.Equal(t, http.StatusOK, resp.Code)

		var m struct {
			Widgets []string `json:"widgets"`
		}
		err := json.NewDecoder(resp.Body).Decode(&m)
		assert.Nil(t, err, "Could not decode body: %v", err)
		assert.Equal(t, 0, len(m.Widgets), "Expected empty widget list.")
	}

	// Test invalid client secret
	{
		resp := testApi.call("GET", reqUrl, "invalid secret", nil)
		assert.Equal(t, http.StatusUnauthorized, resp.Code)
	}

	// Test invalid client
	{
		reqUrl := "/api/clients/unkown-client/requested_widgets"
		resp := testApi.call("GET", reqUrl, "", nil)
		assert.Equal(t, http.StatusNotFound, resp.Code)
	}
}

// mockScheduler to be used as a mock object for the scheduler.
type mockScheduler struct {
	MockRequestUpdateRequest func(client string) chan []string
}

func (m *mockScheduler) RegisterClient(client string) {}
func (m *mockScheduler) HasClient(client string) bool {
	return false
}
func (m *mockScheduler) FulfillUpdateRequest(client string, updated []string) {}
func (m *mockScheduler) RequestWidgets(client string, widgets []string, immediately bool) chan bool {
	return nil
}

func (m *mockScheduler) RequestUpdateRequest(client string) chan []string {
	if m.MockRequestUpdateRequest != nil {
		return m.MockRequestUpdateRequest(client)
	}
	return nil
}
