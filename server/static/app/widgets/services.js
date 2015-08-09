'use strict';

angular.module('dashboard').factory('widgetsManager', [
  '$log',
  function($log) {
    var widgets = [];

    return {
      'register' : function(widget) {
        $log.info('Registering widget', widget.type, 'for', widget.client);
        widgets.push(widget);
      },
    };
  }
]);

angular.module('dashboard').factory('Widget', [
  function() {
    /**
     * Initiate a new Widget instance.
     *
     * @param {string} client The unique identifer of the client to which
     * this widget belongs. (e.g. `'webserver'`)
     *
     * @param {string} type The type of the widget. (e.g. `'uptime'`)
     */
    return function(client, type, gridItem) {
      this.client = client;
      this.type = type;
      this.gridItem = gridItem;
    };
  }
]);
