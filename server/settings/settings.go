package settings

import "time"

var ListenAddress = ":8050"
var StaticAppRoot = "../../app"

var MinimumClientUpdateInterval = 3 * time.Second
var MaximumClientUpdateInterval = 5 * time.Minute

var MaximumWidgetAge = 30 * time.Second
