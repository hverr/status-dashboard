package client

import (
	"github.com/hverr/status-dashboard/widgets"
	"github.com/jmcvetta/napping"
)

type RequestedWidgets struct {
	Widgets []string `json:"widgets"`
}

func GetRequestedWidgets() (RequestedWidgets, error) {
	widgets := RequestedWidgets{}

	url := Configuration.API + "/clients/" + Configuration.Identifier + "/widgets/requested"
	_, err := napping.Get(url, nil, &widgets, nil)
	return widgets, err
}

func PostWidgetUpdate(widget widgets.Widget) error {
	t := widget.Type()
	url := Configuration.API + "/clients/" + Configuration.Identifier + "/widgets/" + t + "/update"

	_, err := napping.Post(url, &widget, nil, nil)
	return err
}
