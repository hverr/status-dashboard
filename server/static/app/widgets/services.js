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
  '$rootScope',
  'api',
  'widgetFactory',
  '$log',
  function($timeout, $q, $rootScope, api, widgetFactory, $log) {
    var self = {
      start: start,
      update: update,

      add: add,
      availableClients: null,

      availableClientsChangedEvent: 'AvailableClientsChangedEvent',
      addWidgetRequestEvent: 'AddWidgetRequestEvent',
    };

    var widgets = {};
    var lastUpdateCall = new Date();

    function start() {
      api.availableClients().then(function(clients) {
        self.availableClients = clients;
        $rootScope.$emit(self.availableClientsChangedEvent);
      });

      var force = true;
      function f() {
        $log.debug('Updating');
        update(force).then(function() {
          force = false;
          $timeout(f, 100);
        }, function() {
          force = false;
          $timeout(f, 10*1000);
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

    function update(force) {
      var done = $q.defer();

      var thisUpdateCall = new Date();
      lastUpdateCall = thisUpdateCall;

      var request = {};
      for(var clientIdentifier in widgets) {
        request[clientIdentifier] = [];
        for(var widgetType in widgets[clientIdentifier]) {
          request[clientIdentifier].push(widgetType);
        }
      }

      api.updateRequest(force, request).then(function(result) {
        if(lastUpdateCall === thisUpdateCall) {
          for(var clientIdentifier in widgets) {
            request[clientIdentifier] = [];
            for(var widgetType in widgets[clientIdentifier]) {
              var widget = widgets[clientIdentifier][widgetType];

              if(!(clientIdentifier in result)) {
                widget.available = false;
              } else if(!(widgetType in result[clientIdentifier])) {
                widget.available = false;
              } else if(!result[clientIdentifier][widgetType]) {
                widget.available = false;
              } else {
                widget.update(result[clientIdentifier][widgetType]);
                widget.available = true;
              }
            }
          }
        }

        done.resolve();
      }, function(reason) {
        done.reject(reason);
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

    function updateRequest(force, widgets) {
      var d = $q.defer();

      var r = resource('/update_request');
      if(force === true) {
        r += '?force=true';
      }

      $http.post(r, widgets).then(function(result) {
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
