'use strict';

angular.module('dashboard').controller('LoadWidgetController', [
  function() {
  }
]);

angular.module('dashboard').factory('LoadWidget', [
  function() {
    return function() {
      return {
        'one' : null,
        'five' : null,
        'fifteen' : null,
      };
    };
  }
]);

angular.module('dashboard').directive('loadWidget', [
  function() {
    return {
      replace: true,
      templateUrl : 'widgets/load/load.html',
    };
  }
]);
