'use strict';

angular.module('dashboard').factory('widgetsManager', [
  '$timeout',
  'api',
  '$log',
  function($timeout, api, $log) {
    var clients = {};

    function registeredWidgets() {
      var widgets = [];
      for(var clientIdentifier in clients) {
        var client = clients[clientIdentifier];
        for(var widgetType in client) {
          widgets.push(client[widgetType]);
        }
      }

      return widgets;
    }

    return {
      register : function(widget) {
        $log.info('Registering', widget);
        if(clients[widget.client] === undefined) {
          clients[widget.client] = {};
        }
        clients[widget.client][widget.identifier] = widget;
      },

      registeredWidgets : registeredWidgets,

      start : function() {
        $log.debug('Should start updating');
      },
    };
  }
]);

angular.module('dashboard').factory('api', [
  '$q',
  '$http',
  '$log',
  function($q, $http, $log) {
    var self = {
      baseURL : '/api',

    };

    function resource(path) {
      return self.baseURL + path;
    }

    return self;
  }
]);
