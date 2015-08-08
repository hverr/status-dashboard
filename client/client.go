package client

import (
	"errors"

	"github.com/hverr/status-dashboard/widgets"
	"github.com/jmcvetta/napping"
)

type RequestedWidgets struct {
	Widgets []string `json:"widgets"`
}

func GetRequestedWidgets() (RequestedWidgets, error) {
	widgets := RequestedWidgets{}

	url := Configuration.API + "/clients/" + Configuration.Identifier + "/widgets/requested"
	resp, err := napping.Get(url, nil, &widgets, nil)
	if err != nil {
		return widgets, err
	} else if resp.HttpResponse().StatusCode != 200 {
		return widgets, errors.New("Could not get requested widgets: " + resp.HttpResponse().Status)
	}

	return widgets, nil
}

func PostWidgetUpdate(widget widgets.Widget) error {
	t := widget.Type()
	url := Configuration.API + "/clients/" + Configuration.Identifier + "/widgets/" + t + "/update"

	resp, err := napping.Post(url, &widget, nil, nil)
	if err != nil {
		return err
	} else {
		return errors.New("Could not post widget update " + resp.HttpResponse().Status)
	}

	return nil
}
