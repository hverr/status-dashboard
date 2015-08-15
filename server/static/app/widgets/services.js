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
        var update = function(force) {
          var request = {};

          for(var clientIdentifier in clients) {
            request[clientIdentifier] = [];
            for(var widgetType in clients[clientIdentifier]) {
              request[clientIdentifier].push(widgetType);
            }
          }

          api.bulk(request, force).then(function(result) {
            if(!result) {
              // A timeout occurred.
              $timeout(update, 100);
              return;
            }

            for(var clientIdentifier in result) {
              if(!clients[clientIdentifier]) {
                continue;
              }

              var client = clients[clientIdentifier];
              for(var widgetType in result[clientIdentifier]) {
                if(!client[widgetType]) {
                  continue;
                }

                client[widgetType].update(result[clientIdentifier][widgetType]);
              }
            }

            $timeout(update, 100);
          }, function() {
            $timeout(update, 10*1000);
          });
        };

        update(true);
      },
    };
  }
]);

angular.module('dashboard').factory('api', [
  '$q',
  '$http',
  '$log',
  function($q, $http, $log) {
    var baseURL = '/api';

    return {
      resource: function(resource) {
        return baseURL + resource;
      },

      bulk: function(widgets, force) {
        var deferred = $q.defer();

        var r = this.resource('/all_widgets');
        if(force) {
          r += '?force=true';
        }

        $http
          .post(r, widgets)
          .then(function(result) {
            $log.debug('Updated widgets:', result);
            if(result.status === 202) {
              // Long poll timout
              deferred.resolve({});

            } else {
              deferred.resolve(result.data);
            }

          }, function(reason) {
            $log.error('error: api.bulk:', reason);
            deferred.reject(reason);
          });

        return deferred.promise;
      },
    };
  }
]);
