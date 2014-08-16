/**
 * Create a system with a name and update url
 *
 * @constructor
 * @param {string} name - The humun-readable of the system
 * @param {string} updateURL - The URL to query for system information
 **/
function System(name, updateURL) {
	this.name = name;
	this.updateURL = updateURL;
	this.components = [];
}

/**
 * General component class
 * @constructor
 **/
function Component() {
	this.key = null;
}

/**
 * Create a component holding memory information about the system
 *
 * @constructor
 * @param {object} [jsonObject] - The JSON object from which to construct
 *    the componenet.
 **/
function MemoryComponent(jsonObject) {
	this.key = 'memory';
	this.total = null;
	this.used = null;

	this.loadJSON = function(jsonObject) {
		if(!jsonObject[this.key]) {
			throw new Error("Expected a member called '" + this.key + "'");
		}

		var memObj = jsonObject[this.key];
		if(typeof memObj['total'] === 'undefined') {
			throw new Error("Expected a member called 'total'");
		} else if(typeof memObj['used'] === 'undefined') {
			throw new Error("Expected a member called 'used'");
		}

		this.total = parseInt(memObj['total']);
		this.used = parseInt(memObj['used']);
	}

	if(jsonObject) {
		this.loadJSON(jsonObject);
	}
}
MemoryComponent.prototype = new Component();

function LoadComponent(jsonObject) {
	this.key = 'load';
	this.cores = null;
	this.load = null;

	this.loadJSON = function(jsonObject) {
		if(!jsonObject[this.key]) {
			throw new Error("Expected a member called '" + this.key + "'");
		}

		var loadObj = jsonObject[this.key];
		if(typeof loadObj['cores'] == 'undefined') {
			throw new Error("Expected a member called 'cores'");
		} else if(typeof loadObj['load'] == 'undefined') {
			throw new Error("Expected a member called 'load'");
		} else if(loadObj['load'].length != 3) {
			throw new Error("Expected 'load' member to have 3 objects");
		}

		this.cores = parseInt(loadObj['cores']);

		var load = loadObj['load'];
		this.load = [parseFloat(load[0]), parseFloat(load[1]), parseFloat(load[2])];
	}

	if(jsonObject) {
		this.loadJSON(jsonObject);
	}
}

function UpTimeComponent(jsonObject) {
	this.key = 'uptime';
	this.days = null;
	this.hours = null;
	this.minutes = null;
	this.seconds = null;

	this.addOneSecond = function() {
		if(this.seconds === null) {
			return;
		}

		this.seconds += 1;
		if(this.seconds == 60) {
			this.seconds = 0;
			this.minutes += 1;
		}
		if(this.minutes == 60) {
			this.hours += 1;
			this.minutes = 0;
		}
		if(this.hours == 24) {
			this.days += 1;
			this.hours = 0;
		}
	}

	this.loadJSON = function(jsonObject) {
		if(!jsonObject[this.key]) {
			throw new Error("Expected a member called '" + this.key + "'");
		}

		var uptimeObj = jsonObject[this.key];
		this.days = uptimeObj.days;
		this.hours = uptimeObj.hours;
		this.minutes = uptimeObj.minutes;
		this.seconds = uptimeObj.seconds;
	}

	if(jsonObject) {
		this.loadJSON(jsonObject);
	}
}

/**
 * Create a class that handles network request to update system information
 **/
function SystemUpdater(system) {
	this.system = system;

	/**
	 * @callback SystemUpdaterCallback
	 *
	 * @param {string} [error] - If the update was not successful, this
	 *     parameter will contain the error description
	 **/
	/**
	 * Update the specified components of this.system
	 *
	 * @function updateComponents
	 * @param {string[]} componentKeys - The components to update, specified by
	 *     their keys
	 * @param {SystemUpdaterCallback} callback - The callback to call when the
	 *     update was completed
	 **/
	this.updateComponents = function(componentKeys, callback) {
		var system = this.system;
		var updateQuery = this.updateQueryForComponents(componentKeys);
		var requestURL = system.updateURL + updateQuery;

		var err = "Could not query " + system.updateURL + ": ";
		console.log("Querying " + requestURL);
		$.get(requestURL).done(function(data, textStatus, jqXHR) {
			console.log("--> " + data);

			/* Parse result */
			var result = JSON.parse(data);
			if(!result) {
				callback(err + "could not load JSON");
				return;
			}

			for(key in result) {
				system.components[key].loadJSON(result);
			}
			callback();

		}).fail(function(jqXHR, textStatus, errorThrown) {
			console.log("--> " + jqXHR.status);
			callback(err + "(" + jqXHR.status + ") " + jqXHR.statusText);
		});
	}

	this.updateQueryForComponents = function(componentKeys) {
		if(componentKeys.length == 0) {
			return "";
		}

		var queryComponents = []
		for(var i = 0; i < componentKeys.length; i++) {
			var encKey = encodeURIComponent(componentKeys[i]);
			queryComponents.push(i.toString() + "=" + encKey);
		}
		return "?" + queryComponents.join('&');
	}
}

/** Widgets **/
function Widget(system, component) {
	this.system = system;
	this.component = component;
	this.widget = null;
	this.id = null;

	this.attach = function(grid, widgetID) {
		this.id = widgetID;
		this.widget = grid.add_widget(this.html());
	}

	this.html = function(widgetID) {
		throw new Error("Method should be overridden");
	}

	this.isAvailable = function() {
		throw new Error("Method should be overridden");
	}

	this.refresh = function() {
		if(!this.isAvailable()) {
			this.disableContent();
			this.enableNotAvailable();
		} else {
			this.disableNotAvailable();
			this.enableContent();
		}
	}

	this.enableContent = function() {
		throw new Error("Method should be overridden");
	}

	this.disableContent = function() {
		throw new Error("Method should be overridden");
	}

	this.enableNotAvailable = function() {
		var notAvailableDiv = this.widget.find('.notavailable');
		if(!notAvailableDiv) {
			var html = '<div class="notavailable">N/A</div>';
			this.widget.find('.content').append(html);
		}
	}

	this.disableNotAvailable = function() {
		this.widget.find('.notavailable').remove();
	}

	this.enableContent = function() {}
}

