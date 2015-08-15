'use strict';

angular.module('dashboard').factory('Widget', [
  function() {
    return function(directive, name) {
      var self = {
        directive: directive,
        height: 1,
        width: 1,
        row: 0,
        col: 0,

        client: null,
        name: name,

        data: null,
        update: function(data) {
          self.data = data;
        },
      };

      return self;
    };
  }
]);

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
  '$q',
  'api',
  'widgetFactory',
  '$log',
  function($timeout, $q, api, widgetFactory, $log) {
    var self = {
      start: start,
      update: update,

      add: add,
    };

    var availableClients;
    var widgets = {};

    function start() {
      api.availableClients().then(function(clients) {
        $log.debug('Got available clients:', clients);
        availableClients = clients;
      });

      function f() {
        update().then(function() {
          $log.debug('Updating');
          $timeout(f, 1000);
        });
      }

      f();
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

    function update() {
      var done = $q.defer();

      var request = {};
      for(var clientIdentifier in widgets) {
        request[clientIdentifier] = [];
        for(var widgetType in widgets[clientIdentifier]) {
          request[clientIdentifier].push(widgetType);
        }
      }

      api.updateRequest(request).then(function(result) {
        $log.debug('Got update:', result);
        for(var clientIdentifier in widgets) {
          request[clientIdentifier] = [];
          for(var widgetType in widgets[clientIdentifier]) {
            var widget = widgets[clientIdentifier][widgetType];

            if(!(clientIdentifier in result)) {
              widget.available = false;
            } else if(!(widgetType in result[clientIdentifier])) {
              widget.available = false;
            } else {
              widget.update(result[clientIdentifier][widgetType]);
              widget.available = true;
            }
          }
        }

        done.resolve();
      }, function() {
        done.resolve();
      });

      return done.promise;
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
      updateRequest: updateRequest,
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

    function updateRequest(widgets) {
      var d = $q.defer();

      $http.post(resource('/update_request'), widgets).then(function(result) {
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
