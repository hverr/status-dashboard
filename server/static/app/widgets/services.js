'use strict';

angular.module('dashboard').factory('widgetsManager', [
  '$interval',
  'api',
  '$log',
  function($interval, api, $log) {
    var widgets = [];

    return {
      register : function(widget) {
        $log.info('Registering widget', widget.type, 'for', widget.client);
        widgets.push(widget);
      },

      registeredWidgets : function() {
        return widgets;
      },

      start : function() {
        var f = function() {
          widgets.forEach(function(w) {
            api.widget(w).then(function(result) {
              w.update(result);
            });
          });
        };

        f();
        $interval(f, 5*1000);
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

      widget: function(w) {
        var deferred = $q.defer();

        var r = this.resource('/clients/' + w.client + '/widgets/' + w.identifier);
        $http
          .get(r)
          .then(function(result) {
            $log.info('Updated', w.identifier, 'for', w.client, ':', result);
            deferred.resolve(result.data);

          }, function(reason) {
            $log.error('error: api.widget:', reason);
          });

        return deferred.promise;
      },
    };
  }
]);
