package client

import (
	"errors"
	"net/http"

	"github.com/hverr/status-dashboard/server"
	"github.com/hverr/status-dashboard/widgets"
	"github.com/jmcvetta/napping"
)

var Session napping.Session

type RequestedWidgets struct {
	Widgets []string `json:"widgets"`
}

func Register(availableWidgets []server.WidgetRegistration) error {
	payload := server.ClientRegistration{
		Name:             Configuration.Name,
		Identifier:       Configuration.Identifier,
		AvailableWidgets: availableWidgets,
	}

	resource := "/clients/" + Configuration.Identifier + "/register"
	resp, err := send("POST", resource, &payload, nil, nil)
	if err != nil {
		return err
	} else if resp.HttpResponse().StatusCode != 200 {
		return errors.New("Could not register client: " + resp.HttpResponse().Status)
	}

	return nil
}

func GetRequestedWidgets() (RequestedWidgets, error) {
	widgets := RequestedWidgets{}

	resource := "/clients/" + Configuration.Identifier + "/requested_widgets"
	resp, err := send("GET", resource, nil, &widgets, nil)
	if err != nil {
		return widgets, err
	} else if resp.HttpResponse().StatusCode != 200 {
		return widgets, errors.New("Could not get requested widgets: " + resp.HttpResponse().Status)
	}

	return widgets, nil
}

func PostWidgetUpdate(widget widgets.Widget) error {
	t := widget.Type()

	resource := "/clients/" + Configuration.Identifier + "/widgets/" + t + "/update"
	resp, err := send("POST", resource, &widget, nil, nil)
	if err != nil {
		return err
	} else if resp.HttpResponse().StatusCode != 200 {
		return errors.New("Could not post widget update " + resp.HttpResponse().Status)
	}

	return nil
}

func PostWidgetBulkUpdate(updates []widgets.BulkElement) error {
	resource := "/clients/" + Configuration.Identifier + "/bulk_update"
	resp, err := send("POST", resource, updates, nil, nil)
	if err != nil {
		return err
	} else if resp.HttpResponse().StatusCode != 200 {
		return errors.New("Could not post bulk update " + resp.HttpResponse().Status)
	}

	return nil
}

func request(method, resource string, payload, result, errMsg interface{}) *napping.Request {
	url := Configuration.API + resource

	req := napping.Request{
		Url:     url,
		Method:  method,
		Payload: payload,
		Result:  result,
		Error:   errMsg,
	}

	if Configuration.Secret != "" {
		req.Header = &http.Header{}
		req.Header.Set("X-Client-Secret", Configuration.Secret)
	}

	return &req
}

func send(method, resource string, payload, result, errMsg interface{}) (*napping.Response, error) {
	return Session.Send(request(method, resource, payload, result, errMsg))
}
