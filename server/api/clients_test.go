package api

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/hverr/status-dashboard/server"
	"github.com/hverr/status-dashboard/server/scheduler"
	"github.com/stretchr/testify/assert"
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
	sched := scheduler.New()
	userAuth := server.UserAuthenticator{conf}

	return &API{serv, userAuth, sched}
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
