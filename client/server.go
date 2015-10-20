package client

import (
	"errors"
	"net/http"

	"github.com/hverr/status-dashboard/server"
	"github.com/hverr/status-dashboard/widgets"
	"github.com/jmcvetta/napping"
)

type Server struct {
	Session       napping.Session
	Configuration Configuration
}

type RequestedWidgets struct {
	Widgets []string `json:"widgets"`
}

func (s *Server) Register(availableWidgets []server.WidgetRegistration) error {
	payload := server.ClientRegistration{
		Name:             s.Configuration.Name,
		Identifier:       s.Configuration.Identifier,
		AvailableWidgets: availableWidgets,
	}

	resource := "/clients/" + s.Configuration.Identifier + "/register"
	resp, err := s.send("POST", resource, &payload, nil, nil)
	if err != nil {
		return err
	} else if resp.HttpResponse().StatusCode != 200 {
		return errors.New("Could not register client: " + resp.HttpResponse().Status)
	}

	return nil
}

func (s *Server) GetRequestedWidgets() (RequestedWidgets, error) {
	widgets := RequestedWidgets{}

	resource := "/clients/" + s.Configuration.Identifier + "/requested_widgets"
	resp, err := s.send("GET", resource, nil, &widgets, nil)
	if err != nil {
		return widgets, err
	} else if resp.HttpResponse().StatusCode != 200 {
		return widgets, errors.New("Could not get requested widgets: " + resp.HttpResponse().Status)
	}

	return widgets, nil
}

func (s *Server) PostWidgetBulkUpdate(updates []widgets.BulkElement) error {
	resource := "/clients/" + s.Configuration.Identifier + "/bulk_update"
	resp, err := s.send("POST", resource, updates, nil, nil)
	if err != nil {
		return err
	} else if resp.HttpResponse().StatusCode != 200 {
		return errors.New("Could not post bulk update " + resp.HttpResponse().Status)
	}

	return nil
}

func (s *Server) request(method, resource string, payload, result, errMsg interface{}) *napping.Request {
	url := s.Configuration.API + resource

	req := napping.Request{
		Url:     url,
		Method:  method,
		Payload: payload,
		Result:  result,
		Error:   errMsg,
	}

	if s.Configuration.Secret != "" {
		req.Header = &http.Header{}
		req.Header.Set("X-Client-Secret", s.Configuration.Secret)
	}

	return &req
}

func (s *Server) send(method, resource string, payload, result, errMsg interface{}) (*napping.Response, error) {
	return s.Session.Send(s.request(method, resource, payload, result, errMsg))
}
