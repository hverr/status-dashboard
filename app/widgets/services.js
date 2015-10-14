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
        configuration: null,

        data: null,
        update: function(data) {
          self.data = data;
        },
        configure: function(c) {
          self.configuration = c;
        },
        identifier: function() {
          return self.type;
        },

        watchValue: function() {
          return {
            directive: self.directive,
            height: self.height, width: self.width,
            row: self.row, col: self.col,
            type: self.type, clientIdentifier: self.clientIdentifier,
            client: self.client, name: self.name,
          };
        },
      };

      return self;
    };
  }
]);

angular.module('dashboard').factory('widgetFactory', [
  'LoadWidget',
  'UptimeWidget',
  'MeminfoWidget',
  'CurrentTimeWidget',
  'CurrentDateWidget',
  'ConnectionsWidget',
  function(LoadWidget,
           UptimeWidget,
           MeminfoWidget,
           CurrentTimeWidget,
           CurrentDateWidget,
           ConnectionsWidget)
  {
    return function(widgetType) {
      switch(widgetType) {
        case 'load':
          return new LoadWidget();
        case 'uptime':
          return new UptimeWidget();
        case 'meminfo':
          return new MeminfoWidget();
        case 'current_time':
          return new CurrentTimeWidget();
        case 'current_date':
          return new CurrentDateWidget();
        case 'connections':
          return new ConnectionsWidget();
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
      updateAvailableClients : updateAvailableClients,

      add: add,
      remove: remove,
      availableClients: null,

      addFrom: addFrom,

      availableClientsChangedEvent: 'AvailableClientsChangedEvent',
      addWidgetRequestEvent: 'AddWidgetRequestEvent',
    };

    var widgets = [];
    var lastUpdateCall = new Date();

    function start() {
      oneSecondService.start();

      updateAvailableClients();

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

    function updateAvailableClients() {
      api.availableClients().then(function(clients) {
        self.availableClients = clients;
        $rootScope.$emit(self.availableClientsChangedEvent);
      });
    }

    function add(clientIdentifier, widget) {
      var w = widgetFactory(widget.type);
      if(w === null) {
        throw 'Unknown widget type: ' + widget.type;
      }

      if(widget.configuration) {
        w.configure(widget.configuration);
      }

      w.clientIdentifier = clientIdentifier;
      w.type = widget.type;
      if(clientIdentifier !== "") {
        // only change widgets for remote clients
        w.available = false;
      }

      w.widgetsManagerHandle = widgets.length;
      widgets.push(w);

      return w;
    }

    function remove(widget) {
      var handle = widget.widgetsManagerHandle;

      widgets.splice(handle, 1);
      if(widgets.length > 0 && widgets.length !== handle) {
        widgets.slice(handle, widgets.length).forEach(function(w) {
          w.widgetsManagerHandle -= 1;
        });
      }
    }

    function update(force) {
      var done = $q.defer();

      var thisUpdateCall = new Date();
      lastUpdateCall = thisUpdateCall;

      var request = {};
      widgets.forEach(function(w) {
        if(w.clientIdentifier === "") {
          return; // built-in widgets don't need to be updated
        }

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

            if(clientIdentifier === "") {
              return; // ignore built-in widgets
            }

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

    function addFrom(json) {
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

angular.module('dashboard').factory('widgetsStore', [
  '$location',
  '$cookies',
  function($location, $cookies) {
    var self = {
      columns: 0,
      widgets: [],

      serialize: serialize,
      deserialize: deserialize,
      saveLayout: saveLayout,
      loadLayout: loadLayout,
      clearLayout: clearLayout,
    };

    function serialize() {
      var json = [];
      self.widgets.forEach(function(w) {
        json.push({
          client: w.clientIdentifier,
          type: w.type,

          height: w.height,
          width: w.width,
          row: w.row,
          col: w.col,
        });
      });
      return {
        columns: self.columns,
        widgets: json,
      };
    }

    function deserialize(json) {
      self.columns = json.columns;
      self.widgets = json.widgets;
    }

    function saveLayout() {
      var data = serialize();
      var expires = new Date();
      expires.setFullYear(expires.getFullYear() + 10);
      $cookies.put('layout', angular.toJson(data), {
        expires: expires,
      });
    }

    function loadLayout() {
      var query = $location.search();
      if('layout' in query) {
        return query.layout;
      }

      var data = $cookies.get('layout');

      return data ? data : null;
    }

    function clearLayout() {
      $cookies.remove('layout');
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
