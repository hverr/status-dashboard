package widgets

// Widget defines the interface for a general widget.
type Widget interface {
	// Name should return a human-readable name for this widget describing what
	// it does.
	Name() string

	// Type should return an identifier for this type of widget. It should be
	// unique among other widgets.
	Type() string

	// HasData should return whether the widget has data to display. This should
	// default to false and become true when the widget gets information from
	// the client for the first time.
	HasData() bool
}

type WidgetInitiator func() Widget

var AllWidgets = map[string]WidgetInitiator{
	LoadWidgetType: func() Widget { return &LoadWidget{} },
}