function MemoryWidget(system) {
	Widget.call(this, system, new MemoryComponent());

	this.isAvailable = function() {
		return (this.component.total !== null && this.component.used !== null);
	}

	this.enableContent = function() {
		var usedSpan = this.widget.find('span.used');
		var totalSpan = this.widget.find('span.total');
		if(!usedSpan.length || !totalSpan.length) {
			usedSpan.remove();
			totalSpan.remove();
			this.widget.find('.content').append([
				'<span class="used"></span><br>',
				'<span class="total"></span>'
			].join('\n'));

			usedSpan = this.widget.find('span.used');
			totalSpan = this.widget.find('span.total');
		}

		var used = this.component.used;
		var total = this.component.total;

		usedSpan.text(used.toString() + " MB");
		totalSpan.text(total.toString() + " MB");

		if(used/total < 0.8) {
			usedSpan.removeClass('red');
			usedSpan.addClass('green');
		} else {
			usedSpan.removeClass('green');
			usedSpan.addClass('red');
		}
	}

	this.disableContent = function() {
		this.widget.find('span.used').remove();
		this.widget.find('span.total').remove();
	}

	this.html = function() {
		return [
		'<div class="new widget memory" id="' + this.id + '">',
		    '<div class="content">',
			'</div>',
			'<div class="system-info">',
				this.system.name.toUpperCase() + " - " +
					this.component.key.toUpperCase(),
			'</div>',
		'</div>'
		].join('\n');
	}
}
MemoryWidget.prototype = new Widget();

function LoadWidget(system) {
	Widget.call(this, system, new LoadComponent());

	this.isAvailable = function() {
		return (this.component.cores != null && this.component.load != null);
	}

	this.enableContent = function() {
		var coresDiv = this.widget.find('.cores');
		var loadDiv = this.widget.find('.load');
		if(!coresDiv.length || !loadDiv.length) {
			coresDiv.remove();
			loadDiv.remove();
			this.widget.find('.content').append([
				'<span class="cores"></span><br>',
				'<div class="load">',
				'<span style="margin-right:5px;"></span>',
				'<span style="margin-right:5px;"></span>',
				'<span></span>',
				'</div>'
			].join('\n'));

			coresDiv = this.widget.find('.cores');
			loadDiv = this.widget.find('.load');
		}

		var cores = this.component.cores;
		var load = this.component.load;

		coresDiv.text(cores.toString() + " cores");
		loadSpans = loadDiv.first().find('span');
		for(var i = 0; i < loadSpans.length; i++) {
			loadSpans.eq(i).text(load[i].toFixed(2));
			if(load[i] > cores) {
				loadSpans.eq(i).removeClass('green');
				loadSpans.eq(i).addClass('red');
			} else {
				loadSpans.eq(i).removeClass('red');
				loadSpans.eq(i).addClass('green');
			}
		}
	}

	this.disableContent = function() {
		this.widget.find('.cores').remove();
		this.widget.find('.load').remove();
	}

	this.html = function() {
		return [
		'<div class="new widget load" id="' + this.id + '">',
		    '<div class="content">',
			'</div>',
			'<div class="system-info">',
				this.system.name.toUpperCase() + " - " +
					this.component.key.toUpperCase(),
			'</div>',
		'</div>'
		].join('\n');
	}
}
LoadWidget.prototype = new Widget();

function UpTimeWidget(system) {
	Widget.call(this, system, new UpTimeComponent());

	var that = this;
	this.updateTime = function() {
		that.component.addOneSecond();
		that.refresh();
		setTimeout(function() {
			that.updateTime();
		}, 1000);
	}

	this.attach = function(grid, widgetID) {
		this.__proto__.attach.call(this, grid, widgetID);
		this.updateTime();
	}

	this.isAvailable = function() {
		var comp = this.component;
		return (comp.days || comp.hours || comp.minutes || comp.seconds);
	}

	this.enableContent = function() {
		var uptimeSpan = this.widget.find('.uptime');
		if(!uptimeSpan.length) {
			this.widget.find('.content').append([
				'<span class="uptime">',
				'<span class="days"></span><br>',
				'<span class="time"></span>',
				'</span>'
			].join('\n'));
			uptimeSpan = this.widget.find('.uptime');
		}

		var days = ""
		if(this.component.days) {
			days += this.component.days;
		} else {
			days += "0";
		}
		days += " days";
		uptimeSpan.first().find('.days').text(days);

		var string = "";
		if(this.component.hours) {
			string += this.component.hours;
		} else {
			string += "0";
		}
		string += "h ";

		if(this.component.minutes) {
			string += this.component.minutes;
		} else {
			string += "0";
		}
		string += "m ";

		if(this.component.seconds) {
			string += this.component.seconds;
		} else {
			string += "0";
		}
		string += "s";

		uptimeSpan.first().find('.time').text(string);
	}

	this.disableContent = function() {
		this.widget.find('.uptime').remove();
	}

	this.html = function() {
		return [
		'<div class="new widget uptime" id="' + this.id + '">',
		    '<div class="content">',
			'</div>',
			'<div class="system-info">',
				this.system.name.toUpperCase() + " - " +
					this.component.key.toUpperCase(),
			'</div>',
		'</div>'
		].join('\n');
	}
}
UpTimeWidget.prototype = new Widget();
