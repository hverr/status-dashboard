/* ===== SETUP GRID ===== */
var g_id_counter = 1;
var gridoptions = {
	widget_selector : "div.widget",
	widget_margins : [20, 20],
	widget_base_dimensions : [250, 250]
};
var g_grid = $(".gridster ul").gridster(gridoptions).data('gridster');
g_grid.disable();

/* ===== MOCKING ===== */
/*$.mockjax({
	url : '*',
	responseText : JSON.stringify({
		memory : {
			'total' : 500,
			'used' : 0
		}
	})
});*/

/* ===== CONFIGURATION ===== */
var g_systems = [];
var g_widgets = [];
var g_configuration = {
	refreshTimeout : 1000*60, /* one minute */
	systems : [
		{
			name : 'webserver',
			updateURL : 'http://status.web.example.org',
			widgets : [
				MemoryWidget,
				LoadWidget,
				UpTimeWidget
			]
		},
        {
            name : 'ftpserver',
            updateURL : 'http://status.ftp.example.org',
            widgets : [
                MemoryWidget,
                LoadWidget,
                UpTimeWidget
            ]
        }
	]
};

function SetupSystemsAndWidgets() {
	var systems = g_configuration['systems'];
	for(var systemIndex = 0; systemIndex < systems.length; systemIndex++) {
		var systemInfo = systems[systemIndex];
		var system = new System(systemInfo.name, systemInfo.updateURL);

		var widgets = systemInfo.widgets;
		for(var widgetIndex = 0; widgetIndex < widgets.length; widgetIndex++) {
			var widget = new widgets[widgetIndex](system);

			system.components[widget.component.key] = widget.component;
			g_widgets.push(widget);
		}

		g_systems.push(system);
	}
}

/* ===== MAIN ===== */
$(document).ready(function() {
	SetupSystemsAndWidgets();

	for(var i = 0; i < g_widgets.length; i++) {
		var id = "widget-" + g_id_counter;
		g_widgets[i].attach(g_grid, id);
		g_id_counter++;
	}

	var refresh = function() {
		RefreshSystems();
		setTimeout(function() {
			refresh();
		}, g_configuration['refreshTimeout']);
	};
	refresh();
});

function RefreshWidgets() {
	for(var i = 0; i < g_widgets.length; i++) {
		g_widgets[i].refresh();
	}
}

function RefreshSystems() {
	for(var i = 0; i < g_systems.length; i++) {
		var system = g_systems[i];
		var su = new SystemUpdater(system);
		su.updateComponents([], function(error) {
			if(error) {
				ReportError(error);
			}
			RefreshWidgets();
		});
	}
}

function ReportError(e) {
    var error = "Error at " + new Date() + ": " + e;
    console.error(error);

    var div = $('<div/>');
    div.text(error);
    div.attr('class', 'alert alert-danger');
    div.attr('role', 'alert');
    $(lasterror).html(div.html());
    $(lasterror).show();
    setTimeout(function() {
        $(lasterror).hide();
    }, 10*1000);
}
