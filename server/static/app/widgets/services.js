'use strict';

angular.module('dashboard').factory('widgetFactory', [
  'LoadWidget',
  'UptimeWidget',
  function(LoadWidget, UptimeWidget) {
    return function(widgetType) {
      switch(widgetType) {
        case 'load':
          return new LoadWidget();
        case 'uptime':
          return new UptimeWidget();
        default:
          return null;
      }
    };
  }
]);

angular.module('dashboard').factory('widgetsManager', [
  '$timeout',
  'api',
  'widgetFactory',
  '$log',
  function($timeout, api, widgetFactory, $log) {
    var self = {
      start: start,

      add: add,
    };

    var availableClients;
    var widgets = {};

    function start() {
      api.availableClients().then(function(clients) {
        $log.debug('Got available clients:', clients);
        availableClients = clients;
      });
    }

    function add(clientIdentifier, widgetType) {
      if(!widgets[clientIdentifier]) {
        widgets[clientIdentifier] = {};
      }

      if(!widgets[clientIdentifier][widgetType]) {
        var widget = widgetFactory(widgetType);
        if(widget === null) {
          throw 'Unknown widget type: ' + widgetType;
        }

        widget.referenceCounter = 1;
        widget.available = false;
        widgets[clientIdentifier][widgetType] = widget;

      } else {
        widgets[clientIdentifier][widgetType].referenceCounter += 1;
      }

      return widgets[clientIdentifier][widgetType];
    }

    return self;
  }
]);

angular.module('dashboard').factory('api', [
  '$q',
  '$http',
  '$log',
  function($q, $http, $log) {
    var self = {
      baseURL : '/api',

      error: defaultError,

      availableClients: availableClients,
    };

    function resource(path) {
      return self.baseURL + path;
    }

    function availableClients() {
      var d = $q.defer();

      $http.get(resource('/available_clients')).then(function(result) {
        d.resolve(result.data);
      }, function(reason) {
        self.error(reason);
        d.reject(reason);
      });

      return d.promise;
    }

    function defaultError(reason) {
      $log.error('HTTP error:', reason);
    }

    return self;
  }
]);
