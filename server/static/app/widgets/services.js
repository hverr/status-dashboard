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

        type: null,
        clientIdentifier: null,
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
  'oneSecondService',
  '$log',
  function($timeout, $q, $rootScope, api, widgetFactory, oneSecondService, $log) {
    var self = {
      start: start,
      update: update,

      add: add,
      availableClients: null,

      serialize: serialize,
      deserialize: deserialize,

      availableClientsChangedEvent: 'AvailableClientsChangedEvent',
      addWidgetRequestEvent: 'AddWidgetRequestEvent',
    };

    var widgets = [];
    var lastUpdateCall = new Date();

    function start() {
      oneSecondService.start();

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
      var w = widgetFactory(widgetType);
      if(w === null) {
        throw 'Unknown widget type: ' + widgetType;
      }

      w.clientIdentifier = clientIdentifier;
      w.type = widgetType;
      w.available = false;

      widgets.push(w);
      return w;
    }

    function update(force) {
      var done = $q.defer();

      var thisUpdateCall = new Date();
      lastUpdateCall = thisUpdateCall;

      var request = {};
      widgets.forEach(function(w) {
        if(!(w.clientIdentifier in request)) {
          request[w.clientIdentifier] = [];
        }

        if(!(w.type in request[w.clientIdentifier])) {
          request[w.clientIdentifier].push(w.type);
        }
      });

      api.updateRequest(force, request).then(function(result) {
        if(lastUpdateCall === thisUpdateCall) {
          widgets.forEach(function(widget) {
            var clientIdentifier = widget.clientIdentifier;
            var widgetType = widget.type;

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
          });
        }

        done.resolve();
      }, function(reason) {
        done.reject(reason);
      });

      return done.promise;
    }

    function serialize() {
      var json = [];
      widgets.forEach(function(w) {
        json.push({
          client: w.clientIdentifier,
          type: w.type,

          height: w.height,
          width: w.width,
          row: w.row,
          col: w.col,
        });
      });
      return json;
    }

    function deserialize(json) {
      var result = [];

      $log.debug('Loading widgets:', json);

      json.forEach(function(data) {
        var w = add(data.client, data.type);
        w.client = data.client;

        w.width = data.width;
        w.height = data.height;
        w.row = data.row;
        w.col = data.col;

        result.push(w);
      });

      return result;
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

angular.module('dashboard').factory('oneSecondService', [
  '$interval',
  function($interval) {
    var self = {
      add: add,
      remove: remove,
      start: start,
    };

    var counter = 0;
    var functions = {};

    function add(f) {
      var handle = counter;
      functions[handle] = f;
      counter += 1;
      return handle;
    }

    function remove(handle) {
      delete functions[handle];
    }

    function start() {
      $interval(function() {
        for(var handle in functions) {
          functions[handle]();
        }
      }, 1000);
    }

    return self;
  }
]);
