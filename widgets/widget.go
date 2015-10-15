package widgets

import "encoding/json"

// Widget defines the interface for a general widget.
type Widget interface {
	// Name should return a human-readable name for this widget describing what
	// it does.
	Name() string

	// Type should return the type of the widget.
	// unique among other widgets.
	Type() string

	// Identifier should return an unique identifier for the widget. It should
	// be based on its type and configuration.
	Identifier() string

	// HasData should return whether the widget has data to display. This should
	// default to false and become true when the widget gets information from
	// the client for the first time.
	HasData() bool

	// Update should update widget with the most recent information
	Update() error

	// Configure should configure a widget using data from the configuration
	// file.
	Configure(data json.RawMessage) error

	// Configuration should return the widget configuration
	Configuration() interface{}

	// Start should start any background processes needed to e.g. gather data,
	// if necessary.
	Start() error
}

type WidgetInitiator func() Widget

var AllWidgets = map[string]WidgetInitiator{
	LoadWidgetType:        func() Widget { return &LoadWidget{} },
	UptimeWidgetType:      func() Widget { return &UptimeWidget{} },
	MeminfoWidgetType:     func() Widget { return &MeminfoWidget{} },
	ConnectionsWidgetType: func() Widget { return &ConnectionsWidget{} },
	NetworkWidgetType:     func() Widget { return &NetworkWidget{} },
}

type BulkElement struct {
	Type       string      `json:"type"`
	Identifier string      `json:"identifiier"`
	Widget     interface{} `json:"widget"`
}
