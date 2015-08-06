package widgets

// Widget defines the interface for a general widget.
type Widget interface {
	// System should return a human-readable name of the system to which this
	// widget belongs.
	System() string

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
