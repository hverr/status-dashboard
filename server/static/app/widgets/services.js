'use strict';

angular.module('dashboard').factory('widgetsManager', [
  '$timeout',
  'api',
  '$log',
  function($timeout, api, $log) {
    var clients = {};
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

      availableWidgets: availableWidgets,
    };

    function resource(path) {
      return self.baseURL + path;
    }

    function availableWidgets() {
      var d = $q.defer();

      $http.get(resource('/available_widgets')).then(function(result) {
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
